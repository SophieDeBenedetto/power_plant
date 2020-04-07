package sensors

import (
	"bytes"
	"encoding/gob"
	"flag"
	"log"

	// "math/rand"
	"strconv"
	"time"

	"github.com/sophiedebenedetto/power_plant/distributed/dto"
	"github.com/sophiedebenedetto/power_plant/distributed/messaging"
	"github.com/sophiedebenedetto/power_plant/distributed/sensors/calculations"
)

var name = flag.String("name", "sensor", "name of sensor")
var freq = flag.Uint("freq", 5, "update frequency in cycle/seconds")
var max = flag.Float64("max", 5.0, "maximum value for generated readings")
var min = flag.Float64("min", 1.0, "minimum value for generated readings")
var stepSize = flag.Float64("step", 0.1, "maximum allowage change per measurement")

func main() {
	flag.Parse()

	dur, _ := time.ParseDuration(strconv.Itoa(1000/int(*freq)) + "ms")
	signal := time.Tick(dur)

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)

	server, publisher := getPublisher(*name)
	defer publisher.Stop()
	defer server.Close()

	var value float64
	for range signal {
		value = calculations.Calculate(max, min, stepSize, value)
		reading := dto.SensorMessage{
			Name:      *name,
			Value:     value,
			Timestamp: time.Now(),
		}
		buf.Reset()
		enc.Encode(reading)
		msg := publisher.Message("text/plain", buf.Bytes())
		publisher.Publish(msg)
		log.Printf("Reading sent. Value: %v", value)
	}
}

func getPublisher(queue string) (*messaging.Server, *messaging.Publisher) {
	rabbitServer := messaging.NewRabbitMQServer("guest", "guest", "localhost:5672")
	rabbitServer.Connect()

	return rabbitServer, messaging.NewPublisher(rabbitServer, queue)
}

// func setUpBuffer() (*) {
// 	buf := new(bytes.Buffer)
// 	enc := gob.NewEncoder(buf)
// 	return buf, enc
// }
