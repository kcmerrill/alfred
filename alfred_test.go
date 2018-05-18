package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

func TestAlfredComponents(t *testing.T) {
	tt := map[string]TestComponent{
		"|show.tasks":                  TestComponent{expectedOutput: []string{" summary            | TESTING SUMMARY"}, failWithOutput: []string{"hidden.task"}},
		"invalid.task":                 TestComponent{shouldFail: true},
		"invalid.task|formatted":       TestComponent{shouldFail: true, expectedOutput: []string{"0s"}},
		"invalid.task|not.formatted":   TestComponent{args: "--no-formatting", shouldFail: true, failWithOutput: []string{"0s"}},
		"summary":                      TestComponent{expectedOutput: []string{"TESTING SUMMARY"}},
		"hidden.task":                  TestComponent{expectedOutput: []string{"Testing a hidden task"}},
		"command":                      TestComponent{expectedOutput: []string{"HELLO ALFRED", "HELLO NEWLINE"}},
		"commands":                     TestComponent{shouldFail: true, expectedOutput: []string{"HELLO ALFRED"}, failWithOutput: []string{"THIS LINE NOT SHOWN"}},
		"exit":                         TestComponent{shouldFail: true},
		"arguments|without":            TestComponent{shouldFail: true},
		"arguments":                    TestComponent{params: "ARG1", expectedOutput: []string{"ARG::ARG1"}},
		"default.arguments":            TestComponent{expectedOutput: []string{"ARG::ARG1"}},
		"required.arguments":           TestComponent{shouldFail: true, failWithOutput: []string{"ARG::ARG2"}},
		"required.arguments|with_args": TestComponent{params: "ARG1", expectedOutput: []string{"ARG::ARG1 ARG::ARG2"}},
		"ok":                TestComponent{expectedOutput: []string{"TEST::OK", "ARG::ARG1", "ARG::OKTEST"}},
		"ok|witharg":        TestComponent{params: "DEFAULTARG", expectedOutput: []string{"TEST::OK", "ARG::DEFAULTARG", "ARG::OKTEST"}},
		"fail":              TestComponent{expectedOutput: []string{"TEST::FAIL", "ARG::ARG1", "ARG::FAILTEST"}},
		"tasks":             TestComponent{expectedOutput: []string{"ARG::ARG1", "ARG::TASKTEST"}},
		"multitask":         TestComponent{expectedOutput: []string{"ARG::ARG1", "ARG::MULTITASKTEST"}},
		"check|should_skip": TestComponent{params: "alfred.yml", expectedOutput: []string{"skipped"}, failWithOutput: []string{"TEST::CHECK"}},
		"check|should_run":  TestComponent{params: "doesnotexist", expectedOutput: []string{"TEST::CHECK"}},
		"config|file":       TestComponent{params: "config.yml", expectedOutput: []string{"FOO::BAR", "FIZZ::BUZZ"}},
		"config|text":       TestComponent{params: "'foo: BAR\nfizz: BUZZ'", expectedOutput: []string{"FOO::BAR", "FIZZ::BUZZ"}},
		"dir":               TestComponent{params: "/tmp/", expectedOutput: []string{"PWD::/private/tmp|PWD::/tmp"}},
		"dir|with_alfred_folder": TestComponent{dir: "catalog2", failWithOutput: []string{"tests/catalog2/alfred"}},
		"wait":           TestComponent{expectedOutput: []string{"wait wait 2s"}},
		"template|sprig": TestComponent{expectedOutput: []string{"HELLO!HELLO!HELLO!HELLO!HELLO!"}},
		"for":            TestComponent{expectedOutput: []string{"ARG::0", "ARG::1", "ARG::2", "ARG::3", "ARG::4"}, failWithOutput: []string{"ARG::5"}},
		"include":        TestComponent{expectedOutput: []string{"taska.alfred.yml", "taskb.alfred.yml"}},
		"@catalog":       TestComponent{expectedOutput: []string{"taska"}, failWithOutput: []string{"taskb"}},
		"@catalog:taska": TestComponent{expectedOutput: []string{"taska.alfred.yml"}},
		"env":            TestComponent{expectedOutput: []string{"ARG::TESTENV"}},
		"register":       TestComponent{expectedOutput: []string{"REGISTER::var", "WHOAMI::(kcmerrill|root)"}},

		// Tests that _should_ be working(read: bugs)
		"arguments|empty_log":   TestComponent{skip: true, shouldFail: true, args: "--log scratch/test_empty.log", filesExist: []string{"scratch/"}},
		"arguments|log_written": TestComponent{skip: true, args: "--log scratch/test.log", filesExist: []string{"scratch/test.log"}},
	}

	/* SHOULD NOT HAVE TO CHANGE TOO MUCH BELOW. Add tests above ^^^^ */

	// lets be in the correct directory
	err := os.Chdir("tests/")
	if err != nil {
		t.Fatalf("Could not switch directories. Bailing ...")
	}

	// lets clean up our scratch directory
	testCommand("rm -rf scratch/", ".")

	// lets cycle through our test map and get to testing!
	for task, tc := range tt {
		if tc.skip {
			// we should skip this test
			continue
		}
		taskBits := strings.Split(task, "|")
		taskName := taskBits[0]
		output, ok := runAlfredCommand(taskName, tc)

		if tc.shouldFail && tc.shouldFail != !ok {
			fmt.Println(output)
			t.Fatalf("[" + task + "] Expected task failure ...")
		}

		if !tc.shouldFail && !ok {
			fmt.Println(output)
			t.Fatalf("[" + task + "] Expected task success ...")
		}

		for _, expectedOutput := range tc.expectedOutput {
			found := false
			for _, actualOutput := range strings.Split(output, "\n") {

				if matched, _ := regexp.MatchString(expectedOutput, actualOutput); matched {
					found = true
					break
				}
			}

			if !found {
				t.Fatalf("["+task+"]Output: Expected to find '%s' in: \n\n%s", expectedOutput, output)
			}
		}

		for _, ignoreOutput := range tc.failWithOutput {
			found := false
			for _, actualOutput := range strings.Split(output, "\n") {
				if strings.Contains(actualOutput, ignoreOutput) {
					found = true
					break
				}
			}

			if found {
				t.Fatalf("["+task+"]Output: Expected to _NOT_ find '%s' in: \n\n%s", ignoreOutput, output)
			}
		}

		for _, file := range tc.filesExist {
			if _, exists := os.Stat(file); exists != nil {
				t.Fatalf("["+task+"]File: Expected the file '%s' to exist ...", file)
			}
		}
	}
}

type TestComponent struct {
	filesExist     []string
	fileContents   map[string]string
	tasksStarted   []string
	expectedOutput []string
	failWithOutput []string
	params         string
	args           string
	shouldFail     bool
	skip           bool
	dir            string
}

func testCommand(command, dir string) (string, bool) {
	cmd := exec.Command("bash", "-c", command)
	cmd.Dir = dir
	cmdOutput, error := cmd.CombinedOutput()
	if error != nil {
		return string(cmdOutput), false
	}
	return string(cmdOutput), true
}

func runAlfredCommand(task string, tc TestComponent) (string, bool) {
	cmd := "alfred --no-colors"
	if tc.args != "" {
		cmd += " " + tc.args
	}

	cmd += " " + task

	if tc.params != "" {
		cmd += " " + tc.params
	}

	dir := "."
	if tc.dir != "" {
		dir = tc.dir
	}
	return testCommand(cmd, dir)
}
