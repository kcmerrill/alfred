package alfred

import (
	"strings"
)

func result(task Task, context *Context, tasks map[string]Task) {

	args := ""
	if len(context.Args) >= 1 {
		args = " (" + strings.Join(context.Args, ", ") + ")"
	}

	if context.Ok {
		output("\n{{ .Text.Success }}~~~~~~~~~~~~~~~~~~~~{{ .Text.Reset }}", task, context)
		output("{{  .Text.Success }}{{ .Text.SuccessIcon }} {{ .TaskName }}"+args+" Ok{{ .Text.Reset}}\n\n", task, context)
		return
	}

	output("\n{{ .Text.Failure }}~~~~~~~~~~~~~~~~~~~~{{ .Text.Reset }}", task, context)
	output("{{  .Text.Failure }}{{ .Text.FailureIcon }} {{ .TaskName }}"+args+" Failed{{ .Text.Reset}}\n\n", task, context)
}
