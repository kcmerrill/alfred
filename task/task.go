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
	Summary  string
	Command  string
	Usage    string
	Dir      string
	Tasks    string
	Every    string
	Wait     string
	Ok       string
	Fail     string
	Args     []string
	Time     *time.Time
	Modules  map[string]string `yaml:",inline"`
	Defaults []string
	Alias    string
	Private  bool
}

func (t *Task) IsPrivate() bool {
	return t.Private
}

func (t *Task) IsAlias(name string) bool {
	for _, alias := range t.Aliases() {
		if name == alias {
			return true
		}
	}
	return false
}

func (t *Task) Aliases() []string {
	return strings.Fields(t.Alias)
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

func (t *Task) RunCommand(cmd string, args []string) bool {
	if cmd != "" {
		cmd := exec.Command("bash", "-c", cmd)
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
	t.Args = t.Defaults

	/* override defaults with the args */
	for index, value := range args {
		if len(t.Args) > index {
			t.Args[index] = value
		} else {
			t.Args = append(t.Args, value)
		}
	}

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

	if dir_ok, dir_translated := t.template(t.Dir); dir_ok {
		t.Dir = dir_translated
	} else {
		return false
	}

	/* All of the modules */
	for key, value := range t.Modules {
		if module_ok, module_translated := t.template(value); module_ok {
			t.Modules[key] = module_translated
		} else {
			return false
		}
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
