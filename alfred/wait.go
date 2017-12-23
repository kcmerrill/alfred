package alfred

import (
	"time"
)

func wait(task Task, context *Context, tasks map[string]Task) {
	if task.Wait == "" {
		return
	}

	dur, err := time.ParseDuration(translate(task.Wait, context))
	if err != nil {
		context.Ok = false
		outFail("waiting", "Unable to parse duration", context)
		return
	}

	outOK("wait", task.Wait, context)

	// get to waiting!
	<-time.After(dur)
}
