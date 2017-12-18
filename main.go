package main

import (
	"os"

	. "github.com/kcmerrill/alfred/alfred"
)

func main() {
	tasks := make(map[string]Task)
	task, args := CLI(os.Args)
	context := &Context{
		Args: args,
	}
	NewTask(task, context, tasks)
}
