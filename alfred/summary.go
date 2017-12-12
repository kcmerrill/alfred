package alfred

import (
	"fmt"

	event "github.com/kcmerrill/hook"
)

func summary(task Task, context *Context, tasks map[string]Task) {
	event.Trigger("output", "{{  .Text.Success }}{{ .TaskName }}{{ .Text.Reset}}", task, context)
	if task.Summary != "" {
		event.Trigger("output", "    "+task.Summary, task, context)
	}
	event.Trigger("output", "    Args: "+fmt.Sprintf("%v", context.Args), task, context)
	event.Trigger("output", "{{ .Text.Success }}~~~~~~~~~~~~~~~~~~~~{{ .Text.Reset }}\n", task, context)
}

func result(task Task, context *Context, tasks map[string]Task) {
	if context.Ok {
		event.Trigger("output", "\n{{ .Text.Success }}~~~~~~~~~~~~~~~~~~~~{{ .Text.Reset }}", task, context)
		event.Trigger("output", "{{  .Text.Success }}{{ .Text.SuccessIcon }} {{ .TaskName }} Ok{{ .Text.Reset}}\n\n", task, context)
		return
	}

	event.Trigger("output", "\n{{ .Text.Failure }}~~~~~~~~~~~~~~~~~~~~{{ .Text.Reset }}", task, context)
	event.Trigger("output", "{{  .Text.Failure }}{{ .Text.FailureIcon }} {{ .TaskName }} Failed{{ .Text.Reset}}\n\n", task, context)
}
