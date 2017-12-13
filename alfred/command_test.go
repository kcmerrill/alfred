package alfred

import (
	"testing"
)

func TestCommandComponent(t *testing.T) {
	tasks := _testSampleTasks()
	context := &Context{}
	command(tasks["ls"], context, tasks)
	//t.Fatalf("Expected to see main.go")
}
