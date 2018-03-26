package alfred

import (
	"testing"
)

func TestSimpleTaskParser(t *testing.T) {
	file, task := TaskParser("simple.task", "alfred:list")
	if file != "" {
		t.Fatalf("A simple task has a local file")
	}

	if task != "simple.task" {
		t.Fatalf("simple.task was passed in as a basic task")
	}
}
func TestDefaultTaskParser(t *testing.T) {
	file, task := TaskParser("", "alfred:list")
	if file != "" {
		t.Fatalf("No task(list) should be local")
	}

	if task != "alfred:list" {
		t.Fatalf("No task was passed, in a default task should have been returned")
	}
}

func TestRemoteTaskParser(t *testing.T) {
	file, task := TaskParser("/remote:new.task", "alfred:list")
	if file != "https://raw.githubusercontent.com/kcmerrill/alfred-tasks/master/remote.yml" {
		t.Fatalf("A remote task should return the master alfred github repository: " + file)
	}

	if task != "new.task" {
		t.Fatalf("The remote file and new.task should have been returned")
	}

	file, task = TaskParser("/remote", "alfred:list")
	if file != "https://raw.githubusercontent.com/kcmerrill/alfred-tasks/master/remote.yml" {
		t.Fatalf("A remote task should return the master alfred github repository: " + file)
	}

	if task != "alfred:list" {
		t.Fatalf("The default should have been returned")
	}
}

func TestHTTPTaskParser(t *testing.T) {
	file, task := TaskParser("http://someplace.com/whatever.yml:some.task", "alfred:list")
	if file != "http://someplace.com/whatever.yml" {
		t.Fatalf("Expected someplace.com to be returned")
	}

	if task != "some.task" {
		t.Fatalf("some.task should have been returned")
	}

	file, task = TaskParser("http://someplace.com/whatever.yml", "alfred:list")
	if file != "http://someplace.com/whatever.yml" {
		t.Fatalf("Expected someplace.com to be returned")
	}

	if task != "alfred:list" {
		t.Fatalf("Expected alfred:list to be returned when a task is not defined")
	}
}

func TestMagicTaskURL(t *testing.T) {
	tm := map[string]string{
		"/slack": "https://raw.githubusercontent.com/kcmerrill/alfred-tasks/master/slack.yml:",
		"slack":  "",
		"http://somesite.com/sometask.yml:sometask": "http://somesite.com/sometask.yml:",
	}

	for task, result := range tm {
		url := MagicTaskURL(task)
		if url != result {
			t.Fatal("expected", result, "actual", url)
		}
	}
}
