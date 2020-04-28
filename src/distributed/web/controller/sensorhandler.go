package controller

import (
	"github.com/sophiedebenedetto/power_plant/src/distributed/web/model"
	"github.com/streadway/amqp"
)

// SensorHandler handles messages from the sensors web exchange
type SensorHandler struct {
	wsc *WebsocketController
}

// NewSensorHandler returns a new SensorHandler
func NewSensorHandler(wsc *WebsocketController) *SensorHandler {
	return &SensorHandler{
		wsc: wsc,
	}
}

// Handle handles messages
func (h *SensorHandler) Handle(msg amqp.Delivery) {
	sensor, _ := model.GetSensorByName(string(msg.Body))
	h.wsc.SendMessage(message{
		Type: "source",
		Data: sensor,
	})
}
