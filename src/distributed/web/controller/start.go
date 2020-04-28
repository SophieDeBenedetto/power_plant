package controller

import (
	"net/http"

	"github.com/sophiedebenedetto/power_plant/src/distributed/messaging"
)

// Initialize registers the routes and file servers
func Initialize(rabbitServer *messaging.Server) {
	registerRoutes(rabbitServer)
	registerFileServers()
}

func registerRoutes(rabbitServer *messaging.Server) {
	ws := NewWebsocketController(rabbitServer)
	http.HandleFunc("/ws", ws.handleMessage)
}

func registerFileServers() {
	http.Handle("/public/",
		http.FileServer(http.Dir("assets")))
	http.Handle("/public/lib",
		http.StripPrefix("/public/lib/",
			http.FileServer(http.Dir("node_modules"))))
}
