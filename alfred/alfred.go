package alfred

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kcmerrill/alfred/remote"
	"github.com/kcmerrill/alfred/task"
	"gopkg.in/yaml.v2"
)

/* Our main object that holds the yaml file, tasks etc */
type Alfred struct {
	/* The contents of the yaml file */
	contents []byte
	/* Where the alfred.yml file was found */
	location string
	/* Variables */
	Vars map[string]string `yaml:"alfred.vars"`
	/* All of the tasks parsed from the yaml file */
	Tasks map[string]*task.Task `yaml:",inline"`
	/* Alfred remotes(private/public repos) */
	remote *remote.Remote
	/* Originating directory */
	dir string
	/* Alfred configuration */
	config struct {
		Remote map[string]string
	}
}

/* Setup our alfred struct */
func New() {
	a := new(Alfred)

	/* Grab our configuration */
	a.Config()

	/* Setup our remotes */
	a.remote = remote.New(a.config.Remote)

	/* Grab the current directory and save if off */
	a.dir, _ = os.Getwd()

	/* Try to find alfred.yml remotely(easy, needs a /) or find it locally */
	if a.findRemote() || a.findLocal() || a.findPrivate() {
		err := yaml.Unmarshal([]byte(a.contents), &a)
		if err == nil {
			/* Setup our aliases/promote commands */
			a.prepare()
			/* Ok, so we have instructions ... do we have a task to run? */
			a.findTask()
		} else {
			/* A problem with the yaml file */
			say("ERROR", err.Error())
			os.Exit(1)
		}
	} else {
		/* Bummer ... nothing found */
		say("ERROR", "Unable to find a job.")
		os.Exit(1)
	}
}

/* Depending on how you call alfred, depends on how it needs to find he task
   `alfred taskname` called taskname locally
   `alfred common/taskname` called taskname on the remote called common, finding a folder with taskname with an alfred.yml file in it.
   Remote files REQUIRE a "/"
*/
func (a *Alfred) findTask() {
	switch {
	/* Look locally, List tasks within its alfred.yml file */
	case len(os.Args) == 1:
		a.List()
		break
	/* Look remotely and list the tasks within it's alfred.yml file */
	case len(os.Args) == 2 && a.isRemote():
		a.List()
		break
	/* Called a local task */
	case len(os.Args) >= 2 && !a.isRemote():
		if a.isValidTask(os.Args[1]) && !a.Tasks[os.Args[1]].IsPrivate() {
			if !a.runTask(os.Args[1], os.Args[2:], false) {
				os.Exit(1)
			}
		} else {
			say(os.Args[1], "invalid task.")
			os.Exit(1)
		}
		break
	/* Called a remote task */
	case len(os.Args) >= 3 && a.isRemote():
		if a.isValidTask(os.Args[2]) && !a.Tasks[os.Args[2]].IsPrivate() {
			if !a.runTask(os.Args[2], os.Args[3:], false) {
				os.Exit(1)
			}
		} else {
			say(os.Args[2], "invalid task.")
			os.Exit(1)
		}
		break
	}
}

