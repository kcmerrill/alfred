package alfred

func fail(task Task, context *Context, tasks map[string]Task) {
	tg := task.ParseTaskGroup(task.Fail)
	execTaskGroup(tg, task, context, tasks)
}
