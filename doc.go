/*
Package pubsub provides Pub/Sub pattern implementations.

The message broker adds new topics, registers new subscribers of
these topics and sends them messages.

Adding a subscriber to a topic:

	broker := NewBroker()
	sub := broker.AddSubscriber()
	broker.Subscribe(sub, topicName)

Each subscriber is waiting for new messages from all topics:

	go sub.Listen()

Anyone who wants to send a message through a broker for a specific topic needs
to do the following:

	msg := []byte("some data")
	topicName := "topic_1"
	broker.Publish(topicName, msg)
*/
package pubsub
