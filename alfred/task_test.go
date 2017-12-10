package alfred

import (
	"testing"
)

func _sampleTasks() map[string]Task {
	tasks := make(map[string]Task)
	tasks["hello.world"] = Task{
		Summary: "Hello world! How are you!",
		Command: "whoami && sleep 1",
		Ok:      "hello.world",
	}
	return tasks
}
func TestNewTask(t *testing.T) {
	tasks := _sampleTasks()
	NewTask("hello.world", InitialContext([]string{}), tasks)

	t.Fatalf("----stop")
}
