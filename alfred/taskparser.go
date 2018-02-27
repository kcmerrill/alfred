package alfred

import (
	"fmt"
	"strings"
)

// TaskParser returns a url and a task(default if necessary)
func TaskParser(task, defaultTask string) (string, string) {
	// does it start with http?
	if strings.HasPrefix(task, "http") && strings.Contains(task, "://") {
		// we have to get the http: colon out of the way :(
		bits := strings.SplitN(task, ":", 3)
		url := strings.Join(bits[0:2], ":")
		if len(bits) >= 3 {
			// alright, so we have tasks and args ...
			return url, bits[2]
		}
		return url, defaultTask
	}

	// does it start with a /? Aka a remote task?
	if strings.HasPrefix(task, "/") {
		bits := strings.Split(task, ":")
		url := fmt.Sprintf("https://raw.githubusercontent.com/kcmerrill/alfred-tasks/master%s.yml", bits[0])

		if len(bits) >= 2 {
			return url, bits[1]
		}
		return url, defaultTask
	}

	// lets check if this is a github file
	if strings.Contains(task, "/") {
		bits := strings.Split(task, ":")
		url := fmt.Sprintf("https://raw.githubusercontent.com/%s/master/alfred.yml", bits[0])
		if len(bits) >= 2 {
			return url, bits[1]
		}
		return url, defaultTask
	}

	if task == "" {
		return "", defaultTask
	}

	if strings.HasPrefix(task, "!") {
		return "", "!exec.command"
	}

	// alright, so it's not a url, it's not a github repo, it must be just a regular local task
	return "", task
}
