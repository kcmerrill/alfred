package alfred

import "strings"

func setup(task Task, context *Context, tasks map[string]Task) {
	tgs := task.ParseTaskGroup(task.Setup)

	tgsNames := make([]string, 0)
	for _, tg := range tgs {
		tgsNames = append(tgsNames, tg.Name)
	}

	if len(tgsNames) != 0 {
		outOK("setup", strings.Join(tgsNames, ", "), context)
		execTaskGroup(tgs, task, context, tasks)
	}
}
