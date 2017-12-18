package alfred

import "fmt"

func result(task Task, context *Context, tasks map[string]Task) {
	fmt.Println(context.Ok)
	if context.Ok {
		output("\n{{ .Text.Success }}~~~~~~~~~~~~~~~~~~~~{{ .Text.Reset }}", task, context)
		output("{{  .Text.Success }}{{ .Text.SuccessIcon }} {{ .TaskName }} Ok{{ .Text.Reset}}\n\n", task, context)
		return
	}

	output("\n{{ .Text.Failure }}~~~~~~~~~~~~~~~~~~~~{{ .Text.Reset }}", task, context)
	output("{{  .Text.Failure }}{{ .Text.FailureIcon }} {{ .TaskName }} Failed{{ .Text.Reset}}\n\n", task, context)
}
