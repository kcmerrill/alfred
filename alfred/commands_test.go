package alfred

import (
	"testing"
)

func TestCommandsComponent(t *testing.T) {
	tasks := make(map[string]Task)
	task := Task{
		Dir:      "/tmp/alfred/commands",
		Commands: "false\ntouch myfile.txt",
	}

	c := InitialContext([]string{})
	commands(task, c, tasks)

	if c.Ok {
		// it shouldn't be ok!
		t.Fatalf("the false command should have failed the task!")
	}
}
