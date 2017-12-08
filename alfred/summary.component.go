package alfred

// SummaryComponent prints out the task summary
func (t Task) SummaryComponent(context Context, event *EventStream, tasks Tasks) Context {
	event.Emit(Event{
		Action:    "print",
		Component: "summary",
		Payload:   "{{ .successColor }}{{ .okIcon }} " + t.Summary + " {{ .resetColor }}",
	})
	return context
}
