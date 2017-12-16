package alfred

func tasksC(task Task, context *Context, tasks map[string]Task) {
	tg := task.ParseTaskGroup(task.Tasks)
	execTaskGroup(tg, task, context, tasks)
}