/* Meat and potatoes. Finds the task and runs it */
func (a *Alfred) runTask(task string, args []string, formatted bool) bool {
	/* Verify again it's a valid task */
	if !a.isValidTask(task) {
		say(task, "Invalid task.")
		return false
	}

	/* Infinite loop Used for the every command */
	for {
		taskok := true
		/* change to the original directory */
		err := os.Chdir(a.dir)
		if err != nil {
			fmt.Println(err.Error())
		}

		/* Lets prep it, and if it's bunk, lets see if we can pump out it's usage */
		if !a.Tasks[task].Prepare(args, a.Vars) {
			say(task+":error", "Missing argument(s).")
			return false
		}

		/* Lets change the directory if set */
		if a.Tasks[task].Dir != "" {
			if err := os.Chdir(a.Tasks[task].Dir); err != nil {
				if err := os.MkdirAll(a.Tasks[task].Dir, 0755); err != nil {
					say(task+":dir", "Invalid directory")
					return false
				} else {
					os.Chdir(a.Tasks[task].Dir)
				}
			}
		}

		/* Go through each of the modules ...
		- before command, docker stop for example
		*/
		for module, cmd := range a.Tasks[task].Modules {
			if !a.Tasks[task].RunCommand(os.Args[0]+" "+a.remote.ModulePath(module)+" "+cmd, task, formatted) {
				/* It failed :( */
				taskok = false
				break
			}
		}

		/* First, lets show the summary */
		if a.Tasks[task].Summary != "" {
			fmt.Println("")
			say(task, a.Tasks[task].Summary)
		}

		/* Lets execute the command if it has one, and add retry logic*/
		for x := 0; x < a.Tasks[task].Retry || x == 0; x++ {
			taskok = a.Tasks[task].RunCommand(a.Tasks[task].Command, task, formatted)
			if taskok {
				break
			}
		}

		/* The task failed ... */
		if !taskok {
			/* Failed? Lets run the failed tasks */
			for _, failed := range a.Tasks[task].FailedTasks() {
				if !a.runTask(failed, args, formatted) {
					break
				}
			}
		}

		/* Handle skips ... */
		if !taskok && a.Tasks[task].Skip {
			return false
		}

		/* Handle exits ... */
		if !taskok && a.Tasks[task].Exit != "" {
			if exitCode, err := strconv.Atoi(a.Tasks[task].Exit); err == nil {
				os.Exit(exitCode)
			}
			return false
		}

		/* Wait ... */
		if wait_duration, wait_err := time.ParseDuration(a.Tasks[task].Wait); wait_err == nil {
			<-time.After(wait_duration)
		}

		var wg sync.WaitGroup
		/* Do we have any tasks we need to run in parallel? */
		for _, t := range a.Tasks[task].MultiTask() {
			wg.Add(1)
			go func(t string, args []string) {
				defer wg.Done()
				a.runTask(t, args, true)
			}(t, args)
		}
		wg.Wait()

		/* Ok, we made it here ... Is this task a task group? */
		for _, t := range a.Tasks[task].TaskGroup() {
			if !a.runTask(t, args, formatted) {
				break
			}
		}

		/* Woot! Lets run the ok tasks */
		if taskok {
			for _, ok_tasks := range a.Tasks[task].OkTasks() {
				if !a.runTask(ok_tasks, args, formatted) {
					break
				}
			}

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

/* Ensure that the task exists */
func (a *Alfred) isValidTask(task string) bool {
	if _, exists := a.Tasks[task]; exists {
		return true
	}
	return false
}

/* The first argument MUST contain a "/" to be considered remote */
func (a *Alfred) isRemote() bool {
	if len(os.Args) >= 2 && strings.Contains(os.Args[1], "/") {
		return true
	}
	return false
}

/* Bounce around the web until we find something, or we don't .. */
func (a *Alfred) findRemote() bool {
	/* Make sure remote is a valid possibility */
	if a.isRemote() {
		remote, module := a.remote.Parse(os.Args[1])

		/* default to plain jane github */
		url := "https://raw.githubusercontent.com/" + os.Args[1] + "/master/alfred.yml"

		/* Does a remote exist? If so, we should use the remote syntax */
		if a.remote.Exists(remote) {
			url = a.remote.URL(remote, module)
		}

		/* try to fetch the alfred file */
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != 200 {
			say("error", "Unknown module "+os.Args[1])
			say("url", url)
			return true
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			/* We found something ... lets use it! */
			//a.contents = append(append(a.contents, []byte("\n")...), body...)
			a.contents = body
			a.location = url
			return true
		}
		return true
	}
	return false
}

/* This function will look locally for an alfred.yml file. First
starting in the working directory and then going back up the parent
directory until either an alfred.yml file is found, or you're in the
root directory */
func (a *Alfred) findLocal() bool {
	/* Grab the current directory */
	dir, err := os.Getwd()
	if err == nil {
		/* Just keep going ... */
		for {
			/* Keep going up a directory */
			if _, stat_err := os.Stat(dir + "/alfred.yml"); stat_err == nil {
				if contents, read_err := ioutil.ReadFile(dir + "/alfred.yml"); read_err == nil {
					/* Sweet. We found an alfred file. Lets save it off and return */
					//a.contents = append(append(a.contents, []byte("\n")...), contents...)
					a.contents = contents
					a.location = dir + "/alfred.yml"
					/* Be sure that we ar relative to where we found the config file */
					a.dir = dir
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
	/* We didn't find anything. /cry */
	return false
}

/* Look in the user's home directory for an alfred.yml file */
func (a *Alfred) findPrivate() bool {
	/* Grab the current directory */
	dir, err := os.Getwd()
	if err == nil {
		/* Set the current directory if all is good */
		a.dir = dir

		usr, err := user.Current()
		if err != nil {
			return false
		}

		privateFile := usr.HomeDir + "/.alfred/alfred.yml"

		if _, stat_err := os.Stat(privateFile); stat_err == nil {
			if contents, read_err := ioutil.ReadFile(privateFile); read_err == nil {
				a.contents = contents
				a.location = dir + "/alfred.yml"
				return true
			}
		}
	}

	/* We didn't find anything. /cry */
	return false
}

/* List out all of the available commands we can run */
func (a *Alfred) List() {
	/* Get/Sort list of tasks ... */
	t := []string{}
	for task, _ := range a.Tasks {
		t = append(t, task)
	}
	sort.Strings(t)

	promoted := false

	for _, which := range []string{"basic", "promoted"} {
		if which == "promoted" && promoted {
			fmt.Println("")
			fmt.Println("----")
			fmt.Println("")
		}
		for _, name := range t {
			task := a.Tasks[name]
			if task.IsAlias(name) || task.IsPrivate() {
				continue
			}

			if which == "basic" && strings.HasSuffix(name, "*") {
				promoted = true
				continue
			}

			if which == "promoted" && !strings.HasSuffix(name, "*") {
				continue
			}

			say(strings.Replace(name, "*", "", -1), task.Summary)

			if task.Alias != "" {
				fmt.Println("  ", "- Alias:", task.Alias)
			}

			if task.Usage != "" {
				fmt.Println("  ", "- Usage:", task.Usage)
			}

			if task.Tasks != "" {
				fmt.Println("  ", "- Tasks:", task.Tasks)
			}
		}
	}

}

/* If any commands have aliases, lets copy the tasks to their new names
   Also, if we have an astrick in the name, lets promote it
*/
func (a *Alfred) prepare() {
	for name, task := range a.Tasks {
		/* Does this task have an alias? If so, lets create it! */
		if len(task.Aliases()) > 0 {
			for _, alias := range task.Aliases() {
				a.Tasks[alias] = task
			}
		}

		/* Should this task be promoted? */
		if strings.HasSuffix(name, "*") {
			a.Tasks[strings.Replace(name, "*", "", -1)] = task
		}
	}
}

/* Alfred speaks! */
func say(task, msg string) {
	fmt.Println("["+task+"]", msg)
}
