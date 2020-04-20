package messaging

import (
	"fmt"

	"github.com/streadway/amqp"
)

// HandlerInterface the handler interface
type HandlerInterface interface {
	Handle(amqp.Delivery)
}

// Consumer receives messages from the queue
type Consumer struct {
	Server  *Server
	Channel *amqp.Channel
	Queue   amqp.Queue
	Handler HandlerInterface
}

// NewConsumer returns a consumer struct
func NewConsumer(s *Server, queue string, autoDelete bool, handler HandlerInterface) *Consumer {
	ch, err := s.Conn.Channel()
	FailOnError(err, "Failed to open channel")
	fmt.Println("Creating queue with name: ", queue)
	q, err := ch.QueueDeclare(queue, false, autoDelete, false, false, nil)
	FailOnError(err, "Failed to declare and connect to queue")
	return &Consumer{
		Server:  s,
		Queue:   q,
		Channel: ch,
		Handler: handler,
	}
}

// QueueBind binds the queue to an exchange
func (c *Consumer) QueueBind(key string, exchange string) {
	err := c.Channel.QueueBind(c.Queue.Name, key, exchange, false, nil)
	if err != nil {
		FailOnError(err, "Failed to bind queue to exchange")
	}
}

// Consume starts listening for messages from a queue
func (c *Consumer) Consume(autoAck bool, exclusive bool) {
	msgs, err := c.Channel.Consume(c.Queue.Name, "", autoAck, exclusive, false, false, nil)
	FailOnError(err, "Failed to start consumer")
	for msg := range msgs {
		c.Handler.Handle(msg)
	}
}

// ExchangeDeclare declares the exchange
func (c *Consumer) ExchangeDeclare(exchange string) {
	err := c.Channel.ExchangeDeclare(
		exchange, // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	FailOnError(err, "Failed to declare exchange")
}

// Stop closes the channel connection
func (c *Consumer) Stop() {
	c.Channel.Close()
}
