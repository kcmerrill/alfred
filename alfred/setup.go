package alfred

func setup(task Task, context *Context, tasks map[string]Task) {
	tg := task.ParseTaskGroup(task.Setup)
	execTaskGroup(tg, task, context, tasks)
}
