package coordinator

import (
	"fmt"

	"github.com/sophiedebenedetto/power_plant/src/distributed/dto"
	"github.com/sophiedebenedetto/power_plant/src/distributed/messaging"
	"github.com/streadway/amqp"
)

// New returns a new coordinator struct
func New() *Coordinator {
	return &Coordinator{
		SensorRegistry: make(map[string]*messaging.Consumer),
		Handler:        &Handler{},
	}
}

// Coordinator knows the sensor registry
type Coordinator struct {
	SensorRegistry map[string]*messaging.Consumer
	Handler        *Handler
}

// QueueIsRegistered returns true if the queue is registered
func (c *Coordinator) QueueIsRegistered(queue string) bool {
	if c.SensorRegistry[queue] != nil {
		return true
	}
	return false
}

// RegisterQueue registers the queue
func (c *Coordinator) RegisterQueue(queue string, consumer *messaging.Consumer) {
	c.SensorRegistry[queue] = consumer
}

// HandleMessages processes the messages coming over the channel
func (c *Coordinator) HandleMessages(msgs <-chan amqp.Delivery) {
	var sensorData *dto.SensorMessage
	for msg := range msgs {
		sensorData = c.Handler.Handle(msg)
		fmt.Printf("Received message, Name: %v, Value: %v, Timestamp: %v \n", sensorData.Name, sensorData.Value, sensorData.Timestamp)
	}
}
