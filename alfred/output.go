package alfred

import (
	"fmt"
)

func output(text string, task Task, context *Context) {
	t := task.Template(text, context)
	if !context.Silent {
		fmt.Println(t)
	}
}
