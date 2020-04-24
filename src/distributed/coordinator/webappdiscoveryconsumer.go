package coordinator

import (
	"github.com/sophiedebenedetto/power_plant/src/distributed/messaging"
	"github.com/streadway/amqp"
)

//WebAppDiscoveryQueue the queue over which the web app requests sensor lists
var WebAppDiscoveryQueue = "WebAppDiscovery"

// WebAppDiscoveryConsumer listens for discovery requests from the web app
type WebAppDiscoveryConsumer struct {
	server   *messaging.Server
	consumer *messaging.Consumer
}

type discoveryHandler struct {
	eventRaiser EventRaiser
}

func (handler *discoveryHandler) Handle(msg amqp.Delivery) {
	handler.eventRaiser.PublishEvent("WebAppDiscoveryRequest", string(msg.Body))
}

// NewWebAppDiscoveryConsumer returns the new WebAppDiscoveryConsumer
func NewWebAppDiscoveryConsumer(server *messaging.Server, eventRaiser EventRaiser) *WebAppDiscoveryConsumer {
	handler := &discoveryHandler{
		eventRaiser: eventRaiser,
	}
	return &WebAppDiscoveryConsumer{
		server:   server,
		consumer: messaging.NewConsumer(server, WebAppDiscoveryQueue, false, handler),
	}
}

// Run runs the consumer
func (webConsumer *WebAppDiscoveryConsumer) Run() {
	webConsumer.consumer.Consume(true, false)
}
