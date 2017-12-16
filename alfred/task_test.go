package alfred

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTask(t *testing.T) {
	tasks := make(map[string]Task)
	tasks["hello.world"] = Task{
		Summary: "Hello world! How are you!",
		Command: "whoami && sleep 1",
	}
	NewTask("hello.world", InitialContext([]string{}), tasks)
}


func TestTaskIsPrivate(t *testing.T) {
	task := Task{
		Description: "A non private task here",
	}
	assert.Equal(t, false, task.IsPrivate(), "Task should be private")

	task = Task{
		Usage: "Usage goes here",
	}
	assert.Equal(t, false, task.IsPrivate(), "Task should be private")

	task = Task{}
	assert.Equal(t, true, task.IsPrivate(), "Task should be private")
}
