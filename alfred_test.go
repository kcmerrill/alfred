package main

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

/*
Lets add some integration tests.

We will build alfred, run some simple tests end to end
*/

func execAlfred(taskname string) (string, bool) {
	cmd := exec.Command("bash", "-c", "alfred --no-colors "+taskname)
	cmd.Dir = "./demo"
	cmdOutput, error := cmd.CombinedOutput()
	if error != nil {
		return error.Error(), false
	}
	return string(cmdOutput), true
}

func TestSummary(t *testing.T) {
	results, _ := execAlfred("summary")
	if !strings.Contains(results, "Testing out the summary") {
		t.Fatalf("Summary should be printed out")
	}
}

func TestArguments(t *testing.T) {
	results, ok := execAlfred("arguments hello world")
	if !ok {
		fmt.Println(results)
		t.Fatalf("Argument task w/args should be ok")
	}

	if !strings.Contains(results, "args:hello world") {
		fmt.Println(results)
		t.Fatalf("Both arguments should be displayed")
	}

	// now lets test it without passing arguments
	results, ok = execAlfred("arguments hello")
	if ok {
		fmt.Println(results)
		t.Fatalf("1/2 argument should cause a failure")
	}
}

func TestDefaults(t *testing.T) {
	results, ok := execAlfred("defaults hello")
	if !ok {
		fmt.Println(results)
		t.Fatalf("Default argument not working")
	}

	if !strings.Contains(results, "defaults:hello world") {
		fmt.Println(results)
		t.Fatalf("Argument world should have been set by default")
	}

	results, _ = execAlfred("defaults hello kc")
	if !strings.Contains(results, "defaults:hello kc") {
		fmt.Println(results)
		t.Fatalf("Argument 2 was passedin and not set properly")
	}
}

func TestStdin(t *testing.T) {
	results, ok := execAlfred("stdin")
	if !ok {
		fmt.Println(results)
		t.Fatalf("No errors should have been reported")
	}

	if !strings.Contains(results, "826662fa3f864c2fccc1d53d85db60c7") {
		fmt.Println(results)
		t.Fatalf("Expected '826662fa3f864c2fccc1d53d85db60c7' as an md5 hash")
	}
}

func TestDir(t *testing.T) {
	results, ok := execAlfred("ls ../pkg/alfred")
	if !ok {
		fmt.Println(results)
		t.Fatalf("Expected a directory listing of ../alfred")
	}

	if !strings.Contains(results, "defaults_test.go") {
		fmt.Println(results)
		t.Fatalf("Expected to see 'defaults_test.go' in ../alfred directory")
	}
}

func TestConfig(t *testing.T) {
	results, ok := execAlfred("config.text")
	if !ok {
		fmt.Println(results)
		t.Fatalf("config.text should not have resulted in an error")
	}

	if !strings.Contains(results, "kc:val-merrill") {
		fmt.Println(results)
		t.Fatalf("config.text kc: merrill not found")
	}

	if !strings.Contains(results, "hello:val-world") {
		fmt.Println(results)
		t.Fatalf("config.text hello: world not found")
	}
}
