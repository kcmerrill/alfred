package alfred

import "testing"

func TestNewTask(t *testing.T) {
	tasks := tasksTestHelper()
	NewTask("hello.world", Context{}, tasks)
	t.Fatalf("Something failed")
}
