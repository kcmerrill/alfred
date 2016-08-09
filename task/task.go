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
	Modules map[string]string `yaml:",inline"`
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

	/* get to translating */
	if every_ok, every_translated := t.template(t.Every); every_ok {
		t.Every = every_translated
	} else {
		return false
	}

	if cmd_ok, cmd_translated := t.template(t.Command); cmd_ok {
		t.Command = cmd_translated
	} else {
		return false
	}

	/* if we made it here, then we are good to go */
	return true
}

func (t *Task) template(translate string) (bool, string) {
	template := template.Must(template.New("translate").Parse(translate))
	b := new(bytes.Buffer)
	err := template.Execute(b, t)
	if err == nil {
		return true, b.String()
	} else {
		return false, translate
	}
	return true, translate
}
