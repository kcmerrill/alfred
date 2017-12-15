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
	for {
		fmt.Println("watching", "---", dir, "---")
		matched := filepath.Walk("/tmp/", func(path string, f os.FileInfo, err error) error {
			fmt.Println("path", path)
			fmt.Println(f.Name())
			if f.ModTime().After(time.Now().Add(-2 * time.Second)) {
				fmt.Println("Name:", f.Name())
				m, _ := regexp.Match(task.Watch, []byte(path))
				if m {
					fmt.Println("matches:", m)
					// If not a match ...
					return nil
				}
			}
			fmt.Println("Nothing")
			return fmt.Errorf("No matches found")
		})
		<-time.After(time.Second)
		fmt.Println(matched)
		//if matched != nil {
		//fmt.Println(".... lets try again")
		//<-time.After(1 * time.Second)
		//} else {
		//fmt.Println("good to go ...")
		//break
		//}
	}
}
