package alfred

import (
	"strings"
)

func forC(task Task, context *Context, tasks map[string]Task) {
	dir, _ := task.dir(context)
	// alright, lets figure out our new lines
	args := strings.Split(evaluate(translate(task.For.Args, context), dir), "\n")
	tg := make([]TaskGroup, 0)
	// if our tasks list isn't empty, lets loop through it
	if task.For.Tasks != "" {
		// parse our task group(should only be space separated ...)
		tgs := task.ParseTaskGroup(translate(task.For.Tasks, context))
		for _, taskGroup := range tgs {
			// loop through each of our arguments
			for _, arg := range args {
				// append the taskname and the argument
				tg = append(tg, TaskGroup{
					Name: translate(taskGroup.Name, context),
					Args: []string{translate(arg, context)},
				})
			}
		}
		// now, we have all our tasks, lets run them as a task group
		execTaskGroup(tg, task, context, tasks)
	}

	tg = make([]TaskGroup, 0)
	// if our tasks list isn't empty, lets loop through it
	if task.For.MultiTask != "" {
		// parse our task group(should only be space separated ...)
		tgs := task.ParseTaskGroup(translate(task.For.MultiTask, context))
		for _, taskGroup := range tgs {
			// loop through each of our arguments
			for _, arg := range args {
				// append the taskname and the argument
				tg = append(tg, TaskGroup{
					Name: translate(taskGroup.Name, context),
					Args: []string{translate(arg, context)},
				})
			}
		}
		// now, we have all our tasks, lets run them as a task group
		goExecTaskGroup(tg, task, context, tasks)
	}
}
