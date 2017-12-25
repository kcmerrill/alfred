package alfred

import (
	"fmt"
	"os"
)

func log(task Task, context *Context, tasks map[string]Task) {
	if task.Log != "" {
		f, err := os.OpenFile(task.Log, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if err == nil {
			context.Log[task.Log] = f
			fmt.Println("kk, so we are here.")
		} else {
			outFail("log", err.Error(), context)
		}
	}
}

func logger(text string, context *Context) {
	c := *context
	// strip away all the color
	c.Text = TextConfig{}
	for _, f := range context.Log {
		f.WriteString(translate(text, context) + "\n")
	}
}
