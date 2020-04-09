package main

import (
	"flag"

	// "math/rand"

	"github.com/sophiedebenedetto/power_plant/src/distributed/messaging"
	"github.com/sophiedebenedetto/power_plant/src/distributed/sensors"
)

var name = flag.String("name", "sensor", "name of sensor")
var freq = flag.Uint("freq", 5, "update frequency in cycle/seconds")
var max = flag.Float64("max", 5.0, "maximum value for generated readings")
var min = flag.Float64("min", 1.0, "minimum value for generated readings")
var stepSize = flag.Float64("step", 0.1, "maximum allowage change per measurement")

func main() {
	flag.Parse()

	rabbitServer := messaging.NewRabbitMQServer("guest", "guest", "localhost:5672")
	rabbitServer.Connect()
	defer rabbitServer.Close()

	manager := sensors.New(rabbitServer, *name, min, max, stepSize, *freq)
	manager.PublishNewSensorQueue()
	manager.Run()
}
