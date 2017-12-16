package alfred

func ok(task Task, context *Context, tasks map[string]Task) {
	tg := task.ParseTaskGroup(task.Ok)
	execTaskGroup(tg, task, context, tasks)
}
