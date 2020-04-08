package coordinator

import (
	"bytes"
	"encoding/gob"

	"github.com/sophiedebenedetto/power_plant/src/distributed/dto"
	"github.com/streadway/amqp"
)

// Handler knows how to decode messages
type Handler struct {
}

// Handle decodes messages
func (h *Handler) Handle(msg amqp.Delivery) *dto.SensorMessage {
	reader := bytes.NewReader(msg.Body)
	decoder := gob.NewDecoder(reader)
	sensorData := new(dto.SensorMessage)
	decoder.Decode(sensorData)
	return sensorData
}
