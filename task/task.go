package task

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/satori/go.uuid"
)

/* Contains our task information
A brief explination:
 - Summary: text showing what the command is about
 - Test: the shell command to run before continuing ..
 - Command: the shell command to run
 - Commands: the shell command to run(except each line is it's own separate command)
 - Usage: explains how to use the command
 - Dir: Change to this directory for this particular command
 - Tasks: a space seperated list of strings that represent tasks to run
 - Every: a string representation of golang duration to run this task for.
 - Wait: a string representation of golang duration that pauses the task before continuing
 - Ok: Similiar to Tasks, except only when the task is succesful
 - Fail: Similiar to Tasks, except only when the task fails
 - AllArgs: Contains all of the arguments passed into alfred(that are not alfred specific)
 - Args: Used for templates to pass in arguments
 - CleanArgs: Only alphanumeric cleaned up arguments
 - UUID: A random UUID string
 - Vars: Used for templates to pass in variables
 - Time: Used primarily for templates inside alfred.yml files
 - Modules: Used for private/public repos and resuseable tasks
 - Defaults: Used to set default values
 - Alias: Space seperated strings determining how else the task can be invoked
 - Private: Some methods should only be called via other tasks and not as a standalone. If so, private accomplishes this
 - Exit: on failure, exit completely with given status code
 - Log: When set, will write to a file(assuming directory structure exists)
 - Retry: When set, will attempt to retry the command X number of times
 - Watch: A regular expression of changed files
 - Setup: Similiar to tasks but these get run _before_ the command/task group gets called
*/
type Task struct {
	Summary   string
	Test      string
	Command   string
	Commands  string
	Usage     string
	Dir       string
	Tasks     string
	Setup     string
	Multitask string
	Every     string
	Wait      string
	Ok        string
	Fail      string
	AllArgs   string
	Args      []string
	CleanArgs []string
	UUID      string
	Vars      map[string]string
	Time      time.Time
	Modules   map[string]string `yaml:",inline"`
	Defaults  []string
	Alias     string
	Private   bool
	Exit      string
	Skip      bool
	Log       string
	Retry     int
	Watch     string
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

/* Return a list of multitasks to run */
func (t *Task) MultiTask() []string {
	return strings.Fields(t.Multitask)
}

/* Return a list of setup tasks to run */
func (t *Task) SetupTasks() []string {
	return strings.Fields(t.Setup)
}

/* Execute a task ... */
func (t *Task) RunCommand(cmd, name string, formatted bool) bool {
	if cmd != "" {
		if t.Log == "" && formatted == false {
			if !t.CommandBasic(cmd) {
				return false
			}
		} else {
			if !t.CommandComplex(cmd, name) {
				return false
			}
		}
		return true
	}
	/* If there was no command to run, then don't fail the task */
	return true
}

/* Execute a complex command(we need to either log, and/or format output ) */
func (t *Task) CommandComplex(cmd, name string) bool {
	if cmd != "" {

		var l *os.File
		var err error

		/* If log is set ... lets use it */
		if t.Log != "" {
			l, err = os.OpenFile(t.Log, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
			if err != nil {
				/* Don't quit ... but don't log either */
				t.Log = ""
			}
			defer l.Close()
		}

		cmd := exec.Command("bash", "-c", cmd)
		cmdReaderStdOut, errStdOut := cmd.StdoutPipe()

		if errStdOut != nil {
			return false
		}

		scannerStdOut := bufio.NewScanner(cmdReaderStdOut)
		go func() {
			for scannerStdOut.Scan() {
				s := fmt.Sprintf("%s\n", scannerStdOut.Text())
				fmt.Printf("%s| %s", name, s)
				if t.Log != "" {
					l.WriteString(s)
				}

			}
		}()

		cmdReaderStdErr, errStdErr := cmd.StderrPipe()

		if errStdErr != nil {
			return false
		}

		scannerStdErr := bufio.NewScanner(cmdReaderStdErr)
		go func() {
			for scannerStdErr.Scan() {
				s := fmt.Sprintf("%s\n", scannerStdErr.Text())
				fmt.Printf("%s:error| %s", name, s)
				if t.Log != "" {
					l.WriteString(s)
				}

			}
		}()

		err = cmd.Start()
		if err != nil {
			return false
		}

		err = cmd.Wait()
		if err != nil {
			return false
		}
	}
	/* If there was no command to run, then don't fail the task */
	return true
}

/* Execute a task ... */
func (t *Task) CommandBasic(cmd string) bool {
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

/* Test
   Currently lets just run a command, BUT, in the future I'd love to see:
   file.exists "filename"
   dir.exists "dirname"
   dir.age 15m
   file.age 15m
   etc etc ..
*/
func (t *Task) TestF(tst string) bool {
	if tst != "" {
		cmd := exec.Command("bash", "-c", tst)
		if cmd.Run() == nil {
			/* Was it succesful? */
			return true
		} else {
			return false
		}
	}
	/* If there was no command to run, then don't fail the test */
	return true
}

/* Evaluate */
func (t *Task) Eval(cmd string) string {
	out, err := exec.Command(cmd).Output()
	if err != nil {
		return cmd
	}
	return string(out)
}

/* Setup a bunch of things, including templates and argument defeaults */
func (t *Task) Prepare(args []string, vars map[string]string) bool {
	t.Args = t.Defaults

	if t.Vars == nil {
		t.Vars = make(map[string]string, 0)
	}

	/* override variable defaults with actual vars */
	for key, value := range vars {
		t.Vars[key] = t.Eval(value)
	}

	/* override defaults with the args */
	for index, value := range args {
		if len(t.Args) > index {
			t.Args[index] = value
		} else {
			t.Args = append(t.Args, value)
		}
	}

	/* Any null values? If so, bail ... */
	for _, value := range t.Args {
		if value == "" {
			return false
		}
	}

	// Cleanup our values
	reg, _ := regexp.Compile("[^A-Za-z0-9]+")
	for _, value := range t.Args {
		t.CleanArgs = append(t.CleanArgs, reg.ReplaceAllString(value, ""))
	}

	// Setup a UUID
	t.UUID = uuid.NewV4().String()

	// Setup time
	t.Time = time.Now()

	/* All of the modules */
	for key, value := range t.Modules {
		if module_ok, module_translated := t.template(value); module_ok {
			t.Modules[key] = module_translated
		} else {
			return false
		}
	}

	/* get to translating */
	if every_ok, every_translated := t.template(t.Every); every_ok {
		t.Every = every_translated
	} else {
		return false
	}

	if allargs_ok, allargs_translated := t.template(strings.Join(args, " ")); allargs_ok {
		t.AllArgs = allargs_translated
	} else {
		return false
	}

	if cmd_ok, cmd_translated := t.template(t.Command); cmd_ok {
		t.Command = cmd_translated
	} else {
		return false
	}

	if cmd_ok, cmd_translated := t.template(t.Commands); cmd_ok {
		t.Commands = cmd_translated
	} else {
		return false
	}

	if dir_ok, dir_translated := t.template(t.Dir); dir_ok {
		t.Dir = dir_translated
	} else {
		return false
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
