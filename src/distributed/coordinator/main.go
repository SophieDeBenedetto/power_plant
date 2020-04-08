package main

import (
	"fmt"

	"github.com/sophiedebenedetto/power_plant/src/distributed/coordinator/coordinator"
	"github.com/sophiedebenedetto/power_plant/src/distributed/messaging"
)

func main() {
	rabbitServer := messaging.NewRabbitMQServer("guest", "guest", "localhost:5672")
	rabbitServer.Connect()
	newSensorConsumer := messaging.NewConsumer(rabbitServer, "sensor_list")
	msgs := newSensorConsumer.Consume()

	var sensorQueueName string
	coord := coordinator.New()
	for msg := range msgs {
		sensorQueueName = string(msg.Body)
		fmt.Println("Received message on sensor list queue")
		fmt.Println(sensorQueueName)
		if !coord.QueueIsRegistered(sensorQueueName) {
			sensorConsumer := messaging.NewConsumer(rabbitServer, sensorQueueName)
			sensorConsumer.QueueBind("", "amq.fanout")
			coord.RegisterQueue(sensorQueueName, sensorConsumer)
			sensorMessages := sensorConsumer.Consume()
			go coord.HandleMessages(sensorMessages)
		}
	}
}
