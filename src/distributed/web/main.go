package main

import (
	"net/http"

	"github.com/sophiedebenedetto/power_plant/src/distributed/messaging"
	"github.com/sophiedebenedetto/power_plant/src/distributed/web/controller"
)

func main() {
	server := messaging.NewRabbitMQServer("guest", "guest", "localhost:5672")
	server.Connect()
	defer server.Close()
	controller.Initialize(server)
	http.ListenAndServe(":3000", nil)
}
