package alfred

import "os/exec"

func evaluate(command, dir string) string {
	results, ok := execute(command, dir)
	if ok {
		return results
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
	if error == nil {
		return error.Error(), false
	}
	return string(cmdOutput), true
}
