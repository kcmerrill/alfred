package alfred

import (
	"bufio"
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

	cmdFailed := false

	// wait for output to be completed before moving on
	var wg sync.WaitGroup
	cmdReaderStdOut, _ := cmd.StdoutPipe()
	scannerStdOut := bufio.NewScanner(cmdReaderStdOut)
	go func() {
		wg.Add(1)
		for scannerStdOut.Scan() {
			cmdOK(scannerStdOut.Text(), context)
			cmdFailed = false
		}
		wg.Done()
	}()

	cmdReaderStdErr, _ := cmd.StderrPipe()
	scannerStdErr := bufio.NewScanner(cmdReaderStdErr)
	go func() {
		wg.Add(1)
		for scannerStdErr.Scan() {
			cmdFailed = true
			cmdFail(scannerStdErr.Text(), context)
		}
		wg.Done()
	}()

	err := cmd.Start()
	if err != nil {
		cmdFail(scannerStdErr.Text(), context)
	}
	statusCode := cmd.Wait()
	wg.Wait()
	if statusCode != nil {
		if !cmdFailed {
			// was the last thing we saw a not failure?
			outFail("command", "failed", context)
		}
		task.Exit(context, tasks)
	}
}
