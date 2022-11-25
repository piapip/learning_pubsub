package pubsub

import (
	"log"
	"sync"
)

type Subscriber struct {
	id      string          // id of the subscriber
	message chan *Message   // message channel
	topics  map[string]bool // topics it is subscribed to.
	active  bool            // the status of the subscriber
	mutex   sync.RWMutex    // lock
}

func NewSubscriber(id string) *Subscriber {
	return &Subscriber{
		id:      id,
		message: make(chan *Message),
		topics:  make(map[string]bool),
		active:  true,
	}
}

type Subscribers map[string]*Subscriber

func (s *Subscriber) AddTopic(topic string) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	s.topics[topic] = true
}

func (s *Subscriber) RemoveTopic(topic string) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if s.topics[topic] {
		delete(s.topics, topic)
	}
}

func (s *Subscriber) GetTopics() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	topics := make([]string, 0)

	for topic, _ := range s.topics {
		topics = append(topics, topic)
	}

	return topics
}

func (s *Subscriber) SelfDestruct() {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	s.active = false
	close(s.message)
}

func (s *Subscriber) Listen() {
	// Listens to the message channel, prints once received.
	for {
		if msg, ok := <-s.message; ok {
			log.Printf("Subscriber %s, receive: %s from topic %s\n", s.id, msg.GetBody(), msg.GetTopic())
		}
	}
}
