package alfred

import (
	"fmt"
)

func output(text string, task Task, context *Context) {
	t := task.Template(text, context)
	fmt.Println(t)
}
