package alfred

import (
	"os"
	"strings"
)

// Log will set an external logger
func Log(filename string, context *Context) {
	tasks := make(map[string]Task)
	log(Task{
		Log: filename,
	}, context, tasks)
}
func log(task Task, context *Context, tasks map[string]Task) {
	if task.Log != "" {
		l := translate(strings.TrimSpace(task.Log), context)

		f, err := os.OpenFile(l, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
		if err == nil {
			context.Lock.Lock()
			context.Log[l] = f
			context.Lock.Unlock()
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
	for _, f := range context.Log {
		f.WriteString(translate(text, context))
	}
}
