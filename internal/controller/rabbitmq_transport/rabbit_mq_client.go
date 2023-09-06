package rabbitmq_transport

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func StartMessageConsumer(name string) {
	conn, err := amqp.Dial("amqp://guest:rmpassword@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	} else {
		log.Printf("Succesfully connected to rabbit mq")
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	queueName := fmt.Sprintf("user_%s", name)
	q, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Printf("Waiting for messages...")

	for msg := range msgs {
		log.Printf("Received a message: %s", msg.Body)
	}
}
