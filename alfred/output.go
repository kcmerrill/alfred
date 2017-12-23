package alfred

import (
	"fmt"
	"time"
)

func outOK(component, text string, context *Context) {
	output("Success", component, text+"\n", context)
}

func outFail(component, text string, context *Context) {
	output("Failure", component, text+"\n", context)
}

func outWarn(component, text string, context *Context) {
	output("Warning", component, text+"\n", context)
}

func output(color, component, text string, context *Context) {
	date := "{{ .Text.Grey }}(" + time.Now().Format(time.RFC822) + "){{ .Text.Reset }}"
	out := date + " {{ .Text.Task }}" + context.TaskName + "{{ .Text.Reset }} {{ .Text." + color + " }}" + component + " {{ .Text.Reset}}" + text
	t := translate(out, context)
	fmt.Print(t)
}

func cmdOK(text string, context *Context) {
	outputCommand("Command", "command", text, context)
}

func cmdFail(text string, context *Context) {
	outputCommand("Failure", "command", text, context)
}

func outputCommand(color, component, text string, context *Context) {
	date := "{{ .Text.Grey }}(" + time.Now().Format(time.RFC822) + "){{ .Text.Reset }}"
	out := date + " {{ .Text." + color + " }}" + text
	t := translate(out, context)
	fmt.Println(t)
}
