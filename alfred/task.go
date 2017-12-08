package alfred

import (
	"reflect"
	"strings"
)

// Task holds all of our task components
type Task struct {
	Summary     string
	Description string
	Dir         string
	Command     string
	Script      string
}

// NewTask will create a new task
func NewTask(name string, context Context, event *EventStream, tasks Tasks) {
	// component order
	co := []string{
		"summary",
	}
	// does the task exist?
	if task, exists := tasks.Get(name); exists {
		event.Emit(Event{Action: "start.task", Payload: "name"})
		for _, component := range co {
			event.Emit(Event{Action: "start.component", Component: component, Payload: name})
			context = task.Component(component, context, event, tasks)
			event.Emit(Event{Action: "finished.component", Component: component, Payload: name})
		}
		return
	}

	event.Emit(Event{
		Action:  "error.msg",
		Payload: "Invalid task: " + name + "",
	})

	event.Emit(Event{
		Action:  "exit",
		Payload: "1",
	})
}

// Component executes a task component
func (t *Task) Component(component string, args ...interface{}) Context {
	params := make([]reflect.Value, len(args))
	for idx := range args {
		params[idx] = reflect.ValueOf(args[idx])
	}
	if context, err := reflect.ValueOf(t).MethodByName(strings.Title(strings.ToLower(component)) + "Component").Call(params)[0].Interface().(Context); err == true {
		return context
	}
	return Context{}
}
