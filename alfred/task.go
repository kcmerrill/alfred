package alfred

import (
	"os"
	"strings"

	event "github.com/kcmerrill/hook"
)

// NewTask will execute a task
func NewTask(task string, context *Context, loadedTasks map[string]Task) {
	dir, t, tasks := FetchTask(task, context, loadedTasks)

	// Skip the task, if we need to skip
	if t.skip {
		return
	}

	// switch the directory
	os.Chdir(dir)

	// innocent until proven guilty
	context.Ok = true

	// set our taskname
	context.TaskFile, context.TaskName = TaskParser(task, ":default")

	components := []Component{
		Component{"log", log},
		Component{"defaults", defaults},
		Component{"summary", summary},
		Component{"stdin", stdin},
		Component{"config", configC},
		Component{"prompt", prompt},
		Component{"register", register},
		Component{"env", env},
		Component{"serve", serve},
		Component{"setup", setup},
		Component{"multitask", multitask},
		Component{"tasks", tasksC},
		Component{"watch", watch},
		Component{"for", forC},
		Component{"command", commandC},
		Component{"commands", commands},
		Component{"httptasks", httptasks},
		Component{"result", result},
		Component{"ok", ok},
		Component{"fail", fail},
		Component{"wait", wait},
		Component{"every", every},
	}

	// cycle through our components ...
	event.Trigger("task.started", t, context, tasks)
	for _, component := range components {
		context.Component = component.Name
		event.Trigger("before."+component.Name, t, context, tasks)
		component.F(t, context, tasks)
		event.Trigger("after."+component.Name, t, context, tasks)
	}
	event.Trigger("task.completed", t, context, tasks)
}

// Task holds all of our task components
type Task struct {
	Aliases  string
	Summary  string
	Usage    string
	Args     []string
	Setup    string
	Defaults []string
	Dir      string
	For      struct {
		Tasks     string
		MultiTask string
		Args      string
	}
	HTTPTasks struct {
		Port     string
		Password string
	} `yaml:"http.tasks"`
	Config      string
	Log         string
	Every       string
	Command     string
	Retry       int
	Register    map[string]string
	Env         map[string]string
	Commands    string
	Serve       string
	Script      string
	Stdin       string
	Prompt      map[string]string
	Tasks       string
	MultiTask   string
	Ok          string
	Fail        string
	Wait        string
	Watch       string
	Private     bool
	ExitCode    int `yaml:"exit"`
	skip        bool
	Interactive bool
}

// Exit determins whether a task should exit or not
func (t *Task) Exit(context *Context, tasks map[string]Task) {
	context.Ok = false
	if t.ExitCode != 0 {
		outFail("["+strings.Join(context.Args, ", ")+"]", "{{ .Text.Failure }}{{ .Text.FailureIcon }} exiting ...", context)
		NewTask("__exit", context, tasks)
		os.Exit(t.ExitCode)
	}
}

// IsPrivate determines if a task is private
func (t *Task) IsPrivate() bool {
	// I like the idea of not needing to put an astrick next to a task
	// ... Descriptions automagically qualify for "important tasks"
	// No descriptions means it's filler, or private
	// Summaries WITH private: true are private
	if t.Summary == "" || t.Private {
		return true
	}

	return false
}
