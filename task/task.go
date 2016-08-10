package task

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"
)

/* Contains our task information
A brief explination:
 - Summary: text showing what the command is about
 - Command: the shell command to run
 - Usage: explains how to use the command
 - Dir: Change to this directory for this particular command
 - Tasks: a space seperated list of strings that represent tasks to run
 - Every: a string representation of golang duration to run this task for.
 - Wait: a string representation of golang duration that pauses the task before continuing
 - Ok: Similiar to Tasks, except only when the task is succesful
 - Fail: Similiar to Tasks, except only when the task fails
 - Args: Used for templates to pass in arguments
 - Time: Used primarily for templates inside alfred.yml files
 - Modules: Used for private/public repos and resuseable tasks
 - Defaults: Used to set default values
 - Alias: Space seperated strings determining how else the task can be invoked
 - Private: Some methods should only be called via other tasks and not as a standalone. If so, private accomplishes this
*/
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

/* Is the task private? */
func (t *Task) IsPrivate() bool {
	return t.Private
}

/* Is this task an alias? */
func (t *Task) IsAlias(name string) bool {
	for _, alias := range t.Aliases() {
		if name == alias {
			return true
		}
	}
	return false
}

/* Grab a list of aliases */
func (t *Task) Aliases() []string {
	return strings.Fields(t.Alias)
}

/* Grab a list of failed tasks */
func (t *Task) FailedTasks() []string {
	return strings.Fields(t.Fail)
}

/* Grab a list of tasks to run when succesful */
func (t *Task) OkTasks() []string {
	return strings.Fields(t.Ok)
}

/* Return a list of tasks to run */
func (t *Task) TaskGroup() []string {
	return strings.Fields(t.Tasks)
}

/* Execute a task ... */
func (t *Task) RunCommand(cmd string) bool {
	if cmd != "" {
		cmd := exec.Command("bash", "-c", cmd)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if cmd.Run() == nil {
			/* Was it succesful? */
			return true
		} else {
			return false
		}
	}
	/* If there was no command to run, then don't fail the task */
	return true
}

/* Setup a bunch of things, including templates and argument defeaults */
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

/* Translate a string to a template */
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
