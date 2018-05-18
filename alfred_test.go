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
		"invalid.task":                 TestComponent{fail: true},
		"invalid.task:formatted":       TestComponent{fail: true, expectedOutput: []string{"0s"}},
		"invalid.task:not.formatted":   TestComponent{args: "--no-formatting", fail: true, noOutput: []string{"0s"}},
		"summary":                      TestComponent{expectedOutput: []string{"TESTING SUMMARY"}},
		"command":                      TestComponent{expectedOutput: []string{"HELLO ALFRED", "HELLO NEWLINE"}},
		"commands":                     TestComponent{fail: true, expectedOutput: []string{"HELLO ALFRED"}, noOutput: []string{"THIS LINE NOT SHOWN"}},
		"exit":                         TestComponent{fail: true},
		"arguments:without":            TestComponent{fail: true},
		"arguments":                    TestComponent{params: "ARG1", expectedOutput: []string{"ARG::ARG1"}},
		"default.arguments":            TestComponent{expectedOutput: []string{"ARG::ARG1"}},
		"required.arguments":           TestComponent{fail: true, noOutput: []string{"ARG::ARG2"}},
		"required.arguments:with_args": TestComponent{params: "ARG1", expectedOutput: []string{"ARG::ARG1 ARG::ARG2"}},
		"ok":                TestComponent{expectedOutput: []string{"TEST::OK", "ARG::ARG1", "ARG::OKTEST"}},
		"ok:witharg":        TestComponent{params: "DEFAULTARG", expectedOutput: []string{"TEST::OK", "ARG::DEFAULTARG", "ARG::OKTEST"}},
		"fail":              TestComponent{expectedOutput: []string{"TEST::FAIL", "ARG::ARG1", "ARG::FAILTEST"}},
		"tasks":             TestComponent{expectedOutput: []string{"ARG::ARG1", "ARG::TASKTEST"}},
		"multitask":         TestComponent{expectedOutput: []string{"ARG::ARG1", "ARG::MULTITASKTEST"}},
		"check:should_skip": TestComponent{params: "alfred.yml", expectedOutput: []string{"skipped"}, noOutput: []string{"TEST::CHECK"}},
		"check:should_run":  TestComponent{params: "doesnotexist", expectedOutput: []string{"TEST::CHECK"}},
		"config:file":       TestComponent{params: "config.yml", expectedOutput: []string{"FOO::BAR", "FIZZ::BUZZ"}},
		"config:text":       TestComponent{params: "'foo: BAR\nfizz: BUZZ'", expectedOutput: []string{"FOO::BAR", "FIZZ::BUZZ"}},
		"dir":               TestComponent{params: "/tmp/", expectedOutput: []string{"PWD::/private/tmp|PWD::/tmp"}},

		// Tests that _should_ be working(read: bugs)
		"arguments:empty_log":   TestComponent{skip: true, fail: true, args: "--log scratch/test_empty.log", filesExist: []string{"scratch/"}},
		"arguments:log_written": TestComponent{skip: true, args: "--log scratch/test.log", filesExist: []string{"scratch/test.log"}},

		// Tests to write
		"dir:with_alfred_folder":        TestComponent{skip: true, fail: true},
		"dir:without_alfred_folder":     TestComponent{skip: true, fail: true},
		"dir:catalog":                   TestComponent{skip: true, fail: true},
		"dir:catalog_on_same_dir_level": TestComponent{skip: true, fail: true},
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
		taskBits := strings.Split(task, ":")
		taskName := taskBits[0]
		output, ok := runAlfredCommand(taskName, tc)

		if !tc.fail != ok {
			t.Fatalf("[" + task + "]TaskResult: Expected task to fail ...")
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

		for _, ignoreOutput := range tc.noOutput {
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
	noOutput       []string
	params         string
	args           string
	fail           bool
	skip           bool
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

	fmt.Println("Command:", cmd)
	return testCommand(cmd, ".")
}
