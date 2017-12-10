package alfred

import "fmt"

func taskGroup(group string, task Task, context *Context, tasks map[string]Task) {
	tgs := task.ParseTaskGroup(group)
	for _, tg := range tgs {
		fmt.Println("taskgroup", tg.Name)
		NewTask(tg.Name, InitialContext([]string{}), tasks)
	}
}
