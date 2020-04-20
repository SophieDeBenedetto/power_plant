package main

import (
	"github.com/sophiedebenedetto/power_plant/src/distributed/coordinator"
	"github.com/sophiedebenedetto/power_plant/src/distributed/datamanager"
	"github.com/sophiedebenedetto/power_plant/src/distributed/messaging"
)

func main() {
	rabbitServer := messaging.NewRabbitMQServer("guest", "guest", "localhost:5672")
	rabbitServer.Connect()
	defer rabbitServer.Close()

	consumer := messaging.NewConsumer(rabbitServer, coordinator.PersistenceQueue, true, datamanager.NewPersistenceHandler())
	consumer.Consume(false, true)
}
