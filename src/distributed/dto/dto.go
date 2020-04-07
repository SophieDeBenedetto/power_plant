package dto

import (
	"encoding/gob"
	"time"
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
