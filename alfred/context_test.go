package alfred

import "testing"

func TestInitialContext(t *testing.T) {
	c := InitialContext([]string{"one", "two"})

	if c.TaskName == "" {
		t.Fatalf("The taskname shouldn't be blank")
	}

	if c.Text.Failure == "" {
		t.Fatalf("The text object should not be empty")
	}

	if c.Args[1] != "two" {
		t.Fatalf("Default args should be set")
	}
}
