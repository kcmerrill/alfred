package alfred

import (
	"bufio"
	"os"
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

	for retry := 0; retry <= task.Retry; retry++ {
		cmd := exec.Command("bash", "-c", translate(commandStr, context))
		cmd.Stdin = os.Stdin

		// set the directory where to run
		cmd.Dir, _ = task.dir(context)

		// wait for output to be completed before moving on
		var wg sync.WaitGroup
		cmdReaderStdOut, _ := cmd.StdoutPipe()
		scannerStdOut := bufio.NewScanner(cmdReaderStdOut)
		go func() {
			wg.Add(1)
			defer wg.Done()
			for scannerStdOut.Scan() {
				cmdOK(scannerStdOut.Text(), context)
			}
		}()

		cmdReaderStdErr, _ := cmd.StderrPipe()
		scannerStdErr := bufio.NewScanner(cmdReaderStdErr)
		go func() {
			wg.Add(1)
			defer wg.Done()
			for scannerStdErr.Scan() {
				cmdFail(scannerStdErr.Text(), context)
			}
		}()

		err := cmd.Start()
		if err != nil {
			cmdFail(scannerStdErr.Text(), context)
		}
		statusCode := cmd.Wait()
		wg.Wait()
		if statusCode != nil {
			task.Exit(context, tasks)
		} else {
			return
		}
	}
	// was the last thing we saw a not failure?
	outFail("command failed", "", context)
}
