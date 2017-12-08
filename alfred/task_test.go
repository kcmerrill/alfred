package alfred

import "testing"

func TestNewTask(t *testing.T) {
	tasks := tasksTestHelper()
	eventStream := NewEventStream(&EventStream{})
	NewTask("hello.world", Context{}, eventStream, tasks)
	t.Fatalf("Something failed")
}
