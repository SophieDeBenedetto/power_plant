package controller

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sophiedebenedetto/power_plant/src/distributed/coordinator"
	"github.com/sophiedebenedetto/power_plant/src/distributed/messaging"
)

type message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// WebsocketController the WS controller
type WebsocketController struct {
	rabbitServer *messaging.Server
	sockets      []*websocket.Conn
	mutex        sync.Mutex
	upgrader     websocket.Upgrader
}

// NewWebsocketController returns a new WS controller
func NewWebsocketController() *WebsocketController {
	server := messaging.NewRabbitMQServer("guest", "guest", "localhost:5672")
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	return &WebsocketController{
		rabbitServer: server,
		upgrader:     upgrader,
	}
}

func (wsc *WebsocketController) handleMessage(w http.ResponseWriter, r *http.Request) {
	socket, _ := wsc.upgrader.Upgrade(w, r, nil)
	wsc.addSocket(socket)
	go wsc.listenForDiscoveryRequests(socket)
}

func (wsc *WebsocketController) addSocket(socket *websocket.Conn) {
	wsc.mutex.Lock()
	wsc.sockets = append(wsc.sockets, socket)
	wsc.mutex.Unlock()
}

func (wsc *WebsocketController) removeSocket(socket *websocket.Conn) {
	wsc.mutex.Lock()
	socket.Close()
	for i := range wsc.sockets {
		if wsc.sockets[i] == socket {
			wsc.sockets = append(wsc.sockets[:i], wsc.sockets[i+1:]...)
		}
	}
	wsc.mutex.Unlock()
}

func (wsc *WebsocketController) listenForDiscoveryRequests(socket *websocket.Conn) {
	for {
		msg := message{}
		err := socket.ReadJSON(&msg)
		if err != nil {
			wsc.removeSocket(socket)
		}

		if msg.Type == "discover" {
			publisher := messaging.NewPublisherWithQueue(wsc.rabbitServer, coordinator.WebAppDiscoveryQueue, false)
			publisher.Publish([]byte(""))
		}
	}
}

// Consume sensor name messages from SensorExchange
func (wsc *WebsocketController) listenForSources() {

}

// Consume sensor reading messages from SensorReadingExchange
func (wsc *WebsocketController) listenForReadings() {

}
