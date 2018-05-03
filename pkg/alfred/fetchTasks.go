package alfred

import (
	"fmt"
	"os"
	"strings"

	"github.com/kcmerrill/common.go/config"
	"github.com/kcmerrill/common.go/file"
	yaml "gopkg.in/yaml.v2"
)

// FetchTask will fetch the tasks
func FetchTask(task string, context *Context, tasks map[string]Task) (string, Task, map[string]Task) {
	if t, exists := tasks[task]; exists {
		return "./", t, tasks
	}

	if strings.HasPrefix(task, "!") {
		context.TaskName = "exec.command"
		return "./", Task{Summary: "Executing Command", Command: task[1:len(task)], ExitCode: 42}, tasks
	}

	var fetched map[string]Task
	var location string
	var contents []byte

	location, task = TaskParser(task, "alfred:list")

	// hmmm, the task does not exist. Lets try to load whatever possible
	if strings.HasPrefix(location, "http") {
		f, err := file.Get(location)
		if err != nil {
			// cannot use output, no task yet ...
			fmt.Println(translate("{{ .Text.Failure }}{{ .Text.FailureIcon }} "+location+"{{ .Text.Reset }}", emptyContext()))
			os.Exit(42)
		}
		contents = f
	} else {
		// must be local? catalog?
		if strings.HasPrefix(location, "@") {
			location = strings.TrimLeft(location, "@") + "/"
		}

		dir, local, err := config.FindAndCombine(location+"alfred", "yml")
		if err != nil {
			// cannot use output, no task yet ...
			fmt.Println(translate("{{ .Text.Failure }}{{ .Text.FailureIcon }} Missing task file.{{ .Text.Reset }}", emptyContext()))
			os.Exit(42)
		}
		contents = local
		location = dir + string(os.PathSeparator) + location
		os.Chdir(location)
	}

	err := yaml.Unmarshal(contents, &fetched)
	if err != nil {
		// cannot use output, no task yet ...
		outFail("yaml", "Unable to unmarshal "+location, context)
		outFail("yaml", "{{ .Text.Failure }}"+err.Error(), context)
		os.Exit(42)
	}

	context.lock.Lock()
	for fetchedTaskName, fetchedTask := range fetched {
		tasks[fetchedTaskName] = fetchedTask
	}
	context.lock.Unlock()

	if task == "__init" || task == "__exit" {
		return location, Task{Skip: true}, tasks
	}

	if task == "alfred:list" {
		list(context, tasks)
		os.Exit(0)
	}

	if t, exists := tasks[task]; exists {
		return "./", t, tasks
	}

	outFail("invalid task", "{{ .Text.Failure }}'"+task+"'", context)
	os.Exit(42)
	return "./", Task{}, tasks
}
