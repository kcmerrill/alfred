package alfred

import (
	"os"
)

func log(task Task, context *Context, tasks map[string]Task) {
	if task.Log != "" {
		f, err := os.OpenFile(task.Log, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if err == nil {
			context.Log[task.Log] = f
		}
	}
}

func logger(text string, context *Context) {
	c := *context
	c.Text = TextConfig{}
}
