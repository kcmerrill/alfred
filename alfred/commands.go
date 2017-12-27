package alfred

import (
	"strings"
)

func commands(task Task, context *Context, tasks map[string]Task) {
	if task.Commands == "" || !context.Ok {
		return
	}

	cmds := strings.Split(task.Commands, "\n")
	for _, cmd := range cmds {
		// the task component
		if len(context.Log) != 0 {
			command(cmd, task, context, tasks)
		} else {
			commandBasic(cmd, task, context, tasks)
		}
		if !context.Ok {
			// command failed?
			break
		}
	}
}
