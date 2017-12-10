package alfred

import (
	"fmt"

	event "github.com/kcmerrill/hook"
)

func summaryHeader(task Task, context *Context) {
	event.Trigger("speak", "{{  .Text.Success }}{{ .TaskName }}{{ .Text.Reset}}", task, context)
	if task.Summary != "" {
		event.Trigger("speak", "    "+task.Summary, task, context)
	}
	event.Trigger("speak", "    Args: "+fmt.Sprintf("%v", context.Args), task, context)
	event.Trigger("speak", "{{ .Text.Success }}~~~~~~~~~~~~~~~~~~~~{{ .Text.Reset }}\n", task, context)
}

func summaryFooter(task Task, context *Context) {
	if context.Ok {
		event.Trigger("speak", "\n{{ .Text.Success }}~~~~~~~~~~~~~~~~~~~~{{ .Text.Reset }}", task, context)
		event.Trigger("speak", "{{  .Text.Success }}{{ .Text.SuccessIcon }} {{ .TaskName }} Ok{{ .Text.Reset}}\n\n", task, context)
		return
	}

	event.Trigger("speak", "\n{{ .Text.Failure }}~~~~~~~~~~~~~~~~~~~~{{ .Text.Reset }}", task, context)
	event.Trigger("speak", "{{  .Text.Failure }}{{ .Text.FailureIcon }} {{ .TaskName }} Failed{{ .Text.Reset}}\n\n", task, context)
}
