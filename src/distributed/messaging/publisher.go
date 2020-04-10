package messaging

import (
	"github.com/streadway/amqp"
)

// Publisher publishes to the queue
type Publisher struct {
	Server       *Server
	Channel      *amqp.Channel
	Queue        amqp.Queue
	QueueName    string
	ExchangeName string
}

// NewPublisherWithQueue returns a new publisher with a declared queue
func NewPublisherWithQueue(s *Server, queue string, autoDelete bool) *Publisher {
	ch, err := s.Conn.Channel()
	FailOnError(err, "Failed to open channel")
	q, err := ch.QueueDeclare(queue, false, autoDelete, false, false, nil)
	FailOnError(err, "Failed to declare and connect to queue")
	return &Publisher{
		Server:       s,
		Queue:        q,
		QueueName:    q.Name,
		ExchangeName: "",
		Channel:      ch,
	}
}

// NewPublisherWithExchange returns a publisher struct with the server and queue
func NewPublisherWithExchange(s *Server, exchange string) *Publisher {
	ch, err := s.Conn.Channel()
	FailOnError(err, "Failed to open channel")
	return &Publisher{
		Server:       s,
		Channel:      ch,
		QueueName:    "",
		ExchangeName: exchange,
	}
}

// Message returns a message to be published
func (p *Publisher) Message(contentType string, body []byte) amqp.Publishing {
	return amqp.Publishing{
		ContentType: contentType,
		Body:        body,
	}
}

// Publish publishes a message
func (p *Publisher) Publish(msg amqp.Publishing) {
	p.Channel.Publish(p.ExchangeName, p.QueueName, false, false, msg)
}

// Stop closes the channel connection
func (p *Publisher) Stop() {
	p.Channel.Close()
}
