package main

import . "github.com/kcmerrill/alfred/alfred"

func main() {
	tasks := make(map[string]Task)
	tasks["hello.world"] = Task{
		Summary: "Hello world! How are you!",
		Watch:   ".*?go$",
		Command: "whoami",
	}
	NewTask("hello.world", InitialContext([]string{}), tasks)
}
