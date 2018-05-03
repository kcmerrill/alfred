package alfred

import (
	"strings"
)

func tasksC(task Task, context *Context, tasks map[string]Task) {
	tgs := task.ParseTaskGroup(task.Tasks)

	tgsNames := make([]string, 0)
	for _, tg := range tgs {
		tgsNames = append(tgsNames, tg.Name)
	}

	if len(tgsNames) != 0 {
		outOK("tasks", strings.Join(tgsNames, ", "), context)
		execTaskGroup(tgs, task, context, tasks)
	}
}
