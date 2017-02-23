package alfred

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/kcmerrill/alfred/remote"
	"github.com/kcmerrill/alfred/task"
	"gopkg.in/yaml.v2"
)

// Alfred is our main object that holds the yaml file, tasks etc
type Alfred struct {
	// The contents of the yaml file
	contents []byte
	// Where the alfred.yml file was found
	location string
	// Variables
	Vars map[string]string `yaml:"alfred.vars"`
	// All of the tasks parsed from the yaml file
	Tasks map[string]*task.Task `yaml:",inline"`
	// Alfred remotes(private/public repos)
	remote *remote.Remote
	// Originating directory
	dir string
	// Alfred configuration
	config struct {
		Remote map[string]string
	}
	// Task/Step Number
	taskNumber int
	// Mutex
	lock *sync.Mutex
}

// New creates and sets up our alfred struct
func New() {
	a := new(Alfred)

	// Grab our configuration
	a.Config()

	// Setup our remotes
	a.remote = remote.New(a.config.Remote)

	// Grab the current directory and save if off
	a.dir, _ = os.Getwd()

	// Setup our mutex
	a.lock = &sync.Mutex{}

	// Try to find alfred.yml remotely(easy, needs a /) or find it locally
	if a.findRemote() || a.findLocal() {
		err := yaml.Unmarshal([]byte(a.contents), &a)
		if err == nil {
			// Setup our aliases/promote commands
			a.prepare()
			// Ok, so we have instructions ... do we have a task to run?
			if !a.findTask() {
				os.Args = append(os.Args[:1], append([]string{"default"}, os.Args[1:]...)...)
				if !a.findTask() {
					say("ERROR", "Invalid task.")
					os.Exit(1)
				}
			}
		} else {
			// A problem with the yaml file
			say("ERROR", err.Error())
			os.Exit(1)
		}
	} else {
		// Bummer ... nothing found
		say("ERROR", "Unable to find a job.")
		os.Exit(1)
	}
}

// findTask will determine if the task is local or remote, or if the user even passed one in
// Depending on how you call alfred, depends on how it needs to find he task
// `alfred taskname` called taskname locally
// `alfred common/taskname` called taskname on the remote called common, finding a folder with taskname with an alfred.yml file in it.
// Remote files REQUIRE a "/"
func (a *Alfred) findTask() bool {
	switch {
	// Look locally, List tasks within its alfred.yml file
	case len(os.Args) == 1:
		a.List()
		break
	// Look remotely and list the tasks within it's alfred.yml file
	case len(os.Args) == 2 && a.isRemote():
		a.List()
		break
	// Called a local task
	case len(os.Args) >= 2 && !a.isRemote():
		if a.isValidTask(os.Args[1]) && !a.Tasks[os.Args[1]].IsPrivate() {
			if !a.runTask(os.Args[1], os.Args[2:], false) {
				return false
			}
		} else {
			return false
		}
		break
	// Called a remote task
	case len(os.Args) >= 3 && a.isRemote():
		if a.isValidTask(os.Args[2]) && !a.Tasks[os.Args[2]].IsPrivate() {
			if !a.runTask(os.Args[2], os.Args[3:], false) {
				return false
			}
		} else {
			say(os.Args[2], "invalid task.")
			return false
		}
		break
	}
	return true
}

// runTask runs the necessary commands required for each task
func (a *Alfred) runTask(task string, args []string, formatted bool) bool {
	// Verify again it's a valid task
	if !a.isValidTask(task) {
		say(task, "Invalid task.")
		return false
	}

	// Infinite loop Used for the every command
	for {
		// Increment the taskNumber
		a.lock.Lock()
		a.taskNumber = a.taskNumber + 1
		a.lock.Unlock()

		// Run our setup tasks
		for _, s := range a.Tasks[task].SetupTasks() {
			if !a.runTask(s, args, formatted) {
				break
			}
		}

		taskok := true

		// change to the original directory
		err := os.Chdir(a.dir)
		if err != nil {
			fmt.Println(err.Error())
		}

		// Lets prep it, and if it's bunk, lets see if we can pump out it's usage
		if !a.Tasks[task].Prepare(args, a.Vars, a.taskNumber) {
			say(task+":error", "Missing argument(s).")
			return false
		}

		// Lets change the directory if set
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

		// We watching for files?
		if a.Tasks[task].Watch != "" {
			// Regardless of what's going on, lets set every to 1s
			a.Tasks[task].Every = "1s"
			for {
				matched := filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
					if f.ModTime().After(time.Now().Add(-2 * time.Second)) {
						m, _ := regexp.Match(a.Tasks[task].Watch, []byte(path))
						if m {
							// If not a match ...
							return errors.New("")
						}
					}
					return nil
				})
				if matched != nil {
					break
				} else {
					<-time.After(1 * time.Second)
				}
			}
		}

		// Go through each of the modules ...
		// before command, docker stop for example
		for module, cmd := range a.Tasks[task].Modules {
			if !a.Tasks[task].RunCommand(os.Args[0]+" "+a.remote.ModulePath(module)+" "+cmd, task, formatted) {
				// It failed :(
				taskok = false
				break
			}
		}

		// First, lets show the summary
		if a.Tasks[task].Summary != "" {
			fmt.Println("")
			say(task, a.Tasks[task].Summary)
		}

		// Test ...
		if a.Tasks[task].TestF(a.Tasks[task].Test) {
			// Lets execute the command if it has one, and add retry logic
			for x := 0; x < a.Tasks[task].Retry || x == 0; x++ {
				taskok = a.Tasks[task].RunCommand(a.Tasks[task].Command, task, formatted)
				if taskok {
					break
				}
			}
		} else {
			// test failed
			taskok = false
		}

		// Commands, not to be misaken for command
		if taskok {
			cmds := strings.Split(a.Tasks[task].Commands, "\n")
			for _, c := range cmds {
				taskok = a.Tasks[task].RunCommand(c, task, formatted)
				if !taskok {
					break
				}
			}
		}

		// Handle Serve ...
		if taskok && a.Tasks[task].Serve != "" {
			go Serve(".", a.Tasks[task].Serve)
		}

		// Wait ...
		if wait_duration, wait_err := time.ParseDuration(a.Tasks[task].Wait); wait_err == nil {
			<-time.After(wait_duration)
		}

		// The task failed ...
		if !taskok {
			red := color.New(color.FgRed).SprintFunc()
			fmt.Println("\n---")
			fmt.Println(red("✘"), task)

			// Failed? Lets run the failed tasks
			for _, failed := range a.Tasks[task].FailedTasks() {
				if !a.runTask(failed, args, formatted) {
					break
				}
			}
		} else {
			green := color.New(color.FgGreen).SprintFunc()
			fmt.Println("\n---")
			fmt.Println(green("✔"), task)
		}

		// Handle skips ...
		if !taskok && a.Tasks[task].Skip {
			return false
		}

		// Handle exits ...
		if !taskok && a.Tasks[task].Exit != "" {
			if exitCode, err := strconv.Atoi(a.Tasks[task].Exit); err == nil {
				os.Exit(exitCode)
			}
			return false
		}

		var wg sync.WaitGroup
		// Do we have any tasks we need to run in parallel?
		for _, t := range a.Tasks[task].MultiTask() {
			wg.Add(1)
			go func(t string, args []string) {
				defer wg.Done()
				a.runTask(t, args, true)
			}(t, args)
		}
		wg.Wait()

		// Ok, we made it here ... Is this task a task group?
		for _, t := range a.Tasks[task].TaskGroup() {
			if !a.runTask(t, args, formatted) {
				break
			}
		}

		// Woot! Lets run the ok tasks
		if taskok {
			for _, ok_tasks := range a.Tasks[task].OkTasks() {
				if !a.runTask(ok_tasks, args, formatted) {
					break
				}
			}

		}

		// Do we need to break or should we keep going?
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

