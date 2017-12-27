package alfred

import "strings"

func ok(task Task, context *Context, tasks map[string]Task) {
	tgs := task.ParseTaskGroup(task.Ok)

	tgsNames := make([]string, 0)
	for _, tg := range tgs {
		tgsNames = append(tgsNames, tg.Name)
	}

	if len(tgsNames) != 0 {
		outOK("ok.tasks", strings.Join(tgsNames, ", "), context)
		execTaskGroup(tgs, task, context, tasks)
	}
}
