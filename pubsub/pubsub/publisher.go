package pubsub

import (
	"fmt"
	"sync"
	"time"
)

type Publisher struct {
	Id    string
	Topic string
}

func CreateNewPublisher(topic string, id string) *Publisher {
	return &Publisher{
		Id:    id,
		Topic: topic,
	}
}

func (c *Publisher) Publish(p *Pubsub, wg *sync.WaitGroup) {
	defer wg.Done()
	topic := c.Topic
	for i := 0; i < 10; i++ { // Reduced from 100 to 10 for cleaner output
		message := fmt.Sprintf("Message #%d from publisher %s", i+1, c.Id)
		p.PublishMessage(topic, message)
		fmt.Printf("Publisher %s published: %s\n", c.Id, message)
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Printf("Publisher %s finished publishing\n", c.Id)
}
