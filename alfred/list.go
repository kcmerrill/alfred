package alfred

import "fmt"

func list(tasks map[string]Task) {
	for label, task := range tasks {
		fmt.Println(translate(label+" : "+task.Summary, emptyContext()))
	}
}
