package datamanager

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/sophiedebenedetto/power_plant/src/distributed/dto"
	"github.com/streadway/amqp"
)

// NewPersistenceHandler persists the message
func NewPersistenceHandler() *PersistenceHandler {
	return &PersistenceHandler{}
}

// PersistenceHandler handles persistence messages
type PersistenceHandler struct {
}

// Handle handles the message
func (ph *PersistenceHandler) Handle(msg amqp.Delivery) {
	reader := bytes.NewReader(msg.Body)
	decoder := gob.NewDecoder(reader)
	sensorData := new(dto.SensorMessage)
	err := decoder.Decode(sensorData)
	if err != nil {
		fmt.Println(fmt.Errorf("Error decoding: %v", err))
	}

	fmt.Printf("Persisting message: %v\n", sensorData)
	err = SaveReading(sensorData)
	if err != nil {
		fmt.Printf("Failed to save message: %v, with error: ", sensorData, err)

	} else {
		msg.Ack(false)
	}
}
