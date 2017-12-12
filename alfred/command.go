package alfred

import (
	"bufio"
	"fmt"
	"os/exec"

	event "github.com/kcmerrill/hook"
)

func command(task Task, context *Context, tasks map[string]Task) {
	if task.Command == "" {
		return
	}

	cmd := exec.Command("bash", "-c", task.Template(task.Command, context))

	// set the directory where to run
	cmd.Dir = task.Dir

	cmdReaderStdOut, _ := cmd.StdoutPipe()
	scannerStdOut := bufio.NewScanner(cmdReaderStdOut)
	go func() {
		for scannerStdOut.Scan() {
			s := fmt.Sprintf("%s", scannerStdOut.Text())
			event.Trigger("speak", s, task, context)
		}
	}()

	cmdReaderStdErr, _ := cmd.StderrPipe()
	scannerStdErr := bufio.NewScanner(cmdReaderStdErr)
	go func() {
		for scannerStdErr.Scan() {
			s := fmt.Sprintf("%s", scannerStdErr.Text())
			event.Trigger("speak", s, task, context)
		}
	}()

	err := cmd.Start()
	if err != nil {
		s := fmt.Sprintf("{{ .Text.Failure }}%s{{ .Text.Reset }}", err.Error())
		event.Trigger("speak", s, task, context)
	}
	statusCode := cmd.Wait()
	if statusCode != nil {
		context.Ok = false
		task.Exit()
	}
}
