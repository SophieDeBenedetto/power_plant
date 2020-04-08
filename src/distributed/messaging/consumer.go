package messaging

import (
	"github.com/streadway/amqp"
)

// Consumer receives messages from the queue
type Consumer struct {
	Server  *Server
	Channel *amqp.Channel
	Queue   amqp.Queue
}

// NewConsumer returns a consumer struct
func NewConsumer(s *Server, queue string) *Consumer {
	ch, err := s.Conn.Channel()
	FailOnError(err, "Failed to open channel")
	q, err := ch.QueueDeclare(queue, false, false, false, false, nil)
	FailOnError(err, "Failed to declare and connect to queue")
	return &Consumer{
		Server:  s,
		Queue:   q,
		Channel: ch,
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
func (c *Consumer) Consume() <-chan amqp.Delivery {
	msgs, err := c.Channel.Consume(c.Queue.Name, "", true, false, false, false, nil)
	FailOnError(err, "Failed to start consumer")
	return msgs
}

// Stop closes the channel connection
func (c *Consumer) Stop() {
	c.Channel.Close()
}
