package alfred

import (
	"fmt"
	"os"
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

	// validate that false stopped the commands component
	if _, exists := os.Stat("/tmp/alfred/commands/myfile.txt"); exists != nil {
		t.Fatalf("file should have been created by commands. ONLY BECAUSE NO EXIT CODE!")
	}

	// OK, controversy ... lets talk
	// This should fail the task, but it shouldn't exit if no exit isn't set
	task = Task{
		Dir:      "/tmp/alfred/commands",
		Commands: "false\ntouch myfile2.txt",
		ExitCode: 42,
	}

	c = InitialContext([]string{})
	commands(task, c, tasks)
	if _, exists := os.Stat("/tmp/alfred/commands/myfile2.txt"); exists != nil {
		fmt.Println(exists.Error())
		t.Fatalf("myfile2.txt should have been created this time")
	}
}
