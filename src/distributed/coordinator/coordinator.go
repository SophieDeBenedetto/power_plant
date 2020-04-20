package coordinator

import (
	"fmt"

	"github.com/sophiedebenedetto/power_plant/src/distributed/messaging"
	"github.com/sophiedebenedetto/power_plant/src/distributed/sensors"
)

// Coordinator knows the sensor registry
type Coordinator struct {
	Server         *messaging.Server
	SensorRegistry map[string]*messaging.Consumer
	EventRaiser    EventRaiser
}

// New returns a new coordinator struct
func New(server *messaging.Server, er EventRaiser) *Coordinator {
	return &Coordinator{
		Server:         server,
		SensorRegistry: make(map[string]*messaging.Consumer),
		EventRaiser:    er,
	}
}

// Run runs the coordinator
func (c *Coordinator) Run(queue string) {
	handler := &SensorListMessageHandler{
		coord:       c,
		eventRaiser: c.EventRaiser,
	}
	newSensorConsumer := messaging.NewConsumer(c.Server, queue, true, handler)
	fmt.Println("Consuming sensor_list...")
	newSensorConsumer.ExchangeDeclare(sensors.SensorList)
	newSensorConsumer.QueueBind("", sensors.SensorList)
	newSensorConsumer.Consume(true, false)
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
