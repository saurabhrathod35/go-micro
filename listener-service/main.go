package main

import (
	"listener-service/event"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// try to connect Rabbitmq
	rabbitmqConn, err := connect()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
		os.Exit(1)
	}
	defer rabbitmqConn.Close()

	//start listning for messages
	log.Println("Starting to listen for messages...")

	// create consumer
	consumer, err := event.NewConsumer(rabbitmqConn, "")
	if err != nil {
		panic(err)
	}
	// watch the queue and consume event
	err = consumer.Listen([]string{"logs.INFO", "logs.ERROR", "logs.WARNING"})
	if err != nil {
		log.Printf("Failed to listen for messages: %s", err)
		os.Exit(1)
	}
}

func connect() (*amqp.Connection, error) {
	var count int64
	var backoff = 1 * time.Second

	var connection *amqp.Connection

	// dont continue until connection is established
	for {
		c, err := amqp.Dial(os.Getenv("RABBITMQ_CONNECTION_STRING"))
		// c, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		if err != nil {
			log.Printf("Failed to connect to RabbitMQ: %s", err)
			count++
		} else {
			log.Println("Connected to RabbitMQ")
			connection = c
			break
		}
		if count > 5 {
			log.Println("Failed to connect to RabbitMQ after 5 attempts, exiting...", err)
			return nil, err
		}
		backoff = time.Duration(math.Pow(float64(count), 2)) * time.Second
		log.Println("Retrying connection in", backoff)
		time.Sleep(backoff)
		continue
	}
	return connection, nil
}
