package messaging

import (
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
func NewConsumer(s *Server, queue string, handler HandlerInterface) *Consumer {
	ch, err := s.Conn.Channel()
	FailOnError(err, "Failed to open channel")
	q, err := ch.QueueDeclare(queue, false, false, false, false, nil)
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
func (c *Consumer) Consume() {
	msgs, err := c.Channel.Consume(c.Queue.Name, "", true, false, false, false, nil)
	FailOnError(err, "Failed to start consumer")
	for msg := range msgs {
		c.Handler.Handle(msg)
	}
}

// Stop closes the channel connection
func (c *Consumer) Stop() {
	c.Channel.Close()
}
