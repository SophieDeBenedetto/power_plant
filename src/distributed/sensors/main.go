package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"log"

	// "math/rand"
	"strconv"
	"time"

	"github.com/sophiedebenedetto/power_plant/src/distributed/dto"
	"github.com/sophiedebenedetto/power_plant/src/distributed/messaging"
	"github.com/sophiedebenedetto/power_plant/src/distributed/sensors/calculations"
)

var name = flag.String("name", "sensor", "name of sensor")
var freq = flag.Uint("freq", 5, "update frequency in cycle/seconds")
var max = flag.Float64("max", 5.0, "maximum value for generated readings")
var min = flag.Float64("min", 1.0, "minimum value for generated readings")
var stepSize = flag.Float64("step", 0.1, "maximum allowage change per measurement")

// SensorList the name of the sensor list queue
var SensorList = "sensor_list"

func main() {
	flag.Parse()

	rabbitServer := messaging.NewRabbitMQServer("guest", "guest", "localhost:5672")
	rabbitServer.Connect()
	defer rabbitServer.Close()

	publishNewSensorQueue(rabbitServer)

	sensorPublisher := messaging.NewPublisher(rabbitServer)
	defer sensorPublisher.Stop()

	pace := setPace(freq)
	buf, enc := setUpBuffer()

	var value float64
	for range pace {
		value = calculations.Calculate(max, min, stepSize, value)
		writeMessageToBuffer(value, buf, enc)
		amqpMsg := sensorPublisher.Message("text/plain", buf.Bytes())
		sensorPublisher.Publish(amqpMsg)
		log.Printf("Reading sent. Value: %v", value)
	}
}

func publishNewSensorQueue(rabbitServer *messaging.Server) {
	publisher := messaging.NewPublisherWithQueue(rabbitServer, SensorList)
	defer publisher.Stop()
	fmt.Print("Publishg message: ", *name)
	msg := publisher.Message("text/plan", []byte(*name))
	publisher.Publish(msg)
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

func writeMessageToBuffer(value float64, buf *bytes.Buffer, enc *gob.Encoder) {
	reading := dto.SensorMessage{
		Name:      *name,
		Value:     value,
		Timestamp: time.Now(),
	}
	buf.Reset()
	enc.Encode(reading)
}
