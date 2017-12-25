package alfred

import (
	"fmt"
	"os"
)

func log(task Task, context *Context, tasks map[string]Task) {
	if task.Log != "" {
		l := translate(task.Log, context)
		f, err := os.OpenFile(task.Log, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if err == nil {
			context.Log[task.Log] = f
		} else {
			outFail("log", err.Error(), context)
			task.Exit(context, tasks)
		}
	}
}

func logger(text string, context *Context) {
	c := *context
	// strip away all the color
	c.Text = TextConfig{}
	fmt.Println("log", context.Log)
	for _, f := range context.Log {
		f.WriteString(translate(text, context) + "\n")
	}
}
