package coordinator

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/sophiedebenedetto/power_plant/src/distributed/dto"
	"github.com/sophiedebenedetto/power_plant/src/distributed/messaging"
	"github.com/streadway/amqp"
)

// SensorListMessageHandler knows how to decode messages
type SensorListMessageHandler struct {
	coord           *Coordinator
	eventaggregator *EventAggregator
}

// SensorDataHandler knows how to decode messages
type SensorDataHandler struct {
	eventaggregator *EventAggregator
}

// Handle decodes messages
func (h *SensorListMessageHandler) Handle(msg amqp.Delivery) {
	sensorQueueName := string(msg.Body)
	fmt.Println("Received message on sensor_list queue")
	if !h.coord.QueueIsRegistered(sensorQueueName) {
		h.eventaggregator.PublishEvent("DataSourceDiscovered", sensorQueueName)
		sensorDataHandler := &SensorDataHandler{
			eventaggregator: NewEventAggregator(),
		}
		sensorConsumer := messaging.NewConsumer(h.coord.Server, sensorQueueName, false, sensorDataHandler)
		sensorConsumer.QueueBind("", "amq.fanout")
		h.coord.RegisterQueue(sensorQueueName, sensorConsumer)
		go sensorConsumer.Consume()
	}
}

// Handle decodes messages
func (h *SensorDataHandler) Handle(msg amqp.Delivery) {
	reader := bytes.NewReader(msg.Body)
	decoder := gob.NewDecoder(reader)
	sensorData := new(dto.SensorMessage)
	err := decoder.Decode(sensorData)
	if err != nil {
		fmt.Println(fmt.Errorf("Error decoding: %v", err))
	}

	fmt.Printf("Received message: %v\n", sensorData)

	eventData := EventData{
		Name:      sensorData.Name,
		Value:     sensorData.Value,
		Timestamp: sensorData.Timestamp,
	}
	h.eventaggregator.PublishEvent("MessageReceived_"+msg.RoutingKey, eventData)
}
