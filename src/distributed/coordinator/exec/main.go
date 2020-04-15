package main

import (
	"github.com/sophiedebenedetto/power_plant/src/distributed/coordinator"
	"github.com/sophiedebenedetto/power_plant/src/distributed/messaging"
)

func main() {
	rabbitServer := messaging.NewRabbitMQServer("guest", "guest", "localhost:5672")
	rabbitServer.Connect()
	defer rabbitServer.Close()

	messaging.NewConsumer(rabbitServer, coordinator.PersistenceQueue, true, coordinator.NewPersistenceHandler())

	dbEventRaiser := coordinator.NewDatabaseEventRaiser(rabbitServer)

	coord := coordinator.New(rabbitServer, dbEventRaiser.EventRaiser)
	coord.Run()
}
