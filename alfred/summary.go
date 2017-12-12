package alfred

import (
	"fmt"

	event "github.com/kcmerrill/hook"
)

func summary(task Task, context *Context, tasks map[string]Task) {
	event.Trigger("speak", "{{  .Text.Success }}{{ .TaskName }}{{ .Text.Reset}}", task, context)
	if task.Summary != "" {
		event.Trigger("speak", "    "+task.Summary, task, context)
	}
	event.Trigger("speak", "    Args: "+fmt.Sprintf("%v", context.Args), task, context)
	event.Trigger("speak", "{{ .Text.Success }}~~~~~~~~~~~~~~~~~~~~{{ .Text.Reset }}\n", task, context)
}

func result(task Task, context *Context, tasks map[string]Task) {
	if context.Ok {
		event.Trigger("speak", "\n{{ .Text.Success }}~~~~~~~~~~~~~~~~~~~~{{ .Text.Reset }}", task, context)
		event.Trigger("speak", "{{  .Text.Success }}{{ .Text.SuccessIcon }} {{ .TaskName }} Ok{{ .Text.Reset}}\n\n", task, context)
		return
	}

	event.Trigger("speak", "\n{{ .Text.Failure }}~~~~~~~~~~~~~~~~~~~~{{ .Text.Reset }}", task, context)
	event.Trigger("speak", "{{  .Text.Failure }}{{ .Text.FailureIcon }} {{ .TaskName }} Failed{{ .Text.Reset}}\n\n", task, context)
}
