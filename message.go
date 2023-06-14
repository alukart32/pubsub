package pubsub

import (
	"errors"
)

var (
	ErrUndefinedTopic = errors.New("undefined topic")
)

// Message represents a topic message.
type Message struct {
	topic string
	body  []byte
}

// NewMessage creates a new message.
func NewMessage(topic string, body []byte) (Message, error) {
	if len(topic) == 0 {
		return Message{}, ErrUndefinedTopic
	}
	return Message{topic: topic, body: body}, nil
}

// Topic returns the message topic.
func (m Message) Topic() string {
	return m.topic
}

// Body returns the message body.
func (m Message) Body() []byte {
	return m.body
}

// BodyAsString returns the message body as string.
func (m Message) BodyAsString() string {
	return string(m.body)
}
