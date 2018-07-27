package alfred

import (
	"fmt"
	"os"

	"github.com/kcmerrill/common.go/config"
)

func include(task Task, context *Context, tasks map[string]Task) {
	if task.Include == "" {
		return
	}

	rp := context.relativePath(task.Include)

	// Does the folder exist? error? boo ...
	if _, err := os.Stat(rp); err != nil {
		context.Ok = false
		return
	}

	_, contents, err := config.FindAndCombine(rp, context.FileName, "yml")
	if err != nil {
		// cannot use output, no task yet ...
		fmt.Println(translate("{{ .Text.Failure }}{{ .Text.FailureIcon }} Missing task file.{{ .Text.Reset }}", emptyContext()))
		os.Exit(42)
	}

	context.rootDir = rp
	tasks = AddTasks(contents, context, tasks)
}
