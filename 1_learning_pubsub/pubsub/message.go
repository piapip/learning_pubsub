package pubsub

type Message struct {
	topic string
	body  string
}

func NewMessage(topic, msg string) *Message {
	return &Message{
		topic: topic,
		body:  msg,
	}
}

func (m *Message) GetTopic() string {
	return m.topic
}
func (m *Message) GetBody() string {
	return m.body
}
