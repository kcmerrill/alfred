package alfred

import (
	"strings"
	"testing"

	event "github.com/kcmerrill/hook"
)

func TestCommandComponent(t *testing.T) {
	tasks := _sampleTasks()
	context := &Context{}

	// lets override output
	out := ""
	event.Register("output", func(text string, task Task, context *Context) {
		out += text + "\n"
	})

	command(tasks["ls"], context, tasks)
	if !strings.Contains(out, "main.go") {
		t.Fatalf("Expected to see main.go")
	}
}
