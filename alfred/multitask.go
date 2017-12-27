package alfred

import "strings"

func multitask(task Task, context *Context, tasks map[string]Task) {
	tgs := task.ParseTaskGroup(task.MultiTask)

	tgsNames := make([]string, 0)
	for _, tg := range tgs {
		tgsNames = append(tgsNames, tg.Name)
	}

	if len(tgsNames) != 0 {
		outOK("multitasks", strings.Join(tgsNames, ", "), context)
		goExecTaskGroup(tgs, task, context, tasks)
	}
}
