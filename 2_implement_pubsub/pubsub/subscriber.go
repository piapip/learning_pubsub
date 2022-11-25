package pubsub

import (
	"fmt"
	"sync"
)

type subscriber struct {
	active        bool
	name          string
	messageStream chan string
	topics        map[string]bool
	mutexLock     sync.RWMutex
}

func NewSubscriber(name string) PubsubSubscriber {
	return &subscriber{
		name:          name,
		active:        true,
		topics:        make(map[string]bool),
		messageStream: make(chan string),
	}
}

func (s *subscriber) GetTopics() []string {
	return []string{}
}

func (s *subscriber) GetName() (string, bool) {
	if s.active {
		return s.name, true
	}

	return s.name, false
}

func (s *subscriber) AddTopic(topic string) {
	s.mutexLock.Lock()
	defer s.mutexLock.Unlock()

	s.topics[topic] = true
}

func (s *subscriber) RemoveTopic(topic string) {
	s.mutexLock.Lock()
	defer s.mutexLock.Unlock()

	s.topics[topic] = false
}

func (s *subscriber) ReceiveMessage(topic string, message string) {
	s.mutexLock.Lock()
	defer s.mutexLock.Unlock()

	if s.active && s.topics[topic] {
		s.messageStream <- message
	}
}

func (s *subscriber) SelfDestruct() {

}

// Listens to the message channel, prints once received.
func (s *subscriber) Listen() {
	go func() {
		for {
			if msg, ok := <-s.messageStream; ok {
				fmt.Printf("\nIncoming message for %s: %s\n", s.name, msg)
			}
		}
	}()
}
