package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "80"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	rabbitmqConn, err := connect()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
		os.Exit(1)
	}
	defer rabbitmqConn.Close()
	app := Config{
		Rabbit: rabbitmqConn,
	}

	log.Printf("Starting broker service on port %s\n", webPort)

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start the server
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
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
