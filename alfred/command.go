package alfred

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"sync"
)

// the task component
func commandC(task Task, context *Context, tasks map[string]Task) {
	if len(context.Log) != 0 {
		command(task.Command, task, context, tasks)
		return
	}
	command(task.Command, task, context, tasks)
	//commandBasic(task.Command, task, context, tasks)
}

func commandBasic(commandStr string, task Task, context *Context, tasks map[string]Task) {
	if commandStr == "" {
		return
	}

	for retry := 0; retry <= task.Retry; retry++ {
		cmd := exec.Command("bash", "-c", translate(commandStr, context))
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
}

// within the context of a task, run a command with proper output
// looking for eval, or simple execs? If so, see utils.go
// this one will hook into the GUI where appropriate
func command(commandStr string, task Task, context *Context, tasks map[string]Task) {
	if commandStr == "" {
		return
	}

	// hack
	cmdOK("", context)

	for retry := 0; retry <= task.Retry; retry++ {
		cmd := exec.Command("bash", "-c", translate(commandStr, context))
		if context.Stdin != "" {
			cmd.Stdin = bytes.NewBufferString(context.Stdin)
		} else {
			cmd.Stdin = os.Stdin
		}

		// set the directory where to run
		cmd.Dir, _ = task.dir(context)

		// wait for output to be completed before moving on
		var wg sync.WaitGroup
		cmdReaderStdOut, _ := cmd.StdoutPipe()
		scannerStdOut := bufio.NewScanner(cmdReaderStdOut)
		scannerStdOut.Split(bufio.ScanBytes)
		go func() {
			wg.Add(1)
			defer wg.Done()
			for scannerStdOut.Scan() {
				cmdOK(scannerStdOut.Text(), context)
			}
		}()

		cmdReaderStdErr, _ := cmd.StderrPipe()
		scannerStdErr := bufio.NewScanner(cmdReaderStdErr)
		scannerStdErr.Split(bufio.ScanBytes)
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
		wg.Wait()
		statusCode := cmd.Wait()

		if statusCode != nil {
			task.Exit(context, tasks)
		} else {
			return
		}
	}
	// was the last thing we saw a not failure?
	outFail("command failed", "", context)
}
