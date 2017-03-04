package alfred

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
 - Serve: A string, denoting port number to serve a static webserver
*/

// Task contains a task definition and all of it's components
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
	Serve     string
}

// TaskDefinition defines the name and params of a task when originally in string format
type TaskDefinition struct {
	Name   string
	Params []string
}

// IsPrivate returns true/false if the task is private(not executable individually)
func (t *Task) IsPrivate() bool {
	return t.Private
}

// IsAlias returns true/false if the task is an alias of an original task or not
func (t *Task) IsAlias(name string) bool {
	for _, alias := range t.Aliases() {
		if name == alias {
			return true
		}
	}
	return false
}

// Aliases returns an array of strings containing the task's aliases
func (t *Task) Aliases() []string {
	return strings.Fields(t.Alias)
}

//TaskGroup takes in a string(bleh(1234) whatever(bleh, woot)) and returns the values and args
func (t *Task) TaskGroup(tasks string, args []string) []TaskDefinition {
	results := make([]TaskDefinition, 0)
	if tasks == "" {
		// If there is nothing, then there is nothing to report
		return results
	}
	if strings.Index(tasks, "(") == -1 {
		// This means we have a regular space delimited list
		tasks := strings.Split(tasks, " ")
		for _, task := range tasks {
			results = append(results, TaskDefinition{Name: task, Params: args})
		}
	} else {
		// This means we have a group of tasks()
		// Not going to do a ton of error checking.
		// Don't be a sad panda and forget to add () to _EVERY_ task!
		definitions := strings.Split(tasks, ")")
		for _, task := range definitions {
			// Clean up the task in case
			task = strings.TrimSpace(task)
			if task == "" {
				// Empty task? Continue ...
				continue
			}
			// Now, lets separate the task from the params
			if len(strings.Split(task, "(")) == 1 {
				// No args
				results = append(results, TaskDefinition{Name: strings.TrimSpace(strings.Split(task, "(")[0]), Params: args})
			} else {
				taskName := strings.TrimSpace(strings.Split(task, "(")[0])
				p := strings.TrimSpace(strings.Split(task, "(")[1])
				params := strings.Split(p, ",")
				for idx, param := range params {
					// lets clean up our params
					params[idx] = strings.TrimSpace(param)
				}
				results = append(results, TaskDefinition{Name: taskName, Params: params})
			}
		}
	}

	return results
}

// RunCommand runs a command, also determining if it needs to be formated(multitasks for example)
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

// CommandComplex takes in a command and will write it to a file, or special formatting(multitask for example)
func (t *Task) CommandComplex(cmd, name string) bool {
	if cmd != "" {

		var l *os.File
		var err error

		// If log is set ... lets use it
		if t.Log != "" {
			l, err = os.OpenFile(t.Log, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
			if err != nil {
				// Don't quit ... but don't log either
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
	// If there was no command to run, then don't fail the task
	return true
}

// CommandBasic runs a basic command, no frills
func (t *Task) CommandBasic(cmd string) bool {
	if cmd != "" {
		cmd := exec.Command("bash", "-c", cmd)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if cmd.Run() == nil {
			// Was it succesful?
			return true
		}
		return false
	}
	// If there was no command to run, then don't fail the task
	return true
}

// TestF runs a command with _NO_ output!
func (t *Task) TestF(tst string) bool {
	if tst != "" {
		cmd := exec.Command("bash", "-c", tst)
		if cmd.Run() == nil {
			/* Was it succesful? */
			return true
		}
		return false
	}
	/* If there was no command to run, then don't fail the test */
	return true
}

// Eval runs a string to see if it's a command or not(depending on it's exit code)
func (t *Task) Eval(cmd string) string {
	out, err := exec.Command(cmd).Output()
	if err != nil {
		return cmd
	}
	return string(out)
}

// Prepare will setup a bunch of things, including templates and argument defeaults
func (t *Task) Prepare(args []string, vars map[string]string) bool {
	t.Args = t.Defaults

	if t.Vars == nil {
		t.Vars = make(map[string]string, 0)
	}

	// override variable defaults with actual vars
	for key, value := range vars {
		t.Vars[key] = t.Eval(value)
	}

	// override defaults with the args
	for index, value := range args {
		if len(t.Args) > index {
			t.Args[index] = value
		} else {
			t.Args = append(t.Args, value)
		}
	}

	// Any null values? If so, bail ...
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

	// All of the modules
	for key, value := range t.Modules {
		if moduleOk, moduleTranslated := t.template(value); moduleOk {
			t.Modules[key] = moduleTranslated
		} else {
			return false
		}
	}

	// get to translating
	if everyOk, everyTranslated := t.template(t.Every); everyOk {
		t.Every = everyTranslated
	} else {
		return false
	}

	if allargsOk, allargsTranslated := t.template(strings.Join(args, " ")); allargsOk {
		t.AllArgs = allargsTranslated
	} else {
		return false
	}

	if cmdOk, cmdTranslated := t.template(t.Command); cmdOk {
		t.Command = cmdTranslated
	} else {
		return false
	}

	if cmdOk, cmdTranslated := t.template(t.Commands); cmdOk {
		t.Commands = cmdTranslated
	} else {
		return false
	}

	if dirOk, dirTranslated := t.template(t.Dir); dirOk {
		t.Dir = dirTranslated
	} else {
		return false
	}

	if tstOk, tstTranslated := t.template(t.Test); tstOk {
		t.Test = tstTranslated
	} else {
		return false
	}

	// if we made it here, then we are good to go
	return true
}

// template a helper function to translate a string to a template
func (t *Task) template(translate string) (bool, string) {
	template := template.Must(template.New("translate").Parse(translate))
	b := new(bytes.Buffer)
	err := template.Execute(b, t)
	if err == nil {
		return true, b.String()
	}
	return false, translate
}
