package alfred

func check(task Task, context *Context, tasks map[string]Task) {
	if task.Check == "" {
		return
	}

	dir, _ := task.dir(context)
	if testCommand(translate(task.Check, context), dir) {
		context.Skip = "check"
	}
}
