package pubsub

import "sync"

type Message struct {
	Topic   string
	Content string
}

type Pubsub struct {
	mu          sync.RWMutex
	Subscribers map[string]map[string]Subscriber
	Topics      map[string]bool // Just track if topic exists
	messageChan chan Message    // Central message channel
	shutdown    chan bool
}

func NewPubsub() *Pubsub {
	ps := &Pubsub{
		Subscribers: make(map[string]map[string]Subscriber),
		Topics:      make(map[string]bool),
		messageChan: make(chan Message, 1000),
		shutdown:    make(chan bool),
	}
	go ps.messageDispatcher() // Start the message dispatcher
	return ps
}

func (ps *Pubsub) CreateTopics(topic string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.Subscribers[topic] = make(map[string]Subscriber)
	ps.Topics[topic] = true
}

// Central message dispatcher that broadcasts messages to all subscribers
func (ps *Pubsub) messageDispatcher() {
	for {
		select {
		case <-ps.shutdown:
			return
		case msg := <-ps.messageChan:
			ps.broadcastMessage(msg)
		}
	}
}

// Broadcast message to all subscribers of a topic
func (ps *Pubsub) broadcastMessage(msg Message) {
	ps.mu.RLock()
	subscribers, exists := ps.Subscribers[msg.Topic]
	if !exists {
		ps.mu.RUnlock()
		return
	}

	// Create a copy of subscribers to avoid holding lock during message delivery
	subscriberList := make([]Subscriber, 0, len(subscribers))
	for _, sub := range subscribers {
		subscriberList = append(subscriberList, sub)
	}
	ps.mu.RUnlock()

	// Send message to all subscribers (non-blocking)
	for _, subscriber := range subscriberList {
		if csub, ok := subscriber.(*CSubscriber); ok {
			select {
			case csub.MessageChan <- msg.Content:
			default:
				// Skip if subscriber's channel is full (prevents blocking)
			}
		}
	}
}

// Publish a message to a topic
func (ps *Pubsub) PublishMessage(topic, content string) {
	ps.mu.RLock()
	topicExists := ps.Topics[topic]
	ps.mu.RUnlock()

	if !topicExists {
		return // Topic doesn't exist
	}

	select {
	case ps.messageChan <- Message{Topic: topic, Content: content}:
	default:
		// Handle case where message channel is full
	}
}

func (ps *Pubsub) Shutdown() {
	close(ps.shutdown)
}
