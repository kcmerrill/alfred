package alfred

func stdin(task Task, context *Context, tasks map[string]Task) {
	if task.Stdin == "" {
		return
	}

	dir, _ := task.dir(context)
	results := evaluate(translate(task.Stdin, context), dir)
	context.Stdin = results
}
