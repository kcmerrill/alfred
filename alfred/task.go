package alfred

import (
	"bytes"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	event "github.com/kcmerrill/hook"
)

// NewTask will execute a task
func NewTask(task string, context *Context, tasks map[string]Task) {
	t, exists := tasks[task]

	if !exists {
		// TODO, Lookup ... then exit
		event.Trigger("speak", "shit is broke", Task{}, context)
		return
	}

	c := context
	c.TaskName = task

	event.Trigger("task.group", t.Setup, t, c, tasks)
	event.Trigger("task.summary.header", t, c)
	event.Trigger("task.command", t, c)
	event.Trigger("task.serve", t, c)
	event.Trigger("task.summary.footer", t, c)
	if c.Ok {
		event.Trigger("task.group", t.Ok, t, c, tasks)
	} else {
		event.Trigger("task.group", t.Fail, t, c, tasks)
	}
	event.Trigger("task.wait", t, c)
}

// Task holds all of our task components
type Task struct {
	Aliases     string
	Summary     string
	Description string
	Args        []string
	Setup       string
	Dir         string
	Command     string
	Serve       string
	Script      string
	Ok          string
	Fail        string
	Private     bool
	Wait        string
	ExitCode    int
}

// TaskGroup contains a task name and it's arguments
type TaskGroup struct {
	Name string
	Args []string
}

// ParseTaskGroup takes in a string, and parses it into a TaskGroup
func (t *Task) ParseTaskGroup(group string) []TaskGroup {
	tg := make([]TaskGroup, 0)
	group = strings.TrimSpace(group)

	if group == "" {
		return tg
	}

	if strings.Index(group, "\n") == -1 {
		// This means we have a regular space delimited list
		tasks := strings.Split(group, " ")
		for _, task := range tasks {
			tg = append(tg, TaskGroup{Name: task, Args: []string{}})
		}
	} else {
		// mix and match here
		tasks := strings.Split(group, "\n")
		for _, task := range tasks {
			re := regexp.MustCompile(`(.*?)\((.*?)\)`)
			results := re.FindStringSubmatch(task)
			if len(results) == 0 {
				tg = append(tg, TaskGroup{Name: strings.TrimSpace(task), Args: []string{}})
			} else {
				args := strings.Split(results[2], ",")
				for idx, a := range args {
					// trim the extra whitespace
					args[idx] = strings.TrimSpace(a)
				}
				tg = append(tg, TaskGroup{Name: strings.TrimSpace(results[1]), Args: args})
			}
		}
	}

	return tg
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
		event.Trigger("speak", "{{ .Text.Failure }} Bad Template: "+err.Error()+"{{ .Text.Reset }}", t, context)
		return ""
	}
	return b.String()
}
