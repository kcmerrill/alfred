package alfred

import (
	"testing"
	"time"
)

func TestWaitComponent(t *testing.T) {
	task := Task{
		Wait: "2s",
	}
	context := &Context{}
	start := time.Now().Unix()
	wait(task, context, _sampleTasks())
	finish := time.Now().Unix()
	if finish-2 != start {
		t.Fatalf("wait() did not wait 2 seconds. #sadpanda")
	}
}
