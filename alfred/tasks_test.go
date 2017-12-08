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

}
