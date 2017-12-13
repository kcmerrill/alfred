package alfred

import (
	"testing"
)

func TestCommandComponent(t *testing.T) {
	tasks := make(map[string]Task)
	context := &Context{}
	command(tasks["ls"], context, tasks)
	// todo: figure out stdout capture
}
