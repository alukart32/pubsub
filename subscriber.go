package pubsub

import (
	"log"
	"sync"

	"github.com/google/uuid"
)

// Subscriber represents the one who expects to receive messages from the broker topic.
type Subscriber struct {
	id     string
	msgs   chan Message
	topics map[string]bool
	mtx    *sync.RWMutex
	done   chan struct{}
}

// NewSubscriber creates a new subscriber. To limit the number of messages received,
// msgThreshold must be greater than 0.
func NewSubscriber(msgThreshold int) *Subscriber {
	if msgThreshold < 0 {
		msgThreshold = 0
	}
	return &Subscriber{
		id:     uuid.New().String(),
		msgs:   make(chan Message, msgThreshold),
		topics: make(map[string]bool),
		mtx:    new(sync.RWMutex),
		done:   make(chan struct{}, 1),
	}
}

// AddTopic adds a new topic.
func (s *Subscriber) AddTopic(topic string) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.topics[topic] = true
}

// RemoveTopic removes the topic.
func (s *Subscriber) RemoveTopic(topic string) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	delete(s.topics, topic)
}

// Topics lists all the subscribers' topics.
func (s *Subscriber) Topics() []string {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	topics := make([]string, len(s.topics))
	var i int
	for k := range s.topics {
		topics[i] = k
		i++
	}
	return topics
}

// Receive proxies a message to subscriber handler.
func (s *Subscriber) Receive(msg Message) {
	select {
	case <-s.done:
		log.Printf("%s can't receive msg, he is done\n", s.id)
	case s.msgs <- msg:
	}
}

// Listen handles subscribers' messages.
func (s *Subscriber) Listen() {
	for msg := range s.msgs {
		log.Printf("%v receive msg: %s from topic: %s\n", s.id, msg.Body(), msg.Topic())
	}
}

// Done prepares the subscriber to finish receiving messages.
func (s *Subscriber) Done() {
	s.done <- struct{}{}
	close(s.msgs)
}
