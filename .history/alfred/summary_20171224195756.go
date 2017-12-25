package alfred

import (
	"fmt"
	"strings"
	"time"
)

func summary(task Task, context *Context, tasks map[string]Task) {
	fmt.Println("batman", context.TaskName, context.TaskFile)
	if task.Summary != "" {
		outOK("["+strings.Join(context.Args, ", ")+"]", task.Summary, context)
	} else {
		outOK("["+strings.Join(context.Args, ", ")+"]", "started ...", context)
	}
}

func result(task Task, context *Context, tasks map[string]Task) {
	if context.Ok {
		outOK("{{ .Text.SuccessIcon }} ok", "args["+strings.Join(context.Args, ", ")+"] {{ .Text.Reset }}elapsed time {{ .Text.Grey }}'{{ .Text.Success }}"+time.Since(context.Started).Round(time.Second).String()+"{{ .Text.Grey }}'", context)
	} else {
		outFail("{{ .Text.FailureIcon }} failed", "args["+strings.Join(context.Args, ", ")+"] {{ .Text.Reset }}elapsed time {{ .Text.Grey }}'{{ .Text.Success }}"+time.Since(context.Started).Round(time.Second).String()+"{{ .Text.Grey }}'", context)
	}
}
