package alfred

import (
	"github.com/kcmerrill/hook"
)

func plugin(task Task, context *Context, tasks map[string]Task) {
	// register our plugins
	if len(task.Plugin) == 0 {
		return
	}

	// sweet ... we have plugins.
	for key, value := range task.Plugin {
		hook.Register(key, value)
		outOK("plugin {{ .Text.Args }}"+key+"{{ .Text.Reset }}", value, context)
	}
}
