package event

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *amqp.Connection
	QueueName string
}

func NewConsumer(conn *amqp.Connection, queueName string) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}
	err := consumer.Setup()
	if err != nil {
		log.Printf("\nFailed to setup consumer: %s", err)
		return Consumer{}, err
	}
	return consumer, nil
}

func (connection *Consumer) Setup() error {
	channel, err := connection.conn.Channel()
	if err != nil {
		log.Printf("\nFailed to open a channel: %s", err)
		return err
	}
	return declareExchange(channel)
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		log.Printf("\nFailed to open a channel: %s", err)
		return err
	}
	defer ch.Close()
	queue, err := declareRandomQueue(ch)
	if err != nil {
		log.Printf("\nFailed to declare queue: %s", err)
		return err
	}
	for _, topic := range topics {
		err = ch.QueueBind(
			queue.Name,   // queue name
			topic,        // routing key
			"logs_topic", // exchange name
			false,        // no wait
			nil,          // arguments
		)
		if err != nil {
			log.Printf("\nFailed to bind queue: %s", err)
			return err
		}

	}
	msgs, err := ch.Consume(
		queue.Name, // queue name
		"",         // consumer tag
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		log.Printf("\nFailed to register a consumer: %s", err)
		return err
	}
	forever := make(chan bool)
	go func() {
		log.Printf("\nWaiting for messages on queue: %s", queue.Name)
		for d := range msgs {
			var payload Payload
			err = json.Unmarshal(d.Body, &payload)
			if err != nil {
				log.Printf("\nFailed to unmarshal message: %s", err)
				continue
			}
			go handlePayload(payload)
		}
	}()
	log.Printf("\nConsumer is listening for messages on queue: %s", queue.Name)
	<-forever
	return nil
}

// func (consumer *Consumer) Listen(topics []string) error {
// 	ch, err := consumer.conn.Channel()
// 	if err != nil {
// 		return err
// 	}
// 	defer ch.Close()

// 	q, err := declareRandomQueue(ch)
// 	if err != nil {
// 		return err
// 	}

// 	for _, s := range topics {
// 		ch.QueueBind(
// 			q.Name,
// 			s,
// 			"logs_topic",
// 			false,
// 			nil,
// 		)

// 		if err != nil {
// 			return err
// 		}
// 	}

// 	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
// 	if err != nil {
// 		return err
// 	}

// 	forever := make(chan bool)
// 	go func() {
// 		for d := range messages {
// 			var payload Payload
// 			_ = json.Unmarshal(d.Body, &payload)

// 			go handlePayload(payload)
// 		}
// 	}()

// 	fmt.Printf("Waiting for message [Exchange, Queue] [logs_topic, %s]\n", q.Name)
// 	<-forever

// 	return nil
// }

func handlePayload(payload Payload) {
	log.Printf("\nHandling payload: %s", payload.Name)
	log.Printf("\nReceived message: %s", payload.Data)
	switch payload.Name {
	case "log", "event":
		err := logEvent(payload)
		if err != nil {
			log.Printf("\nFailed to log event: %s", err)
		}
	case "auth":
		log.Printf("\nReceived auth event: %s", payload.Data)
	default:
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
		log.Printf("\nReceived unknown event type: %s", payload.Name)
	}
}

func logEvent(payload Payload) error {
	jsonData, _ := json.MarshalIndent(payload, "", "\t")
	logServiceURL := "http://logger-service/log"
	req, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("\nFailed to create request: %s", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Printf("\nFailed to send request: %s", err)
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusAccepted {
		log.Printf("\nReceived non-accepted status code: %d", response.StatusCode)
		return errors.New("error calling log service")
	}

	return nil
}
