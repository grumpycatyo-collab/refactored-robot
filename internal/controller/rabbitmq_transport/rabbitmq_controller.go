package rabbitmq_transport

import (
	"log"

	amqp "github.com/streadway/amqp"
)

type AMQController struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewAMQController(amqpURL string) (*AMQController, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &AMQController{
		conn:    conn,
		channel: channel,
	}, nil

}

func (c *AMQController) PublishMessage(exchange, routingKey, message string) error {
	err := c.channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	return err
}

func (c *AMQController) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *AMQController) DeclareQueue(queueName string) error {
	_, err := c.channel.QueueDeclare(
		queueName, // Queue name
		true,      // Durable
		false,     // Delete when unused
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	return err
}

func (c *AMQController) CreateQueueAndPublishMessage(queueName, message string) error {
	err := c.DeclareQueue(queueName)
	if err != nil {
		return err
	}

	exchangeName := ""
	routingKey := queueName
	err = c.PublishMessage(exchangeName, routingKey, message)
	if err != nil {
		return err
	}

	return nil
}

func (c *AMQController) Receiver(name string) {
	msgs, err := c.channel.Consume(
		name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Waiting for messages...")
	for msg := range msgs {
		log.Printf("Received message: %s", msg.Body)
	}
}
