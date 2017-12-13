package alfred

import (
	"fmt"
	"testing"
)

func _testSampleTasks() map[string]Task {
	tasks := make(map[string]Task)
	tasks["hello.world"] = Task{
		Summary: "Hello world! How are you!",
		Command: "whoami && sleep 1",
	}
	tasks["ls"] = Task{
		Summary: "ls current directory",
		Dir:     "../",
		Command: "ls -alh",
	}
	tasks["http.serve"] = Task{
		Serve: "8088",
	}
	return tasks
}

func TestNewTask(t *testing.T) {
	tasks := _testSampleTasks()
	NewTask("hello.world", InitialContext([]string{}), tasks)
}

func TestParseTaskGroup(t *testing.T) {
	task := Task{}
	tg := task.ParseTaskGroup(" task.one task.two task.three ")
	if tg[0].Name != "task.one" {
		t.Fatalf("tg should contain task.one")
	}

	if tg[1].Name != "task.two" {
		t.Fatalf("tg should contain task.two")
	}

	if tg[2].Name != "task.three" {
		t.Fatalf("tg should contain task.three")
	}

	tgWArgs := task.ParseTaskGroup("task.one\n task.two(arg1, arg2 , arg3, {{ arg4 }}) \ntask.two\ntask.three(arg1)")
	fmt.Println(tgWArgs)

	if tgWArgs[0].Name != "task.one" {
		t.Fatalf("task with args name not set")
	}

	if len(tgWArgs[0].Args) != 0 {
		t.Fatalf("task with args, the args should be empty")
	}

	fmt.Println("1.Name", tgWArgs[1].Name)
	if tgWArgs[1].Name != "task.two" {
		t.Fatalf("task.two was expected, with arguments")
	}

	if len(tgWArgs[1].Args) != 4 {
		t.Fatalf("task.two should have 4 arguments")
	}

	if tgWArgs[1].Args[1] != "arg2" {
		t.Fatalf("Expected 'arg2' as argument with no white space")
	}
}
