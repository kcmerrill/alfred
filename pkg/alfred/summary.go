package alfred

import (
	"strings"
	"time"
)

func summary(task Task, context *Context, tasks map[string]Task) {
	argsStr := ""
	if len(context.Args) > 0 {
		argsStr = " [" + strings.Join(context.Args, ", ") + "]"
	}
	if task.Summary != "" {
		outOK("started"+argsStr, task.Summary, context)
	} else {
		outOK("started"+argsStr, "", context)
	}
}

func result(task Task, context *Context, tasks map[string]Task) {
	argsStr := ""
	if len(context.Args) > 0 {
		argsStr = " [" + strings.Join(context.Args, ", ") + "]"
	}
	if context.Ok {
		outOK("{{ .Text.SuccessIcon }} ok"+argsStr, "in {{ .Text.Grey }}{{ .Text.Success }}"+time.Since(context.TaskStarted).Round(time.Second).String()+"{{ .Text.Grey }}", context)
	} else {
		outFail("{{ .Text.FailureIcon }} failed"+argsStr, "elapsed time {{ .Text.Grey }}'{{ .Text.Success }}"+time.Since(context.TaskStarted).Round(time.Second).String()+"{{ .Text.Grey }}'", context)
	}
}
