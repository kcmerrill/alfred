package alfred

import (
	"os"
	"strings"

	event "github.com/kcmerrill/hook"
)

// NewTask will execute a task
func NewTask(task string, context *Context, loadedTasks map[string]Task) {
	dir, t, tasks := FetchTask(task, context, loadedTasks)

	// register plugins ...
	plugin(t, context, tasks)
	event.Trigger("dir", &dir)
	event.Trigger("t", &t)
	event.Trigger("tasks", &tasks)

	// Skip the task, if we need to skip
	if t.Skip {
		return
	}

	// innocent until proven guilty
	context.Ok = true

	// set our taskname
	_, context.TaskName = TaskParser(task, "alfred:list")

	// interactive mode?
	context.Interactive = t.Interactive

	if !context.hasBeenInited {
		context.hasBeenInited = true
		NewTask(MagicTaskURL(task)+"__init", context, tasks)
	}

	components := []Component{
		Component{"log", log},
		Component{"summary", summary},
		Component{"prompt", prompt},
		Component{"register", register},
		Component{"defaults", defaults},
		Component{"stdin", stdin},
		Component{"config", configC},
		Component{"env", env},
		Component{"check", check},
		Component{"watch", watch},
		Component{"serve", serve},
		Component{"setup", setup},
		Component{"multitask", multitask},
		Component{"tasks", tasksC},
		Component{"for", forC},
		Component{"command", commandC},
		Component{"commands", commands},
		Component{"httptasks", httptasks},
		Component{"result", result},
		Component{"include", includeC},
		Component{"ok", ok},
		Component{"fail", fail},
		Component{"wait", wait},
		Component{"every", every},
	}

	// cycle through our components ...
	event.Trigger("task.started", t, context, tasks)
	for _, component := range components {
		context.Component = component.Name
		event.Trigger("before."+component.Name, context)
		component.F(t, context, tasks)
		event.Trigger("after."+component.Name, context)
		if context.Skip != "" {
			outOK(context.Skip, "skipped", context)
			event.Trigger("task.skipped", context)
			return
		}
	}
	event.Trigger("task.completed", context)
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
	SlackSlashCommands struct {
		token string
		port  string
	} `yaml:"slack.slash.commands"`
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
	Skip        bool
	Interactive bool
	Plugin      map[string]string
	Check       string
	Include     string
}

// Exit determins whether a task should exit or not
func (t *Task) Exit(context *Context, tasks map[string]Task) {
	context.Ok = false
	if t.ExitCode != 0 {
		outFail("["+strings.Join(context.Args, ", ")+"]", "{{ .Text.Failure }}{{ .Text.FailureIcon }} exiting ...", context)
		NewTask("__exit", context, tasks)
		os.Exit(t.ExitCode)
	}

	if t.Skip {
		// skip
		context.Skip = "command"
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
