package controller

import (
	"bytes"
	"encoding/gob"

	"github.com/sophiedebenedetto/power_plant/src/distributed/dto"
	"github.com/streadway/amqp"
)

// SensorReadingHandler handles messages from the sensors web exchange
type SensorReadingHandler struct {
	wsc *WebsocketController
}

// NewSensorReadingHandler returns a new SensorHandler
func NewSensorReadingHandler(wsc *WebsocketController) *SensorReadingHandler {
	return &SensorReadingHandler{
		wsc: wsc,
	}
}

// Handle handles messages
func (h *SensorReadingHandler) Handle(msg amqp.Delivery) {
	buf := bytes.NewBuffer(msg.Body)
	dec := gob.NewDecoder(buf)
	sm := dto.SensorMessage{}
	dec.Decode(&sm)
	h.wsc.SendMessage(message{
		Type: "reading",
		Data: sm,
	})
}
