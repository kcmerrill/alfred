package alfred

func fail(task Task, context *Context, tasks map[string]Task) {
	if context.Ok {
		return
	}

	execTaskGroup(task.Fail, task, context, tasks)
}
