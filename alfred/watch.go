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
	for {
		matched := filepath.Walk(task.Dir, func(path string, f os.FileInfo, err error) error {
			if f.ModTime().After(time.Now().Add(-2 * time.Second)) {
				m, _ := regexp.Match(task.Watch, []byte(path))
				if m {
					// If not a match ...
					return nil
				}
			}
			return fmt.Errorf("No matches found")
		})
		if matched != nil {
			break
		} else {
			<-time.After(1 * time.Second)
		}
	}
}
