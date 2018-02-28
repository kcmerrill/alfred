package alfred

import (
	"testing"
	"time"
)

func TestWaitComponent(t *testing.T) {
	task := Task{
		Wait: "2s",
	}
	context := InitialContext([]string{})
	start := time.Now().Unix()
	tasks := make(map[string]Task)
	wait(task, context, tasks)
	finish := time.Now().Unix()
	if finish-2 != start {
		t.Fatalf("wait() did not wait 2 seconds. #sadpanda")
	}
}
