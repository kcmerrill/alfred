package alfred

import (
	"bufio"
	"os"
	"strings"
)

func prompt(task Task, context *Context, tasks map[string]Task) {
	if len(task.Prompt) == 0 {
		return
	}

	context.Interactive = true

	for v, phrase := range task.Prompt {
		reader := bufio.NewReader(os.Stdin)
		outputPrompt("Success", "prompt", phrase+" ", context)
		p, _ := reader.ReadString('\n')
		context.SetVar(v, strings.TrimSpace(p))
	}

	for v := range task.Prompt {
		outOK("registered {{ .Text.Args }}"+v+"{{ .Text.Reset }}", context.GetVar(v, ""), context)
	}
}
