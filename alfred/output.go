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

func output(color, component, text string, context *Context) {
	date := "{{ .Text.Grey }}(" + time.Now().Format(time.RFC822) + "){{ .Text.Reset }}"
	out := "\r" + elapsed(context) + date + " {{ .Text.Task }}" + context.TaskName + "{{ .Text.Reset }} {{ .Text." + color + " }}" + component + " {{ .Text.Reset}}" + text + "{{ .Text.Reset }}"
	t := translate(out, context)
	fmt.Println(t)
}

func cmdOK(text string, context *Context) {
	outputCommand("Command", "command", text, context)
}

func cmdFail(text string, context *Context) {
	outputCommand("Failure", "command", text, context)
}

func outputCommand(color, component, text string, context *Context) {
	if text == "\r\n" || text == "\n" || text == "\r" || text == "" {
		formatting := "\033[1000D\033[K"
		if text == "\r" {
			formatting = "\033[1000D"
		}
		date := "{{ .Text.Grey }}(" + time.Now().Format(time.RFC822) + "){{ .Text.Reset }}"
		out := text + formatting + elapsed(context) + date + " {{ .Text." + color + " }}"
		t := translate(out, context)
		fmt.Print(t)
	} else {
		fmt.Print(text)
	}
	logger(text, context)
}

func elapsed(context *Context) string {
	return "{{ .Text.Grey }}[{{ .Text.Success }}" + padLeft(time.Since(context.Started).Round(time.Second).String(), 6, " ") + "{{ .Text.Grey }}] "
}
