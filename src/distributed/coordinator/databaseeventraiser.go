package coordinator

import "fmt"

// DatabaseEventRaiser handles DB events
type DatabaseEventRaiser struct {
	EventRaiser EventRaiser
	sources     []string
}

// NewDatabaseEventRaiser returns a new handler
func NewDatabaseEventRaiser() *DatabaseEventRaiser {
	dbEventRaiser := &DatabaseEventRaiser{
		EventRaiser: NewEventAggregator(),
		sources:     make([]string, 0),
	}
	dbEventRaiser.EventRaiser.AddListener("DataSourceDiscovered", dbEventRaiser.handleDataSourceDiscovered)
	return dbEventRaiser
}

func (dbEventRaiser *DatabaseEventRaiser) handleDataSourceDiscovered(eventName interface{}) {
	fmt.Println("Handling data source discovered in DB event raiser")
	name, _ := eventName.(string)
	name = string(name)
	for _, source := range dbEventRaiser.sources {
		if source == name {
			return
		}
	}
	fmt.Println("Adding message received listener")
	dbEventRaiser.sources = append(dbEventRaiser.sources, name)
	dbEventRaiser.EventRaiser.AddListener("MessageReceived_"+name, dbEventRaiser.handleMessageReceived)
}

func (dbEventRaiser *DatabaseEventRaiser) handleMessageReceived(event interface{}) {
	// every 5 seconds, publish to persistence queue
	fmt.Println("Handling message received in DB event raiser")
	fmt.Println(event)
}
