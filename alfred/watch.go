package alfred

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

func watch(task Task, context *Context, tasks map[string]Task) {
	if task.Watch == "" {
		return
	}
	dir, _ := task.dir(context)
	outOK("watching", dir, context)
	for {
		matched := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
			if f.ModTime().After(time.Now().Add(-2 * time.Second)) {
				m, _ := regexp.Match(translate(task.Watch, context), []byte(path))
				if m {
					// If not a match ...
					return fmt.Errorf(path)
				}
			}
			// continue on
			return nil
		})

		if matched != nil {
			// seems weird, but we are passing back an error
			outOK("modified", matched.Error(), context)
			break
		} else {
			<-time.After(time.Second)
		}
	}
}
