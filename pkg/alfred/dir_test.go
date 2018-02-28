package alfred

import (
	"testing"
)

func TestTaskDir(t *testing.T) {
	task := Task{
		Dir: "/tmp/alfred/dir-test-{{ .Text.Failure }}",
	}

	c := InitialContext([]string{})
	c.Text.Failure = "red"

	dir, ok := task.dir(c)

	if dir != "/tmp/alfred/dir-test-red" {
		t.Fatalf("Template not working with task.dir()")
	}

	if !ok {
		t.Fatalf("Ability to create directories broken")
	}
}
