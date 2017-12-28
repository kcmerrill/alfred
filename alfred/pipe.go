package alfred

func pipe(task Task, context *Context, tasks map[string]Task) {
	if task.Pipe == "" {
		return
	}

	dir, _ := task.dir(context)
	results, err := execute(translate(task.Pipe, context), dir)
	if err == false {
		context.Stdin = results
	} else {
		task.Exit(context, tasks)
	}
}
