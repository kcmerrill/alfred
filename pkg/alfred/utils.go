package alfred

import (
	"os/exec"
	"strings"
)

func evaluate(command, dir string) string {
	results, ok := execute(command, dir)
	if ok {
		return strings.TrimSpace(results)
	}
	return command
}

func testCommand(command, dir string) bool {
	_, ok := execute(command, dir)
	return ok
}

func execute(command, dir string) (string, bool) {
	cmd := exec.Command("bash", "-c", command)
	cmd.Dir = dir
	cmdOutput, error := cmd.CombinedOutput()
	if error != nil {
		return error.Error(), false
	}
	return string(cmdOutput), true
}

func emptyContext() *Context {
	return InitialContext([]string{})
}

func emptyTaskList() map[string]Task {
	etl := make(map[string]Task)
	return etl
}

func padLeft(word string, size int, padding string) string {
	return pad(word, size, padding) + word
}

func padRight(word string, size int, padding string) string {
	return word + pad(word, size, padding)
}

func pad(word string, size int, padding string) string {
	padded := ""
	if len(word) >= size {
		return padded
	}

	for l := 0; l < size-len(word); l++ {
		padded += padding
	}

	return padded
}
