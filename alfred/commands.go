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
		command(cmd, task, context, tasks)
		if !context.Ok {
			// command failed?
			break
		}
	}
}
