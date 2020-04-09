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
	coord *Coordinator
}

// SensorDataHandler knows how to decode messages
type SensorDataHandler struct {
}

// Handle decodes messages
func (h *SensorListMessageHandler) Handle(msg amqp.Delivery) {
	sensorQueueName := string(msg.Body)
	fmt.Println("Received message on sensor list queue")
	fmt.Println(sensorQueueName)
	if !h.coord.QueueIsRegistered(sensorQueueName) {
		sensorConsumer := messaging.NewConsumer(h.coord.Server, sensorQueueName, &SensorDataHandler{})
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
}
