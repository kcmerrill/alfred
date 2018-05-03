package alfred

import (
	"fmt"
	"time"
)

func outOK(component, text string, context *Context) {
	output("Success", component, text, context)
}

func outFail(component, text string, context *Context) {
	output("Failure", component, text, context)
}

func outWarn(component, text string, context *Context) {
	output("Warning", component, text, context)
}

func outArgs(component, text string, context *Context) {
	output("Args", component, text, context)
}

func outPrefix(color, component, text string, context *Context) string {
	date := "{{ .Text.Grey }}(" + time.Now().Format(time.RFC822) + "){{ .Text.Reset }}"
	out := elapsed(context) + date + "  [" + context.rootDir + "] {{ .Text.Task }}" + context.TaskName + "{{ .Text.Reset }} {{ .Text." + color + " }}" + component + " {{ .Text.Reset}}" + text + "{{ .Text.Reset }}"
	return out
}

func output(color, component, text string, context *Context) {
	if context.Text.DisableFormatting {
		return
	}
	out := outPrefix(color, component, text, context)
	t := translate(out, context)

	fmt.Fprintln(context.Out, t)
}

func outputPrompt(color, component, text string, context *Context) {
	out := outPrefix(color, component, text, context)
	t := translate(out, context)
	fmt.Fprintln(context.Out, t)
}

func cmdOK(text string, context *Context) {
	outputCommand("Command", "command", text, context)
}

func cmdFail(text string, context *Context) {
	outputCommand("Failure", "command", text, context)
}

func outputCommand(color, component, text string, context *Context) {
	if !context.Text.DisableFormatting {
		date := "{{ .Text.Grey }}(" + time.Now().Format(time.RFC822) + "){{ .Text.Reset }}"
		out := elapsed(context) + date + " {{ .Text." + color + " }} " + text
		t := translate(out, context)
		fmt.Fprintln(context.Out, t)
	} else {
		fmt.Fprintln(context.Out, text)
	}

	logger(text, context)
}

func elapsed(context *Context) string {
	return "{{ .Text.Grey }}[{{ .Text.Success }}" + padLeft(time.Since(context.Started).Round(time.Second).String(), 8, " ") + "{{ .Text.Grey }}] "
}
