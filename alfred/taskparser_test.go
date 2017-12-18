package alfred

import (
	"testing"
)

func TestSimpleTaskParser(t *testing.T) {
	file, task := TaskParser("simple.task", ":list")
	if file != ":local" {
		t.Fatalf("A simple task has a local file")
	}

	if task != "simple.task" {
		t.Fatalf("simple.task was passed in as a basic task")
	}
}
func TestDefaultTaskParser(t *testing.T) {
	file, task := TaskParser("", ":list")
	if file != ":local" {
		t.Fatalf("No task(list) should be local")
	}

	if task != ":list" {
		t.Fatalf("No task was passed, in a default task should have been returned")
	}
}

func TestRemoteTaskParser(t *testing.T) {
	file, task := TaskParser("/remote:new.task", ":list")
	if file != "https://raw.githubusercontent.com/kcmerrill/alfred/master/remote/remote.yml" {
		t.Fatalf("A remote task should return the master alfred github repository")
	}

	if task != "new.task" {
		t.Fatalf("The remote file and new.task should have been returned")
	}

	file, task = TaskParser("/remote", ":default")
	if file != "https://raw.githubusercontent.com/kcmerrill/alfred/master/remote/remote.yml" {
		t.Fatalf("A remote task should return the master alfred github repository")
	}

	if task != ":default" {
		t.Fatalf("The default should have been returned")
	}
}

func TestHTTPTaskParser(t *testing.T) {
	file, task := TaskParser("http://someplace.com/whatever.yml:some.task", ":list")
	if file != "http://someplace.com/whatever.yml" {
		t.Fatalf("Expected someplace.com to be returned")
	}

	if task != "some.task" {
		t.Fatalf("some.task should have been returned")
	}

	file, task = TaskParser("http://someplace.com/whatever.yml", ":list")
	if file != "http://someplace.com/whatever.yml" {
		t.Fatalf("Expected someplace.com to be returned")
	}

	if task != ":list" {
		t.Fatalf("Expected :list to be returned when a task is not defined")
	}
}
