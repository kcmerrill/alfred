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
		event.Trigger("output", "shit is broke", Task{}, context)
		return
	}

	// copy our context
	c := context

	// set our taskname
	c.TaskName = task

	// lets setup our task groups
	t.Setup = t.ParseTaskGroup(t.SetupStr)
	t.Tasks = t.ParseTaskGroup(t.TasksStr)
	t.Ok = t.ParseTaskGroup(t.OkStr)
	t.Fail = t.ParseTaskGroup(t.FailStr)

	components := []string{
		"setup",
		"summary",
		"command",
		"serve",
		"result",
		"ok",
		"fail",
		"wait",
	}

	event.Trigger("start.task", task)
	// cycle through our components ...
	for _, component := range components {
		event.Trigger("before."+component, t, context, tasks)
		event.Trigger(component, t, context, tasks)
		event.Trigger("after."+component, t, context, tasks)
	}
	event.Trigger("completed.task", task)
}

// Task holds all of our task components
type Task struct {
	Aliases     string
	Summary     string
	Description string
	Args        []string
	SetupStr    string
	Setup       []TaskGroup
	Dir         string
	Command     string
	Serve       string
	Script      string
	TasksStr    string
	Tasks       []TaskGroup
	OkStr       string
	Ok          []TaskGroup
	FailStr     string
	Fail        []TaskGroup
	Private     bool
	Wait        string
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
		event.Trigger("output", "{{ .Text.Failure }} Bad Template: "+err.Error()+"{{ .Text.Reset }}", t, context)
		return ""
	}
	return b.String()
}
