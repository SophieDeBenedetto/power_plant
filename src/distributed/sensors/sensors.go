package sensors

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/sophiedebenedetto/power_plant/src/distributed/dto"
	"github.com/sophiedebenedetto/power_plant/src/distributed/messaging"
	"github.com/sophiedebenedetto/power_plant/src/distributed/sensors/calculations"
)

// SensorList the name of the sensor list queue
var SensorList = "sensor_list"

// SensorManager the sensor manager
type SensorManager struct {
	Server     *messaging.Server
	Calculator *calculations.Calculator
	Name       string
	Frequency  uint
}

// New creates a new sensor manager
func New(server *messaging.Server, name string, max *float64, min *float64, stepSize *float64, freq uint) *SensorManager {
	return &SensorManager{
		Server:     server,
		Calculator: calculations.New(max, min, stepSize),
		Frequency:  freq,
	}
}

// PublishNewSensorQueue publishes the 'new sensor online' message
func (manager *SensorManager) PublishNewSensorQueue() {
	publisher := messaging.NewPublisherWithQueue(manager.Server, SensorList)
	defer publisher.Stop()
	fmt.Print("Publishg message: ", manager.Name)
	msg := publisher.Message("text/plan", []byte(manager.Name))
	publisher.Publish(msg)
}

// Run takes sensor readings and publishes them
func (manager *SensorManager) Run() {
	sensorPublisher := messaging.NewPublisherWithExchange(manager.Server, "amq.fanout")
	defer sensorPublisher.Stop()

	pace := setPace(&manager.Frequency)
	buf, enc := setUpBuffer()

	var value float64
	for range pace {
		value = manager.Calculator.Calculate(value)
		writeMessageToBuffer(manager.Name, value, buf, enc)
		amqpMsg := sensorPublisher.Message("text/plain", buf.Bytes())
		sensorPublisher.Publish(amqpMsg)
		log.Printf("Reading sent. Value: %v", value)
	}
}

func setPace(freq *uint) <-chan time.Time {
	dur, _ := time.ParseDuration(strconv.Itoa(1000/int(*freq)) + "ms")
	return time.Tick(dur)
}

func setUpBuffer() (*bytes.Buffer, *gob.Encoder) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	return buf, enc
}

func writeMessageToBuffer(name string, value float64, buf *bytes.Buffer, enc *gob.Encoder) {
	reading := dto.SensorMessage{
		Name:      name,
		Value:     value,
		Timestamp: time.Now(),
	}
	buf.Reset()
	enc = gob.NewEncoder(buf)
	enc.Encode(reading)
}