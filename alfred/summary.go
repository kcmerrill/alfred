package alfred

import "strings"

func summary(task Task, context *Context, tasks map[string]Task) {
	if task.Summary != "" {
		outOK("summary", task.Summary, context)
		outOK("args", "["+strings.Join(context.Args, ", ")+"]", context)
	}
}
