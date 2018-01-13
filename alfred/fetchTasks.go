package alfred

import (
	"fmt"
	"os"

	"github.com/kcmerrill/common.go/config"
	"github.com/kcmerrill/common.go/file"
	yaml "gopkg.in/yaml.v2"
)

// FetchTask will fetch the tasks
func FetchTask(task string, context *Context, tasks map[string]Task) (string, Task, map[string]Task) {
	if t, exists := tasks[task]; exists {
		return "./", t, tasks
	}

	var fetched map[string]Task
	var location string
	var contents []byte

	location, task = TaskParser(task, ":list")

	// hmmm, the task does not exist. Lets try to load whatever possible
	if location != ":local" {
		f, err := file.Get(location)
		if err != nil {
			// cannot use output, no task yet ...
			fmt.Println(translate("{{ .Text.Failure }}{{ .Text.FailureIcon }} "+location+"{{ .Text.Reset }}", emptyContext()))
			os.Exit(42)
		}
		contents = f
	} else {
		// must be local
		dir, local, err := config.FindAndCombine("alfred", "yml")
		if err != nil {
			// cannot use output, no task yet ...
			fmt.Println(translate("{{ .Text.Failure }}{{ .Text.FailureIcon }} Missing task file.{{ .Text.Reset }}", emptyContext()))
			os.Exit(42)
		}
		os.Chdir(dir)
		contents = local
	}

	err := yaml.Unmarshal(contents, &fetched)
	if err != nil {
		// cannot use output, no task yet ...
		outFail("yaml", "Unable to unmarshal "+location, context)
		outFail("yaml", "{{ .Text.Failure }}"+err.Error(), context)
		os.Exit(42)
	}

	for fetchedTaskName, fetchedTask := range fetched {
		tasks[fetchedTaskName] = fetchedTask
	}

	if t, exists := tasks[task]; exists {
		return "", t, tasks
	}

	if task == "__init" {
		return "__init", Task{skip: true}, tasks
	}

	if task == "__exit" {
		return "__exit", Task{skip: true}, tasks
	}

	if task == ":list" {
		list(context, tasks)
		os.Exit(0)
	}

	outFail("invalid task", "{{ .Text.Failure }}'"+task+"'", context)
	os.Exit(42)
	return "", Task{}, tasks
}
