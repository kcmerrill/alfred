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
func NewTask(name string, context Context, tasks Tasks) {
	// component order
	co := []string{
		"summary",
	}
	// does the task exist?
	if task, exists := tasks.Get(name); exists {
		for _, component := range co {
			context = task.Component(component, context, tasks)
		}
	}
	// TODO: Task does not exist
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
