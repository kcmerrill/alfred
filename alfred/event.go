package alfred

// Event is a payload to be displayed or acted upon
type Event struct {
	Component string
	Action    string
	Payload   string
}

// NewEvent returns a preconfigured event
func NewEvent(action, component, payload string) Event {
	return Event{
		Action:    action,
		Component: component,
		Payload:   payload,
	}
}
