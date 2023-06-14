package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/alukart32/pubsub"
)

var topics = map[string]string{
	"T1": "TOPIC_1",
	"T2": "TOPIC_2",
	"T3": "TOPIC_3",
	"T4": "TOPIC_4",
}

func publisher(wg *sync.WaitGroup, done chan struct{}, broker *pubsub.Broker) {
	defer wg.Done()

	timer := time.NewTimer(1 * time.Second)
	defer func() {
		if !timer.Stop() {
			<-timer.C
		}
	}()

	randomTopics := make([]string, len(topics))
	var i int
	for _, v := range topics {
		randomTopics[i] = v
		i++
	}

	for {
		select {
		case <-done:
			return
		case <-timer.C:
			topic := randomTopics[rand.Intn(len(topics))]
			msg := fmt.Sprintf("%f", rand.Float64())

			log.Printf("publish msg: %v, for %v", msg, topic)
			err := broker.Publish(topic, []byte(msg))
			if err != nil {
				log.Printf("can't publish msg: %v for %s: %v", msg, topic, err)
			}
		}
		_ = timer.Reset(time.Duration(100+rand.Intn(800)) * time.Millisecond)
	}
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Lmicroseconds)

	broker := pubsub.NewBroker()

	s1 := broker.AddSubscriber(0)
	broker.Subscribe(s1, topics["T1"])
	broker.Subscribe(s1, topics["T2"])

	s2 := broker.AddSubscriber(0)
	broker.Subscribe(s2, topics["T2"])
	broker.Subscribe(s2, topics["T3"])

	s3 := broker.AddSubscriber(0)
	broker.Subscribe(s3, topics["T3"])
	broker.Subscribe(s3, topics["T4"])

	var wg sync.WaitGroup
	wg.Add(4)

	done := make(chan struct{}, 1)
	go publisher(&wg, done, broker)
	go func() {
		defer wg.Done()
		s1.Listen()
	}()
	go func() {
		defer wg.Done()
		s2.Listen()
	}()
	go func() {
		defer wg.Done()
		s3.Listen()
	}()

	time.AfterFunc(5*time.Second, func() {
		broker.Subscribe(s2, topics["T1"])
	})
	time.AfterFunc(15*time.Second, func() {
		broker.Unsubscribe(s2, topics["T2"])
	})
	time.AfterFunc(25*time.Second, func() {
		broker.RemoveSubscriber(s2)
	})

	time.AfterFunc(30*time.Second, func() {
		broker.RemoveSubscriber(s1)
		broker.RemoveSubscriber(s3)
		done <- struct{}{}
	})
	wg.Wait()
}
