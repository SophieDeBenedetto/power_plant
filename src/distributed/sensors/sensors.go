package sensors

import (
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
		Name:       name,
		Server:     server,
		Calculator: calculations.New(max, min, stepSize),
		Frequency:  freq,
	}
}

// PublishNewSensorQueue publishes the 'new sensor online' message
func (manager *SensorManager) PublishNewSensorQueue() {
	publisher := messaging.NewPublisherWithQueue(manager.Server, SensorList, true)
	defer publisher.Stop()
	fmt.Println("Publishing message: ", manager.Name)
	publisher.Publish([]byte(manager.Name))
}

// Run takes sensor readings and publishes them
func (manager *SensorManager) Run() {
	sensorPublisher := messaging.NewPublisherWithExchange(manager.Server, "amq.fanout")
	defer sensorPublisher.Stop()

	pace := setPace(&manager.Frequency)
	sensorPublisher.SetUpWriter()

	var value float64
	for range pace {
		value = manager.Calculator.Calculate(value)
		msg := dto.SensorMessage{
			Name:      manager.Name,
			Value:     value,
			Timestamp: time.Now(),
		}
		sensorPublisher.WriteMessageToBuffer(msg)
		sensorPublisher.Publish(sensorPublisher.MessageBytes())
		log.Printf("Reading sent. Value: %v", value)
	}
}

func setPace(freq *uint) <-chan time.Time {
	dur, _ := time.ParseDuration(strconv.Itoa(1000/int(*freq)) + "ms")
	return time.Tick(dur)
}
