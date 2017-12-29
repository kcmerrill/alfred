package alfred

func stdin(task Task, context *Context, tasks map[string]Task) {
	if task.Stdin == "" {
		return
	}

	dir, _ := task.dir(context)
	results, err := execute(translate(task.Stdin, context), dir)
	if err == false {
		context.Stdin = results
	} else {
		task.Exit(context, tasks)
	}
}
