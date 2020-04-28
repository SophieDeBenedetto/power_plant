package webconsumer

import (
	"github.com/sophiedebenedetto/power_plant/src/distributed/messaging"
)

// WebConsumer consumes web app exchanges
type WebConsumer struct {
	Server   *messaging.Server
	Consumer *messaging.Consumer
}

// NewWebConsumer returns the new consumer
func NewWebConsumer(server *messaging.Server, queue string, exchange string, handler messaging.HandlerInterface) *WebConsumer {
	consumer := messaging.NewConsumer(server, queue, true, handler)
	consumer.ExchangeDeclare(exchange)
	consumer.QueueBind("", exchange)
	return &WebConsumer{
		Server:   server,
		Consumer: consumer,
	}
}

// Run runs the consumer
func (c *WebConsumer) Run() {
	c.Consumer.Consume(true, false)
}
