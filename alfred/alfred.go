package alfred

import (
	"fmt"
	"github.com/kcmerrill/alfred/remote"
	"github.com/kcmerrill/alfred/task"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

type Alfred struct {
	contents []byte
	location string
	Tasks    map[string]*task.Task `yaml:",inline"`
	remote   *remote.Remote
	dir      string
}

func New() {
	a := new(Alfred)
	a.remote = remote.New()

	a.dir, _ = os.Getwd()

	if a.findRemote() || a.findLocal() {
		err := yaml.Unmarshal([]byte(a.contents), &a)
		if err == nil {
			/* Ok, so we have instructions ... do we have a task to run? */
			a.findTask()
		} else {
			say("ERROR", err.Error())
			os.Exit(1)
		}
	} else {
		/* Bummer ... */
		say("ERROR", "Unable to find a job.")
		os.Exit(1)
	}
}

func (a *Alfred) findTask() {
	switch {
	case len(os.Args) == 1:
		a.List()
		break
	case len(os.Args) == 2 && a.isRemote():
		a.List()
		break
	case len(os.Args) >= 2 && !a.isRemote():
		if a.isValidTask(os.Args[1]) {
			if !a.runTask(os.Args[1], os.Args[2:]) {
				os.Exit(1)
			}
		} else {
			say(os.Args[1], "invalid task.")
			os.Exit(1)
		}
		break
	case len(os.Args) >= 3 && a.isRemote():
		if a.isValidTask(os.Args[2]) {
			if !a.runTask(os.Args[2], os.Args[3:]) {
				os.Exit(1)
			}
		} else {
			say(os.Args[2], "invalid task.")
			os.Exit(1)
		}
		break
	}
}

func (a *Alfred) runTask(task string, args []string) bool {
	/* Verify again it's a valid task */
	if !a.isValidTask(task) {
		say(task, "Invalid task.")
		return false
	}

	for {
		/* change to the original directory */
		err := os.Chdir(a.dir)
		if err != nil {
			fmt.Println(err.Error())
		}

		/* First, lets show the summary */
		if a.Tasks[task].Summary != "" {
			fmt.Println("")
			say(task, a.Tasks[task].Summary)
		}

		/* Lets prep it, and if it's bunk, lets see if we can pump out it's usage */
		if !a.Tasks[task].Prepare(args) {
			say(task, "ERROR: Missing argument(s).")
			return false
		}

		/* Lets change the directory if set */
		if a.Tasks[task].Dir != "" {
			if err := os.Chdir(a.Tasks[task].Dir); err != nil {
				say(task+":dir", "Invalid directory")
				return false
			}
		}

		/* Lets execute the command if it has one */
		if !a.Tasks[task].RunCommand(a.Tasks[task].Command, args) {
			/* Failed? Lets run the failed tasks */
			for _, failed := range a.Tasks[task].FailedTasks() {
				if !a.runTask(failed, args) {
					return false
				}
			}
			return false
		} else {
			/* Woot! Lets run the ok tasks */
			for _, ok_tasks := range a.Tasks[task].OkTasks() {
				if !a.runTask(ok_tasks, args) {
					return false
				}
			}
		}

		/* Go through each of the modules ... */
		for module, cmd := range a.Tasks[task].Modules {
			if !a.Tasks[task].RunCommand(os.Args[0]+" "+a.remote.ModulePath(module)+" "+cmd, args) {
				/* It failed :( */
				return false
			}
		}

		/* Ok, we made it here ... Is this task a task group? */
		for _, t := range a.Tasks[task].TaskGroup() {
			if !a.runTask(t, args) {
				return false
			}
		}

		/* Wait ... */
		if wait_duration, wait_err := time.ParseDuration(a.Tasks[task].Wait); wait_err == nil {
			<-time.After(wait_duration)
		}

		/* Do we need to break or should we keep going? */
		if a.Tasks[task].Every != "" {
			if every_duration, every_err := time.ParseDuration(a.Tasks[task].Every); every_err == nil {
				<-time.After(every_duration)
			} else {
				break
			}
		} else {
			break
		}
	}
	return true
}

func (a *Alfred) isValidTask(task string) bool {
	if _, exists := a.Tasks[task]; exists {
		return true
	}
	return false
}

func (a *Alfred) isRemote() bool {
	if len(os.Args) >= 2 && strings.Contains(os.Args[1], "/") {
		return true
	}
	return false
}

func (a *Alfred) findRemote() bool {
	if a.isRemote() {
		remote, module := a.remote.Parse(os.Args[1])
		/* default to plain jane github */
		url := "https://raw.githubusercontent.com/" + os.Args[1] + "/master/alfred.yml"
		if a.remote.Exists(remote) {
			url = a.remote.URL(remote, module)
		}

		/* try to fetch now */
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != 200 {
			say("ERROR", "Unknown module "+os.Args[1])
			return true
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			a.contents = body
			a.location = url
			return true
		}
		return true
	}
	return false
}

func (a *Alfred) findLocal() bool {
	dir, err := os.Getwd()
	if err == nil {
		for {
			/* Keep going up a directory */
			if _, stat_err := os.Stat(dir + "/alfred.yml"); stat_err == nil {
				if contents, read_err := ioutil.ReadFile(dir + "/alfred.yml"); read_err == nil {
					a.contents = contents
					a.location = dir + "/alfred.yml"
					/* Be sure that we ar relative to where we found the config file */
					os.Chdir(dir)
					return true
				}
			}
			dir = path.Dir(dir)
			if dir == "/" {
				/* We've gone too far ... */
				break
			}
		}
	}
	return false
}

func (a *Alfred) List() {
	fmt.Println()
	for name, task := range a.Tasks {
		say(name, task.Summary)

		if task.Usage != "" {
			fmt.Println("  ", "Usage:", task.Usage)
		}

		if task.Tasks != "" {
			fmt.Println("  ", "Tasks:", task.Tasks)
		}

	}
}

func say(task, msg string) {
	fmt.Println("["+task+"]", msg)
}
