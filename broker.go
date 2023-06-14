package pubsub

import (
	"log"
	"sync"
)

type Subscribers map[string]*Subscriber

// Broker represents the message broker. A new topic is added together with
// its first subscriber. The message is published to all subscribers, otherwise
// the message is lost.
type Broker struct {
	subscribers Subscribers
	topics      map[string]Subscribers
	mtx         *sync.RWMutex
}

// NewBroker creates a new Broker.
func NewBroker() *Broker {
	return &Broker{
		subscribers: make(Subscribers),
		topics:      make(map[string]Subscribers),
		mtx:         new(sync.RWMutex),
	}
}

// AddSubscriber creates and adds a new broker subscriber. The limit
// for receiving messages can be specified via msgThreshold.
func (b *Broker) AddSubscriber(msgThreshold int) *Subscriber {
	if msgThreshold < 0 {
		msgThreshold = 0
	}

	b.mtx.Lock()
	defer b.mtx.Unlock()
	s := NewSubscriber(msgThreshold)
	b.subscribers[s.id] = s
	return s
}

// Subscribe adds a new subscriber to the topic.
func (b *Broker) Subscribe(s *Subscriber, topic string) {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	if b.topics[topic] == nil {
		b.topics[topic] = make(Subscribers)
	}

	s.AddTopic(topic)
	b.topics[topic][s.id] = s
	log.Printf("%s subscribed to %s\n", s.id, topic)
}

// Unsubscribe removes the subscriber from the topic.
func (b *Broker) Unsubscribe(s *Subscriber, topic string) {
	b.mtx.RLock()
	defer b.mtx.RUnlock()

	delete(b.topics[topic], s.id)
	s.RemoveTopic(topic)
	log.Printf("%s unsubscribed from %s\n", s.id, topic)
}

// Len counts the current number of subscribers of the topic.
func (b *Broker) Len(topic string) int {
	b.mtx.RLock()
	defer b.mtx.RUnlock()
	return len(b.topics[topic])
}

// RemoveSubscriber removes the subscriber from the broker. It unsubscribes him
// from all topics.
func (b *Broker) RemoveSubscriber(s *Subscriber) {
	for topic := range s.topics {
		b.Unsubscribe(s, topic)
	}
	b.mtx.Lock()
	delete(b.subscribers, s.id)
	b.mtx.Unlock()
	s.Done()
}

// Publish sends a message to all subscribers of the topic.
func (b *Broker) Publish(topic string, msg []byte) error {
	topicMsg, err := NewMessage(topic, msg)
	if err != nil {
		return err
	}

	b.mtx.RLock()
	subs := b.topics[topic]
	b.mtx.RUnlock()
	for _, s := range subs {
		go s.Receive(topicMsg)
	}
	return nil
}
