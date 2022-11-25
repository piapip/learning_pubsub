package pubsub

// Lesson learn:
//   - only do this double interface BS if we are going to make [multi - multi]
//   - For [1 - multi] or [multi - 1] should be done with only 1 interface
//     If we are making 2 interfaces for [multi - 1], it's wasteful.
//   - For [mutli - multi], 2 interfaces are not enough, we'll need a 3rd party.
//     Then the main client will use only the 3rd party to avoid heavy coupling.
type PubsubPublisher interface {
	Subscribe(name string, topic string) error
	Unsubscribe(name string, topic string) error
	Publish(topic, msg string)
	AddSubscriber(subscriber PubsubSubscriber) error
	ListSubscribers(topic string) []string
	GetTopics() []string
}

type PubsubSubscriber interface {
	AddTopic(topic string)
	RemoveTopic(topic string)
	ReceiveMessage(topic string, message string)
	SelfDestruct()
	Listen()
	GetTopics() []string
	GetName() (string, bool)
}
