package alfred

import (
	"bufio"
	"fmt"
	"os/exec"
	"sync"
)

// the task component
func commandC(task Task, context *Context, tasks map[string]Task) {
	command(task.Command, task, context, tasks)
}

// within the context of a task, run a command with proper output
// looking for eval, or simple execs? If so, see utils.go
// this one will hook into the GUI where appropriate
func command(commandStr string, task Task, context *Context, tasks map[string]Task) {
	if commandStr == "" {
		return
	}

	cmd := exec.Command("bash", "-c", translate(commandStr, context))

	// set the directory where to run
	cmd.Dir, _ = task.dir(context)

	// wait for output to be completed before moving on
	var wg sync.WaitGroup

	cmdReaderStdOut, _ := cmd.StdoutPipe()
	scannerStdOut := bufio.NewScanner(cmdReaderStdOut)
	go func() {
		wg.Add(1)
		for scannerStdOut.Scan() {
			s := fmt.Sprintf("%s", scannerStdOut.Text())
			output(s, task, context)
		}
		wg.Done()
	}()

	cmdReaderStdErr, _ := cmd.StderrPipe()
	scannerStdErr := bufio.NewScanner(cmdReaderStdErr)
	go func() {
		wg.Add(1)
		for scannerStdErr.Scan() {
			s := fmt.Sprintf("%s", scannerStdErr.Text())
			output(s, task, context)
		}
		wg.Done()
	}()

	err := cmd.Start()
	if err != nil {
		s := fmt.Sprintf("{{ .Text.Failure }}%s{{ .Text.Reset }}", err.Error())
		output(s, task, context)
	}
	statusCode := cmd.Wait()
	wg.Wait()
	fmt.Println("command", commandStr)
	if statusCode != nil {
		task.Exit(context, tasks)
	}
}
