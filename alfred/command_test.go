package alfred

import (
	"testing"
)

func TestCommandComponent(t *testing.T) {
	tasks := _testSampleTasks()
	context := &Context{}
	command(tasks["ls"], context, tasks)
	// todo: figure out stdout capture
}
