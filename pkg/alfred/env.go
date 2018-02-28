package alfred

import (
	"os"
)

func env(task Task, context *Context, tasks map[string]Task) {
	if len(task.Env) == 0 {
		return
	}

	dir, _ := task.dir(context)

	for key, value := range task.Env {
		keyT := translate(key, context)
		valueT := evaluate(translate(value, context), dir)
		os.Setenv(keyT, valueT)
		outOK("set {{ .Text.Args }}$"+keyT+"{{ .Text.Reset }}", valueT, context)
	}
}
