package alfred

import (
	"fmt"
)

func summary(task Task, context *Context, tasks map[string]Task) {
	output("{{  .Text.Success }}{{ .TaskName }}{{ .Text.Reset}}", task, context)
	if task.Summary != "" {
		output("    "+task.Summary, task, context)
	}
	output("    Args: "+fmt.Sprintf("%v", context.Args), task, context)
	output("{{ .Text.Success }}~~~~~~~~~~~~~~~~~~~~{{ .Text.Reset }}\n", task, context)
}

func result(task Task, context *Context, tasks map[string]Task) {
	if context.Ok {
		output("\n{{ .Text.Success }}~~~~~~~~~~~~~~~~~~~~{{ .Text.Reset }}", task, context)
		output("{{  .Text.Success }}{{ .Text.SuccessIcon }} {{ .TaskName }} Ok{{ .Text.Reset}}\n\n", task, context)
		return
	}

	output("\n{{ .Text.Failure }}~~~~~~~~~~~~~~~~~~~~{{ .Text.Reset }}", task, context)
	output("{{  .Text.Failure }}{{ .Text.FailureIcon }} {{ .TaskName }} Failed{{ .Text.Reset}}\n\n", task, context)
}
