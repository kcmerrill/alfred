package alfred

import (
	"os"

	event "github.com/kcmerrill/hook"
)

// NewTask will execute a task
func NewTask(task string, context *Context, loadedTasks map[string]Task) {
	t, tasks := FetchTask(task, loadedTasks)

	// copy our context
	c := context

	// innocent until proven guilty
	c.Ok = true

	// set our taskname
	c.TaskFile, c.TaskName = TaskParser(task, ":default")

	// cycle through our components
	components := []Component{
		Component{"setup", setup},
		Component{"tasks", tasksC},
		Component{"watch", watch},
		Component{"command", command},
		Component{"serve", serve},
		Component{"result", result},
		Component{"ok", ok},
		Component{"fail", fail},
		Component{"wait", wait},
		Component{"every", every},
	}

	// cycle through our components ...
	event.Trigger("task.started", t, context, tasks)
	for _, component := range components {
		event.Trigger("before."+component.Name, t, context, tasks)
		component.F(t, context, tasks)
		event.Trigger("after."+component.Name, t, context, tasks)
	}
	event.Trigger("task.completed", t, context, tasks)
}

// Task holds all of our task components
type Task struct {
	Aliases     string
	Summary     string
	Description string
	Args        []string
	Setup       string
	Dir         string
	Every       string
	Command     string
	Serve       string
	Script      string
	Tasks       string
	Ok          string
	Fail        string
	Wait        string
	Watch       string
	ExitCode    int
}

// Exit determins whether a task should exit or not
func (t *Task) Exit() {
	if t.ExitCode != 0 {
		os.Exit(t.ExitCode)
	}
}

// IsPrivate determines if a task is private
func (t *Task) IsPrivate() bool {
	// I like the idea of not needing to put an astrick next to a task
	// ... Descriptions automagically qualify for "important tasks"
	// No descriptions means it's filler, or private
	if t.Description != "" {
		return false
	}

	return true
}
