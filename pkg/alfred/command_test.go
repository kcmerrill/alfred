package alfred

import "testing"
import "os"

func TestCommandComponent(t *testing.T) {
	tasks := make(map[string]Task)
	task := Task{
		Dir:     "/tmp/alfred/command",
		Command: "touch myfile.txt",
	}
	c := InitialContext([]string{})
	c.Text.DisableFormatting = true
	commandC(task, c, tasks)

	// validate that it worked and that directory was created
	if _, exists := os.Stat("/tmp/alfred/command/myfile.txt"); exists != nil {
		t.Fatalf("file should have been created by command")
	}
}
