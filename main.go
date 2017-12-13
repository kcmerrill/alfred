package main

import . "github.com/kcmerrill/alfred/alfred"

func main() {
	tasks := make(map[string]Task)
	tasks["hello.world"] = Task{
		Summary: "Hello world! How are you!",
		Command: "whoami && sleep 1",
	}
	tasks["http.serve"] = Task{
		Serve: "8088",
	}

	NewTask("http.serve", InitialContext([]string{}), tasks)
}