// Ensure that the task exists
func (a *Alfred) isValidTask(task string) bool {
	if _, exists := a.Tasks[task]; exists {
		return true
	}
	return false
}

// The first argument MUST contain a "/" to be considered remote
func (a *Alfred) isRemote() bool {
	if len(os.Args) >= 2 && strings.Contains(os.Args[1], "/") {
		return true
	}
	return false
}

// Bounce around the web until we find something, or we don't ..
func (a *Alfred) findRemote() bool {
	// Make sure remote is a valid possibility
	if a.isRemote() {
		remote, module := a.remote.Parse(os.Args[1])

		// default to plain jane github
		url := "https://raw.githubusercontent.com/" + os.Args[1] + "/master/alfred.yml"

		// Does a remote exist? If so, we should use the remote syntax
		if a.remote.Exists(remote) {
			url = a.remote.URL(remote, module)
		}

		// try to fetch the alfred file
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != 200 {
			say("error", "Unknown module "+os.Args[1])
			say("url", url)
			return true
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			// We found something ... lets use it!
			//a.contents = append(append(a.contents, []byte("\n")...), body...)
			a.contents = body
			a.location = url
			return true
		}
		return true
	}
	return false
}

// findLocal This function will look locally for an alfred.yml file. First
// starting in the working directory and then going back up the parent
// directory until either an alfred.yml file is found, or you're in the
// root directory
func (a *Alfred) findLocal() bool {
	// Grab the current directory
	dir, err := os.Getwd()
	if err == nil {
		// Just keep going ...
		for {
			// Did we find a bunch of alfred files?
			patterns := []string{
				dir + "/alfred.yml",
				dir + "/.alfred/*alfred.yml",
				dir + "/alfred/*alfred.yml"}
			for _, pattern := range patterns {
				if alfred_files, files_err := filepath.Glob(pattern); files_err == nil && len(alfred_files) > 0 {
					for _, alfred_file := range alfred_files {
						if contents, read_err := ioutil.ReadFile(alfred_file); read_err == nil {
							// Sweet. We found an alfred file. Lets save it off and return
							//a.contents = append(append(a.contents, []byte("\n")...), contents...)
							a.contents = append(a.contents, []byte("\n\n")...)
							a.contents = append(a.contents, contents...)
							a.location = alfred_file
							// Be sure that we ar relative to where we found the config file
							a.dir = dir
						}
					}
					return true
				} else {
					fmt.Println("Error")
				}
			}

			dir = path.Dir(dir)
			if dir == "/" {
				// We've gone too far ...
				break
			}
		}
	}
	// We didn't find anything. /cry
	return false
}

// List out all of the available commands we can run
func (a *Alfred) List() {
	// Get/Sort list of tasks ...
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

// prepare will cycle through all the tasks and prep them to be displayed, or aliased
// If any commands have aliases, lets copy the tasks to their new names
// Also, if we have an astrick in the name, lets promote it
func (a *Alfred) prepare() {
	for name, task := range a.Tasks {
		// Does this task have an alias? If so, lets create it!
		if len(task.Aliases()) > 0 {
			for _, alias := range task.Aliases() {
				a.Tasks[alias] = task
			}
		}

		// Should this task be promoted?
		if strings.HasSuffix(name, "*") {
			a.Tasks[strings.Replace(name, "*", "", -1)] = task
		}
	}
}

// Alfred speaks!
func say(task, msg string) {
	fmt.Println("["+task+"]", msg)
}
