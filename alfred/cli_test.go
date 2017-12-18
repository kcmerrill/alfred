package alfred

import "testing"

func TestEmptyCLI(t *testing.T) {
	task, args := CLI([]string{"alfred"})
	if task != "" {
		t.Fatalf("Expected: an empty task")
	}

	if len(args) >= 1 {
		t.Fatalf("No args were passed in. Should be len(0)")
	}
}

func TestCLI(t *testing.T) {
	task, args := CLI([]string{"alfred", "/tester:one", "two"})
	if task != "/tester:one" {
		t.Fatalf("Expected /tester:one")
	}

	if args[0] != "two" {
		t.Fatalf("Expected an argument to be passed in.")
	}
}
