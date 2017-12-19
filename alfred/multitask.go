package alfred

func multitask(task Task, context *Context, tasks map[string]Task) {
	tg := task.ParseTaskGroup(task.MultiTask)
	goExecTaskGroup(tg, task, context, tasks)
}
