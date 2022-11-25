package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/piapip/test-pubsub/pubsub"
)

var availableTopics = map[string]string{
	"BTC": "BITCOIN",
	"ETH": "ETHEREUM",
	"DOT": "POLKADOT",
	"SOL": "SOLANA",
}

func pricePublisher(publisher pubsub.PubSub) {
	topicKeys := make([]string, len(availableTopics))
	topicValues := make([]string, len(availableTopics))

	count := 0
	for k, v := range availableTopics {
		topicKeys[count] = k
		topicValues[count] = v

		count++
	}

	for {
		randValue := topicValues[rand.Intn(len(topicValues))] // all topic values.
		msg := fmt.Sprintf("%f", rand.Float64())
		// fmt.Printf("Publishing %s to %s topic\n", msg, randKey)
		go publisher.Publish(randValue, msg)
		// Uncomment if you want to broadcast to all topics.
		// go broker.Broadcast(msg, topicValues)
		r := rand.Intn(4)
		time.Sleep(time.Duration(r) * time.Second) //sleep for random secs.
	}
}

func main() {
	ps := pubsub.NewPubsub()
	s1 := pubsub.NewSubscriber("1")

	ps.Subscribe(s1, availableTopics["BTC"])
	ps.Subscribe(s1, availableTopics["ETH"])

	s2 := pubsub.NewSubscriber("2")
	ps.Subscribe(s2, availableTopics["ETH"])
	ps.Subscribe(s2, availableTopics["SOL"])

	go (func() {
		// sleep for 3 sec, and then subscribe for topic DOT for s2
		time.Sleep(3 * time.Second)
		ps.Subscribe(s2, availableTopics["DOT"])
	})()

	// s2 Unsubcribe SOL
	go (func() {
		// sleep for 5 sec, and then unsubscribe for topic SOL for s2
		time.Sleep(5 * time.Second)
		ps.Subscribe(s2, availableTopics["SOL"])
		fmt.Printf("Total subscribers for topic ETH is %v\n", ps.CountSubscribers(availableTopics["ETH"]))
	})()

	// Concurrently publish the values.
	go pricePublisher(ps)
	// Concurrently listens from s1.
	go s1.Listen()
	// Concurrently listens from s2.
	go s2.Listen()
	// to prevent terminate
	fmt.Scanln()
	fmt.Println("Done!")
}
