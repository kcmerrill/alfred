package alfred

import "testing"

func tasksTestHelper() Tasks {
	raw := make(map[string]Task)
	raw["hello.world"] = Task{
		Summary: "Hello world! How are you!",
		Command: "whoami",
	}
	return Tasks{
		Raw: raw,
	}
}

func TestTasksGet(t *testing.T) {
	tasks := tasksTestHelper()
	if task, exists := tasks.Get("hello.world"); !exists {
		t.Fatalf("Unable to find the task 'hello.world'")
	} else {
		if task.Command != "whoami" {
			t.Fatalf("The incorrect task was found. Expected 'whoami' as the command")
		}
	}
}
