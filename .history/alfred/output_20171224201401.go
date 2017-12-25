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

func outArgs(component, text string, context *Context) {
	output("Args", component, text+"\n", context)
}

func output(color, component, text string, context *Context) {
	date := "{{ .Text.Grey }}(" + time.Now().Format(time.RFC822) + "){{ .Text.Reset }}"
	out := elapsed(context) + date + " {{ .Text.Task }}" + context.TaskName + "{{ .Text.Reset }} {{ .Text." + color + " }}" + component + " {{ .Text.Reset}}" + text + "{{ .Text.Reset }}"
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
	out := elapsed(context) + date + " {{ .Text." + color + " }}" + text
	t := translate(out, context)
	fmt.Println(t)
	logger(text+"\n", context)
}

func elapsed(context *Context) string {
	return "{{ .Text.Grey }}[{{ .Text.Success }}" + padLeft(time.Since(context.Started).Round(time.Second).String(), 3, " ") + "{{ .Text.Grey }}] "
}
