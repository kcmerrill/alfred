package alfred

import (
	"time"
)

func every(task Task, context *Context) {
	// convert task.Every into a duration
	if task.Every == "" {
		return
	}

	dur, err := time.ParseDuration(task.Every)
	if err != nil {
		return
	}

	// pause ...
	<-time.After(dur)
	NewTask(context.Name, context, tasks)
}
