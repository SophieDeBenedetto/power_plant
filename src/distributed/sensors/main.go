package main

import (
	"bytes"
	"encoding/gob"
	"flag"
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

func main() {
	flag.Parse()

	pace := setPace(freq)
	buf, enc := setUpBuffer()

	rabbitServer, publisher := getPublisher(*name)
	defer publisher.Stop()
	defer rabbitServer.Close()

	var value float64
	for range pace {
		value = calculations.Calculate(max, min, stepSize, value)
		writeMessageToBuffer(value, buf, enc)
		amqpMsg := publisher.Message("text/plain", buf.Bytes())
		publisher.Publish(amqpMsg)
		log.Printf("Reading sent. Value: %v", value)
	}
}

func setPace(freq *uint) <-chan time.Time {
	dur, _ := time.ParseDuration(strconv.Itoa(1000/int(*freq)) + "ms")
	return time.Tick(dur)
}

func getPublisher(queue string) (*messaging.Server, *messaging.Publisher) {
	rabbitServer := messaging.NewRabbitMQServer("guest", "guest", "localhost:5672")
	rabbitServer.Connect()

	return rabbitServer, messaging.NewPublisher(rabbitServer, queue)
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
