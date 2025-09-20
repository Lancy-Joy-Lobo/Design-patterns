package main

import (
	"fmt"
	pubsub "ps/pubsub"
	"sync"
	"time"
)

func main() {
	fmt.Println("=== Pub-Sub System Demo ===")

	// Create the pub-sub system
	ps := pubsub.NewPubsub()

	// Create topics
	ps.CreateTopics("news")
	ps.CreateTopics("sports")
	fmt.Println("Created topics: news, sports")

	// Create publishers
	p1 := pubsub.CreateNewPublisher("news", "p1")
	p2 := pubsub.CreateNewPublisher("news", "p2")
	p3 := pubsub.CreateNewPublisher("sports", "p3")
	fmt.Println("Created publishers: p1(news), p2(news), p3(sports)")

	// Create subscribers
	s1 := pubsub.CreateNewSubscriber("s1", "news")
	s2 := pubsub.CreateNewSubscriber("s2", "news")
	s3 := pubsub.CreateNewSubscriber("s3", "sports")
	fmt.Println("Created subscribers: s1(news), s2(news), s3(sports)")

	// Subscribe to topics
	s1.Subscribe(ps)
	s2.Subscribe(ps)
	s3.Subscribe(ps)

	// Start message processing for subscribers
	go s1.ProcessMessage(ps)
	go s2.ProcessMessage(ps)
	go s3.ProcessMessage(ps)

	fmt.Println("\n=== Starting Publishers ===")

	// Start publishers
	wg := sync.WaitGroup{}
	wg.Add(3)

	go p1.Publish(ps, &wg)
	go p2.Publish(ps, &wg)
	go p3.Publish(ps, &wg)

	// Let it run for a bit
	time.Sleep(3 * time.Second)

	// Test unsubscription
	fmt.Println("\n=== Testing Unsubscription ===")
	s1.UnSubscribe(ps)

	// Wait for publishers to finish
	wg.Wait()

	fmt.Println("\n=== Cleaning Up ===")
	// Clean shutdown
	s1.Close()
	s2.Close()
	s3.Close()
	ps.Shutdown()

	fmt.Println("Demo completed successfully!")
}
