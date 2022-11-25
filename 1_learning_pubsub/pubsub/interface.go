package pubsub

type PubSub interface {
	Subscribe(s *Subscriber, topic string)
	Unsubscribe(s *Subscriber, topic string)
	Publish(topic, msg string)
	CountSubscribers(topic string) int
}
