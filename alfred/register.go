package alfred

func register(task Task, context *Context, tasks map[string]Task) {
	if len(task.Register) == 0 {
		return
	}

	dir, _ := task.dir(context)

	for key, value := range task.Register {
		keyT := translate(key, context)
		valueT := evaluate(translate(value, context), dir)
		context.SetVar(keyT, valueT)
		outOK("registered {{ .Text.Args }}"+keyT+"{{ .Text.Reset }}", valueT, context)
	}
}
