package alfred

import (
	"testing"
	"time"
)

func TestWatchComponent(t *testing.T) {
	task := Task{
		Watch: ".*?go$",
	}
	context := InitialContext([]string{})
	tasks := make(map[string]Task)
	changes := make(chan bool)
	go func(changes chan bool) {
		evaluate("touch watch_test.go", ".")
		watch(task, context, tasks)
		changes <- true
	}(changes)

	select {
	case <-changes:
	case <-time.After(2 * time.Second):
		// timeout ... boo!
		t.Fatalf("Watch failed ...")
	}

	<-time.After(time.Second * 3)
}
