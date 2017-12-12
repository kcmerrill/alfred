package alfred

func ok(task Task, context *Context, tasks map[string]Task) {
	if !context.Ok {
		return
	}

	execTaskGroup(task.Ok, task, context, tasks)
}
