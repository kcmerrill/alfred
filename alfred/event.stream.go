package alfred

import (
	"fmt"
)

// EventStream contains a stream of incoming events
type EventStream struct {
	Stream chan Event
}

// NewEventStream returns a functioning event stream
func NewEventStream(es *EventStream) *EventStream {
	es.Stream = make(chan Event)
	go es.Collect()
	return es
}

// Emit adds an event to the stream
func (e *EventStream) Emit(event Event) bool {
	e.Stream <- event
	return true
}

// Collect will take in events for processing
func (e *EventStream) Collect() {
	for {
		event := <-e.Stream
		fmt.Println(event)
	}
}
