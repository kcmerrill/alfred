package alfred

import "os/exec"

func evaluate(command string) string {
	results, ok := execute(command)
	if ok {
		return results
	}
	return command
}

func testCommand(command string) bool {
	_, ok := execute(command)
	return ok
}

func execute(command string) (string, bool) {
	cmd, error := exec.Command("bash", "-c", command).CombinedOutput()
	if error == nil {
		return error.Error(), false
	}
	return string(cmd), true
}
