package event

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventEmitter struct {
	connection *amqp.Connection
}

func (e *EventEmitter) setUp() error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	return declareExchange(channel)
}

func (e *EventEmitter) Push(event string, severity string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		log.Printf("\nFailed to open a channel: %s", err)
	}
	defer channel.Close()
	log.Printf("Pushing event: %s with severity: %s", event, severity)
	err = channel.Publish(
		"logs_topic", // exchange name
		severity,     // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(event),
		},
	)
	if err != nil {
		log.Printf("\nFailed to publish message: %s", err)
		return err
	}
	return nil
}

func NewEventEmitter(connection *amqp.Connection) (EventEmitter, error) {
	emmiter := EventEmitter{
		connection: connection,
	}
	err := emmiter.setUp()
	if err != nil {
		log.Printf("\nFailed to setup event emitter: %s", err)
		return EventEmitter{}, err
	}
	return emmiter, nil
}
