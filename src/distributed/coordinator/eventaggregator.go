package coordinator

import "time"

// EventRaiser raises events
type EventRaiser interface {
	AddListener(name string, f func(interface{}))
	PublishEvent(name string, event interface{})
}

// EventAggregator knows to to respond to events
type EventAggregator struct {
	listeners map[string][]func(interface{})
}

// EventData describes an event
type EventData struct {
	Name      string
	Value     float64
	Timestamp time.Time
}

// NewEventAggregator returns a new NewEventAggregator
func NewEventAggregator() *EventAggregator {
	return &EventAggregator{
		listeners: make(map[string][]func(interface{})),
	}
}

// AddListener adds a listener
func (ea *EventAggregator) AddListener(name string, f func(interface{})) {
	ea.listeners[name] = append(ea.listeners[name], f)
}

// PublishEvent publishes an event
func (ea *EventAggregator) PublishEvent(name string, event interface{}) {
	if ea.listeners[name] != nil {
		for _, f := range ea.listeners[name] {
			f(event)
		}
	}
}
