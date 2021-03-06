package coordinator

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sophiedebenedetto/power_plant/src/distributed/dto"
	"github.com/sophiedebenedetto/power_plant/src/distributed/messaging"
)

var maxRate = 5 * time.Second

// PersistenceQueue is the queue for messages to persist
var PersistenceQueue = "persistence"

// DatabaseEventRaiser handles DB events
type DatabaseEventRaiser struct {
	EventRaiser EventRaiser
	sources     []string
	publisher   *messaging.Publisher
}

// NewDatabaseEventRaiser returns a new handler
func NewDatabaseEventRaiser(rabbitServer *messaging.Server) *DatabaseEventRaiser {

	dbEventRaiser := &DatabaseEventRaiser{
		EventRaiser: NewEventAggregator(),
		sources:     make([]string, 0),
		publisher:   messaging.NewPublisherWithQueue(rabbitServer, PersistenceQueue, true),
	}
	dbEventRaiser.EventRaiser.AddListener("DataSourceDiscovered", dbEventRaiser.handleDataSourceDiscovered)
	return dbEventRaiser
}

func (dbEventRaiser *DatabaseEventRaiser) handleDataSourceDiscovered(eventName interface{}) {
	logrus.WithField("service", "db").Info("Handling data source discovered in DB event raiser")
	name, _ := eventName.(string)
	name = string(name)
	for _, source := range dbEventRaiser.sources {
		if source == name {
			return
		}
	}
	logrus.WithField("service", "db").Info("Adding message received listener")
	dbEventRaiser.sources = append(dbEventRaiser.sources, name)
	dbEventRaiser.EventRaiser.AddListener("MessageReceived_"+name, dbEventRaiser.handleMessageReceived())
}

func (dbEventRaiser *DatabaseEventRaiser) handleMessageReceived() func(interface{}) {
	prevTime := time.Unix(0, 0)
	return func(event interface{}) {
		logrus.WithField("service", "db").Info("Handling message received in DB event raiser")
		logrus.WithField("service", "db").Info(event)
		ed := event.(EventData)
		if time.Since(prevTime) > maxRate {
			prevTime = time.Now()
			msg := dto.SensorMessage{
				Name:      ed.Name,
				Value:     ed.Value,
				Timestamp: ed.Timestamp,
			}
			dbEventRaiser.publisher.SetUpWriter()
			dbEventRaiser.publisher.WriteMessageToBuffer(msg)
			logrus.WithField("service", "db").Info("PUBLISHING PERSISTENCE MSG...")
			dbEventRaiser.publisher.Publish(dbEventRaiser.publisher.MessageBytes())
		}
	}
}
