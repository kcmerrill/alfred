package alfred

import (
	"bytes"
	"os"
	"text/template"

	"github.com/Masterminds/sprig"
	event "github.com/kcmerrill/hook"
)

// NewTask will execute a task
func NewTask(task string, context *Context, tasks map[string]Task) {
	t, exists := tasks[task]

	if !exists {
		// TODO, Lookup ... then exit
		output("shit is broke", Task{}, context)
		return
	}

	// copy our context
	c := context

	// set our taskname
	c.TaskName = task

	// lets setup our task groups
	t.Setup = t.ParseTaskGroup(t.setup)
	t.Tasks = t.ParseTaskGroup(t.tasks)
	t.Ok = t.ParseTaskGroup(t.ok)
	t.Fail = t.ParseTaskGroup(t.fail)

	components := []string{
		"setup",
		"summary",
		"watch",
		"command",
		"serve",
		"result",
		"ok",
		"fail",
		"wait",
		"every",
	}

	event.Trigger("task.started", task)
	// cycle through our components ...
	for _, component := range components {
		event.Trigger("before."+component, t, context, tasks)
		// todo: turn into invoke(s)
		event.Trigger(component, t, context, tasks)
		event.Trigger("after."+component, t, context, tasks)
	}
	event.Trigger("task.completed", task)
}

// Task holds all of our task components
type Task struct {
	Aliases     string
	Summary     string
	Description string
	Usage       string
	Args        []string
	setup       string
	Setup       []TaskGroup
	Dir         string
	Every       string
	Command     string
	Serve       string
	Script      string
	tasks       string
	Tasks       []TaskGroup
	ok          string
	Ok          []TaskGroup
	fail        string
	Fail        []TaskGroup
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

// Template is a helper function to translate a string to a template
func (t *Task) Template(translate string, context *Context) string {
	if translate == "" {
		// Nothing to translate, move along
		return translate
	}
	fmap := sprig.TxtFuncMap()
	te := template.Must(template.New("template").Funcs(fmap).Parse(translate))
	var b bytes.Buffer
	err := te.Execute(&b, context)
	if err != nil {
		output("{{ .Text.Failure }}{{ .Text.FailureIcon }}Bad Template: "+err.Error()+"{{ .Text.Reset }}", *t, context)
		return ""
	}
	return b.String()
}

// IsPrivate determines if a task is private
func (t *Task) IsPrivate() bool {
	// I like the idea of not needing to put an astrick next to a task
	// ... Descriptions and usage automagically qualify for "important tasks"
	// No descriptions, or usage information means it's filler, or private
	if t.Description != "" || t.Usage != "" {
		return false
	}

	return true
}
