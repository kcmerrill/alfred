package main

import . "github.com/kcmerrill/alfred/alfred"

func main() {
	tasks := make(map[string]Task)
	tasks["hello.world"] = Task{
		//Summary: "Hello world! How are you!",
		Command: "bleh && sleep 1",
		Wait:    "10s",
	}
	NewTask("hello.world", InitialContext([]string{}), tasks)
}
