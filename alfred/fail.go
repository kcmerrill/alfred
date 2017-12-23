package alfred

import "strings"

func fail(task Task, context *Context, tasks map[string]Task) {
	if !context.Ok {
		tgs := task.ParseTaskGroup(task.Fail)

		tgsNames := make([]string, 0)
		for _, tg := range tgs {
			tgsNames = append(tgsNames, tg.Name)
		}

		if len(tgsNames) != 0 {
			outFail("fail.tasks", strings.Join(tgsNames, ", "), context)
			execTaskGroup(tgs, task, context, tasks)
		}
	}
}
