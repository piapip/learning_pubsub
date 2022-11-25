package pubsub

import (
	"fmt"
	"sync"
)

type publisher struct {
	subscribers map[string]PubsubSubscriber
	topics      map[string]map[string]PubsubSubscriber
	mutexLock   sync.RWMutex
}

func NewPublisher() PubsubPublisher {
	return &publisher{
		subscribers: make(map[string]PubsubSubscriber),
		topics:      make(map[string]map[string]PubsubSubscriber),
	}
}

func (p *publisher) AddSubscriber(subscriber PubsubSubscriber) error {
	p.mutexLock.Lock()
	defer p.mutexLock.Unlock()

	name, ok := subscriber.GetName()
	if !ok {
		return fmt.Errorf("the subscriber %s is disable", name)
	}

	p.subscribers[name] = subscriber

	return nil
}

func (p *publisher) AddTopic(topic string) {
	p.mutexLock.Lock()
	defer p.mutexLock.Unlock()

	if p.topics[topic] == nil {
		p.topics[topic] = make(map[string]PubsubSubscriber)
	}
}

func (p *publisher) Subscribe(name string, topic string) error {
	p.mutexLock.Lock()
	defer p.mutexLock.Unlock()

	if p.topics[topic] == nil {
		p.topics[topic] = make(map[string]PubsubSubscriber)
	}

	subscriber := p.subscribers[name]

	if p.subscribers[name] == nil {
		return fmt.Errorf("subscriber either doesn't exist or isn't active")
	}

	subscriber.AddTopic(topic)

	p.topics[topic][name] = subscriber

	return nil
}

func (p *publisher) Unsubscribe(name string, topic string) error {
	p.mutexLock.Lock()
	defer p.mutexLock.Unlock()

	for _, subscriber := range p.topics[topic] {
		checkName, exist := subscriber.GetName()
		if exist {
			if name == checkName {
				subscriber.RemoveTopic(topic)
				delete(p.topics[topic], checkName)

				return nil
			}
		}
	}

	return fmt.Errorf("no subscriber with such name or topic exists")
}

func (p *publisher) Publish(topic, msg string) {
	p.mutexLock.RLock()
	subscribers := p.topics[topic]
	p.mutexLock.RUnlock()

	for _, subscriber := range subscribers {
		subscriber.ReceiveMessage(topic, msg)
	}
}

func (p *publisher) ListSubscribers(topic string) []string {
	p.mutexLock.RLock()
	defer p.mutexLock.RUnlock()

	result := make([]string, 0)

	for _, subscriber := range p.subscribers {
		name, ok := subscriber.GetName()
		if ok {
			result = append(result, name)
		}
	}

	return result
}

func (p *publisher) GetTopics() []string {
	p.mutexLock.RLock()
	defer p.mutexLock.RUnlock()

	result := make([]string, 0)

	for topic, _ := range p.topics {
		result = append(result, topic)
	}

	return result
}
