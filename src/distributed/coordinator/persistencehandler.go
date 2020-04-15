package coordinator

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/sophiedebenedetto/power_plant/src/distributed/dto"
	"github.com/streadway/amqp"
)

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

	fmt.Printf("PERSISTING MESSAGE: %v\n", sensorData)
}
