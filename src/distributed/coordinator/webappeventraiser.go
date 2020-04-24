package coordinator

import (
	"github.com/sirupsen/logrus"
	"github.com/sophiedebenedetto/power_plant/src/distributed/dto"
	"github.com/sophiedebenedetto/power_plant/src/distributed/messaging"
)

// SensorExchange is the exchange for sensors
var SensorExchange = "WebAppSensors"

// SensorReadingExchange is the exchange for sensor readings
var SensorReadingExchange = "WebAppSensorReadings"

// WebAppEventRaiser handles DB events
type WebAppEventRaiser struct {
	EventRaiser            EventRaiser
	sources                []string
	sensorPublisher        *messaging.Publisher
	sensorReadingPublisher *messaging.Publisher
}

// NewWebAppEventRaiser returns a new handler
func NewWebAppEventRaiser(rabbitServer *messaging.Server) *WebAppEventRaiser {

	webAppEventRaiser := &WebAppEventRaiser{
		EventRaiser:            NewEventAggregator(),
		sources:                make([]string, 0),
		sensorPublisher:        messaging.NewPublisherWithExchange(rabbitServer, SensorExchange),
		sensorReadingPublisher: messaging.NewPublisherWithExchange(rabbitServer, SensorReadingExchange),
	}
	webAppEventRaiser.EventRaiser.AddListener("DataSourceDiscovered", webAppEventRaiser.handleDataSourceDiscovered)
	webAppEventRaiser.EventRaiser.AddListener("WebAppDiscoveryRequest", webAppEventRaiser.handleWebAppDiscovery)
	return webAppEventRaiser
}

func (webAppEventRaiser *WebAppEventRaiser) handleDataSourceDiscovered(eventName interface{}) {
	logrus.WithField("service", "web").Info("Handling data source discovered in web app event raiser")
	name, _ := eventName.(string)
	name = string(name)
	for _, source := range webAppEventRaiser.sources {
		if source == name {
			return
		}
	}
	logrus.WithField("service", "web").Info("Adding message received listener")
	webAppEventRaiser.sources = append(webAppEventRaiser.sources, name)
	webAppEventRaiser.sensorPublisher.Publish([]byte(name)) // publish new sensor name to WebAppSensors
	webAppEventRaiser.EventRaiser.AddListener("MessageReceived_"+name, webAppEventRaiser.handleMessageReceived)
}

func (webAppEventRaiser *WebAppEventRaiser) handleMessageReceived(event interface{}) {
	logrus.WithField("service", "web").Info("Handling message received in web app event raiser")
	ed := event.(EventData)
	msg := dto.SensorMessage{
		Name:      ed.Name,
		Value:     ed.Value,
		Timestamp: ed.Timestamp,
	}
	webAppEventRaiser.sensorReadingPublisher.SetUpWriter()
	webAppEventRaiser.sensorReadingPublisher.WriteMessageToBuffer(msg)
	logrus.WithField("service", "web").Info("Publishing web app message...")
	// publish sensor reading to WebAppSensorReadings
	webAppEventRaiser.sensorReadingPublisher.Publish(webAppEventRaiser.sensorReadingPublisher.MessageBytes())
}

func (webAppEventRaiser *WebAppEventRaiser) handleWebAppDiscovery(event interface{}) {
	logrus.WithField("service", "web").Info("Handling web app discovery message received in web app event raiser")
	for _, source := range webAppEventRaiser.sources {
		webAppEventRaiser.sensorPublisher.Publish([]byte(source))
	}
}
