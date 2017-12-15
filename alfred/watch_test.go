package alfred

import (
	"testing"
	"time"
)

func TestWatchComponent(t *testing.T) {
	task := Task{
		Watch: ".*?go$",
	}
	context := &Context{}
	tasks := make(map[string]Task)
	changes := make(chan bool)
	go func(changes chan bool) {
		watch(task, context, tasks)
		changes <- true
	}(changes)

	select {
	case <-changes:
		t.Fatalf("here ...")
	case <-time.After(2 * time.Second):
		// timeout ... boo!
		t.Fatalf("Watch failed ...")
	}
}
