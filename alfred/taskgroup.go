package alfred

func taskGroup(group string, task Task, context *Context, tasks map[string]Task) {
	tgs := task.ParseTaskGroup(group)
	for _, tg := range tgs {
		NewTask(tg.Name, InitialContext([]string{}), tasks)
	}
}
