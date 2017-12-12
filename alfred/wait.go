package alfred

import "time"

func wait(task Task, context *Context, tasks map[string]Task) {
	if task.Wait == "" {
		return
	}

	dur, err := time.ParseDuration(task.Wait)
	if err != nil {
		return
	}

	// get to waiting!
	<-time.After(dur)
}
