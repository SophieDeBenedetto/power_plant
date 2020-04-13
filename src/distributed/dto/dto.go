package dto

import (
	"encoding/gob"
	"time"

	"github.com/sophiedebenedetto/power_plant/src/distributed/messaging"
)

// SensorMessage is a message from a sensor
type SensorMessage struct {
	Name      string
	Value     float64
	Timestamp time.Time
}

func init() {
	gob.Register(SensorMessage{})
}

// Encode codes the message
func (message SensorMessage) Encode(writer *messaging.Writer) {
	writer.Reset()
	writer.Enc.Encode(message)
}
