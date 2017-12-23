package alfred

import (
	"time"
)

func every(task Task, context *Context, tasks map[string]Task) {
	e := task.Every

	// override every if we are watching
	if task.Watch != "" {
		e = "1s"
	}

	// convert task.Every into a duration
	if e == "" {
		return
	}

	dur, err := time.ParseDuration(e)
	if err != nil {
		return
	}

	outOK("every", e, context)
	// pause ...
	<-time.After(dur)
	NewTask(context.TaskName, context, tasks)
}
