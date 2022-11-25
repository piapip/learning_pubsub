package pubsub

import (
	"log"
	"sync"
)

type publisher struct {
	subcribers Subscribers            // map of subscribers id:subscribers ??? (shouldn't it be map then?)
	topics     map[string]Subscribers // map of topic to subscribers topic
	mutextLock sync.RWMutex           // mutext lock
}

func NewPubsub() PubSub {
	return &publisher{
		subcribers: Subscribers{},
		topics:     make(map[string]Subscribers),
	}
}

func (p *publisher) Subscribe(s *Subscriber, topic string) {
	p.mutextLock.RLock()
	defer p.mutextLock.RUnlock()

	if p.topics[topic] == nil {
		p.topics[topic] = Subscribers{}
	}

	if p.subcribers[s.id] == nil {
		p.subcribers[s.id] = s
	}

	s.AddTopic(topic)

	p.topics[topic][s.id] = s

	log.Printf("%s subscribed for topic: %s\n", s.id, topic)
}

func sendSignal(s *Subscriber, msg *Message) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if s.active {
		s.message <- msg
	}
}

func (p *publisher) Unsubscribe(s *Subscriber, topic string) {
	p.mutextLock.RLock()
	defer p.mutextLock.RUnlock()

	// if there's no topic, then stop
	if p.topics[topic] == nil || len(p.topics[topic]) == 0 {
		return
	}

	// if the topic doesn't relate to the subscriber, then stop
	if p.topics[topic][s.id] == nil {
		return
	}

	s.RemoveTopic(topic)
	delete(p.topics[topic], s.id)

	log.Printf("%s unsubscribed for topic: %s\n", s.id, topic)
}

func (p *publisher) Publish(topic, msg string) {
	p.mutextLock.RLock()
	pSubscribers := p.topics[topic]
	p.mutextLock.RUnlock()
	newMsg := NewMessage(topic, msg)
	for _, s := range pSubscribers {
		if !s.active {
			return
		}

		go (func() {
			sendSignal(s, newMsg)
		})()
	}
}

func (p *publisher) CountSubscribers(topic string) int {
	p.mutextLock.RLock()
	defer p.mutextLock.RUnlock()

	return len(p.topics[topic])
}
