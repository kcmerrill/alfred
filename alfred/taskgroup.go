package alfred

import (
	"regexp"
	"strings"
	"sync"
)

// TaskGroup contains a task name and it's arguments
type TaskGroup struct {
	Name string
	Args []string
}

// ParseTaskGroup takes in a string, and parses it into a TaskGroup
func (t *Task) ParseTaskGroup(group string) []TaskGroup {
	tg := make([]TaskGroup, 0)
	group = strings.TrimSpace(group)

	if group == "" {
		return tg
	}

	// TODO: This is pretty terrible ... but until I get something better it stays
	// I need to research tokenizers
	if strings.Index(group, "\n") == -1 && !strings.Contains(group, "(") {
		// This means we have a regular space delimited list, probably
		tasks := strings.Split(group, " ")
		for _, task := range tasks {
			tg = append(tg, TaskGroup{Name: task, Args: []string{}})
		}
	} else {
		// mix and match here
		tasks := strings.Split(group, "\n")
		for _, task := range tasks {
			re := regexp.MustCompile(`(.*?)\((.*?)\)`)
			results := re.FindStringSubmatch(task)
			if len(results) == 0 {
				tg = append(tg, TaskGroup{Name: strings.TrimSpace(task), Args: []string{}})
			} else {
				args := strings.Split(results[2], ",")
				for idx, a := range args {
					// trim the extra whitespace
					args[idx] = strings.TrimSpace(a)
				}
				tg = append(tg, TaskGroup{Name: strings.TrimSpace(results[1]), Args: args})
			}
		}
	}
	return tg
}

func execTaskGroup(taskGroups []TaskGroup, task Task, context *Context, tasks map[string]Task) {
	for _, tg := range taskGroups {
		c := copyContex(context, translateArgs(tg.Args, context))
		NewTask(tg.Name, c, tasks)
	}
}

func goExecTaskGroup(taskGroups []TaskGroup, task Task, context *Context, tasks map[string]Task) {
	var wg sync.WaitGroup
	for _, tg := range taskGroups {
		wg.Add(1)
		go func(tg TaskGroup) {
			c := copyContex(context, translateArgs(tg.Args, context))
			NewTask(tg.Name, c, tasks)
			wg.Done()
		}(tg)
	}
	wg.Wait()
}
