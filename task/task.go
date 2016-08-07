package task

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"
)

type Task struct {
	Summary string
	Command string
	Usage   string
	Tasks   string
	Every   string
	Wait    string
	Ok      string
	Fail    string
	Args    []string
	Time    *time.Time
}

func (t *Task) FailedTasks() []string {
	return strings.Fields(t.Fail)
}

func (t *Task) OkTasks() []string {
	return strings.Fields(t.Ok)
}

func (t *Task) TaskGroup() []string {
	return strings.Fields(t.Tasks)
}

func (t *Task) RunCommand(args []string) bool {
	if t.Command != "" {
		cmd := exec.Command("bash", "-c", t.Command)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if cmd.Run() == nil {
			return true
		} else {
			return false
		}
	}
	return true
}

func (t *Task) Prepare(args []string) bool {
	t.Args = args
	t.Time = new(time.Time)

	/* Make sure the every isn't empty ... */
	if t.Every != "" {
		template := template.Must(template.New("every").Parse(t.Every))
		b := new(bytes.Buffer)
		err := template.Execute(b, t)
		if err == nil {
			t.Every = b.String()
		} else {
			return false
		}
	}

	/* Make sure the command isn't empty ... */
	if t.Command != "" {
		template := template.Must(template.New("command").Parse(t.Command))
		b := new(bytes.Buffer)
		err := template.Execute(b, t)
		if err == nil {
			t.Command = b.String()
		} else {
			return false
		}
	}
	/* If there isn't a command, then no failures can be there if we try to prep it */
	return true
}
