package alfred

import (
	"fmt"
	"strings"
)

func commands(task Task, context *Context, tasks map[string]Task) {
	if task.Commands == "" {
		return
	}

	cmds := strings.Split(task.Commands, "\n")
	for _, cmd := range cmds {
		// the task component
		fmt.Println("command to run:", cmd)
		command(cmd, task, context, tasks)
	}
}
