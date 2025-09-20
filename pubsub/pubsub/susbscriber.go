package pubsub

import (
	"fmt"
)

type Subscriber interface {
	Subscribe(p *Pubsub) error
	UnSubscribe(p *Pubsub) error
	ProcessMessage(p *Pubsub)
	Close()
}

type CSubscriber struct {
	Id          string
	Topic       string
	MessageChan chan string // Individual message channel for this subscriber
	Signal      chan bool   // Signal channel for shutdown
}

func CreateNewSubscriber(id, topic string) *CSubscriber {
	return &CSubscriber{
		Id:          id,
		Topic:       topic,
		MessageChan: make(chan string, 100), // Buffered channel for messages
		Signal:      make(chan bool),
	}
}

func (subscriber *CSubscriber) Subscribe(p *Pubsub) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	_, ok := p.Subscribers[subscriber.Topic]
	if !ok {
		fmt.Printf("Error: topic '%s' doesn't exist for subscriber %s\n", subscriber.Topic, subscriber.Id)
		return fmt.Errorf("topic doesnt exist")
	}
	p.Subscribers[subscriber.Topic][subscriber.Id] = subscriber
	fmt.Printf("Subscriber %s subscribed to topic '%s'\n", subscriber.Id, subscriber.Topic)
	return nil
}

func (subscriber *CSubscriber) UnSubscribe(p *Pubsub) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	subscribers, ok := p.Subscribers[subscriber.Topic]
	if !ok {
		return fmt.Errorf("topic doesnt exist")
	}
	_, ok = subscribers[subscriber.Id]
	if ok {
		delete(subscribers, subscriber.Id)
		subscriber.Signal <- true
	}
	fmt.Printf("Subscriber %s unsubscribed from topic '%s'\n", subscriber.Id, subscriber.Topic)
	return nil
}

func (subscriber *CSubscriber) ProcessMessage(p *Pubsub) {
	for {
		select {
		case <-subscriber.Signal:
			fmt.Printf("Shutting down subscriber %s\n", subscriber.Id)
			return
		case msg := <-subscriber.MessageChan:
			fmt.Printf("Subscriber %s received message: %s\n", subscriber.Id, msg)
		}
	}
}

func (subscriber *CSubscriber) Close() {
	close(subscriber.Signal)
	close(subscriber.MessageChan)
}
