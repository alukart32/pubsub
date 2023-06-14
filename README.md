# Pub/Sub service

This project implements a Pub/Sub messaging pattern using in-process communication between several goroutines over a channel.

## Model

The following entities are defined to send and receive messages on certain topics:

- Broker
- Subscriber
- Message

**Broker** is the central link in the transmission of messages from topics to subscribers.
He keeps track of all his subscribers and topics with their subscribers.

**Subscriber** knows which topics he needs to listen to and how to receive messages from them.
In this implementation, the subscriber outputs the received message to the default log.

**Message** is a container for transmitting data from a specific topic.

## Example

To run the example, you need to execute

    go run ./test/main.go

or run bin

    ./test/main (./test/main.exe)
