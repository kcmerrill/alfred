package alfred

import (
	"fmt"
	"os"

	"github.com/kcmerrill/common.go/config"
)

func includeC(task Task, context *Context, tasks map[string]Task) {
	if task.Include == "" {
		return
	}

	// Does the folder exist? error? boo ...
	if _, err := os.Stat(task.Include); err != nil {
		context.Ok = false
		return
	}

	_, contents, err := config.FindAndCombine(task.Include, "alfred", "yml")
	if err != nil {
		// cannot use output, no task yet ...
		fmt.Println(translate("{{ .Text.Failure }}{{ .Text.FailureIcon }} Missing task file.{{ .Text.Reset }}", emptyContext()))
		os.Exit(42)
	}

	context.rootDir = task.Include
	tasks = AddTasks(contents, context, tasks)
}
