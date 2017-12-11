package main

import . "github.com/kcmerrill/alfred/alfred"

func main() {
	tasks := make(map[string]Task)
	tasks["hello.world"] = Task{
		Serve: "8080",
	}
	NewTask("hello.world", InitialContext([]string{}), tasks)
}
