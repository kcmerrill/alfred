package alfred

import (
	"bufio"
	"fmt"
	"os/exec"
	"sync"
)

func command(task Task, context *Context, tasks map[string]Task) {
	if task.Command == "" {
		return
	}

	cmd := exec.Command("bash", "-c", translate(task.Command, context))

	// set the directory where to run
	cmd.Dir = task.Dir

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
	if statusCode != nil {
		context.Ok = false
		task.Exit()
	}
}
