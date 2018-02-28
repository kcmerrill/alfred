package alfred

import (
	"testing"
)

func TestDefaults(t *testing.T) {
	task := Task{
		Defaults: []string{"one", "two", "three"},
	}
	context := InitialContext([]string{"one", "two"})
	context.Text.DisableFormatting = true
	tasks := make(map[string]Task)

	defaults(task, context, tasks)

	if context.Args[2] != "three" {
		t.Fatalf("default for arg #3 is three")
	}
}
