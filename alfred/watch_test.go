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
	go func() {
		<-time.After(2 * time.Second)
		execute("touch watch_test.go", ".")
	}()
	watch(task, context, tasks)
}
