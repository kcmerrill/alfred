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

func TestTaskTemplate(t *testing.T) {
	context := &Context{
		Text: TextConfig{
			Success:     "green",
			SuccessIcon: "checkmark",
			Failure:     "red",
			FailureIcon: "x",
		},
	}

	task := Task{}

	if task.Template("{{ .Text.Success }}", context) != "green" {
		t.Fatalf(".Text.Success should be green")
	}

	if task.Template("{{ .Text.SuccessIcon }}", context) != "checkmark" {
		t.Fatalf(".Text.SuccessIcon should be icon")
	}

	if task.Template("{{ .Text.Failure }}", context) != "red" {
		t.Fatalf(".Text.Failure should be red")
	}

	assert.Equal(t, task.Template("{{ .Text.FailureIcon }}", context), "x", "Should be X")
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
