package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

/*
The test suite is really a functional test.
Once I get the functional suite working, then I'll add unit tests
but first I'd like to get a happy path.

I say all that to say this: This package requires alfred in your path,
and it also requires alfred to already be compiled

-kc
*/

func init() {
	os.Chdir("./examples/demo-everything/")
}

func run(cmd string, t *testing.T) (string, error) {
	if out, err := exec.Command("bash", "-c", cmd).CombinedOutput(); err == nil {
		fmt.Println(string(out))
		return strings.Trim(string(out), "\n"), nil
	} else {
		return strings.Trim(string(out), "\n"), err
	}

}

func TestCurrentDirectory(t *testing.T) {
	if cur, err := os.Getwd(); err == nil {
		if !strings.Contains(cur, "examples/demo-everything") {
			t.Logf("Unable to properly switch to examples/demo-everything")
			t.Fail()
		}
	}
}

func TestListing(t *testing.T) {
	sut, err := run("alfred", t)

	/* Was it able to run? */
	if err != nil {
		t.Logf("Unable to run alfred. Is it compiled and in your path?")
		t.Fail()
	}

	/* Check summary */
	if !strings.Contains(sut, "[one] Displaying the task name") {
		t.Logf("Task one was not in the alfred listing")
		t.Fail()
	}

	/* Check promoted/alphabetic ordering ... */
	if !strings.HasSuffix(sut, `[two] A simple echo

----

[twentysix] Notice the astrick? It means it's a "main" task. Useful for a long alfred file`) {
		t.Logf("Alphabetical ordering is off, and so is promoted")
		t.Fail()
	}

	/* Check Usage */
	if !strings.Contains(sut, "- Usage: alfred fourteen foldername") {
		t.Logf("Usage was not displayed ...")
		t.Fail()
	}

	/* Check Aliases */
	if !strings.Contains(sut, "- Alias: six six.one six.two") {
		t.Logf("Aliases were not displayed ...")
		t.Fail()
	}

	/* Check Private */
	if strings.Contains(sut, "[eight]") {
		t.Logf("Private task [eight] should not be found")
		t.Fail()
	}
}

func TestSimpleCommand(t *testing.T) {
	sut, _ := run("alfred two", t)

	fmt.Println("system", sut)

	/* Check summary */
	if !strings.Contains(sut, "[two] A simple echo") {
		t.Logf("Summary was not displayed")
		t.FailNow()
	}

	/* Check command was run */
	if !strings.Contains(sut, "A simple echo command") {
		t.Logf("Command [two] was not run succesfully")
		t.FailNow()
	}
}

func TestDirKey(t *testing.T) {
	sut, _ := run("alfred three", t)

	if !strings.Contains(sut, "/tmp") {
		t.Logf("The directory should have changed to /tmp")
		t.FailNow()
	}
}

func TestAlias(t *testing.T) {
	sut, _ := run("alfred six", t)

	/* Verify the summary changed */
	if !strings.Contains(sut, "[six] Step five, but aliased as step six too! Space seperated") {
		t.Logf("The task [five] alias should actually be [six] now")
		t.FailNow()
	}

	/* Make sure the ls command ran */
	if !strings.Contains(sut, "emptydir") {
		t.Logf("ls should contain emptydir")
		t.FailNow()
	}

	/* Make sure the ls has emptydir */
	if !strings.Contains(sut, "emptydir") {
		t.Logf("ls should contain the emptydir as well")
		t.FailNow()
	}
}

func TestOkTask(t *testing.T) {
	sut, _ := run("alfred seven", t)

	/* Verify task seven was called */
	if !strings.Contains(sut, "[seven]") {
		t.Logf("Task [seven] should have been called")
		t.FailNow()
	}

	/* Ok, [eight] should have been called too */
	if !strings.Contains(sut, "[eight]") {
		t.Logf("Task [eight] should have been called")
		t.FailNow()
	}

	/* Failure should _not_ have been called */
	if strings.Contains(sut, "[ten]") {
		t.Logf("Task [ten] should not have been called")
		t.FailNow()
	}
}

func TestTaskFail(t *testing.T) {
	sut, err := run("alfred nine", t)

	if err != nil {
		fmt.Println("FAIL---", err)
		t.Logf("While the task failed, alfred should report 0 exit code")
		t.FailNow()
	}

	/* Verify that task [nine] ran */
	if !strings.Contains(sut, "[nine]") {
		t.Logf("Task [nine] should have been called")
		t.FailNow()
	}

	/* Make sure we failed, and we see the output */
	if !strings.Contains(sut, "No such file or directory") {
		t.Logf("ls /kcwashere should not have existed ...")
		t.FailNow()
	}

	/* Verify step [ten] was run as it was the fail task*/
	if !strings.Contains(sut, "[ten]") {
		t.Logf("Because kcwashere doesn't exist, step 10 should have run")
		t.FailNow()
	}

	/* Verify step [eight] was _not_ run, given it was the ok task */
	if strings.Contains(sut, "[eight]") {
		t.Logf("Step [eight] was run, and it should _not_ have")
		t.FailNow()
	}
}

func TestTaskGroup(t *testing.T) {
	sut, _ := run("alfred eleven", t)

	/* task eleven is a task group */
	if !strings.Contains(sut, "[eleven]") {
		t.Logf("Task [eleven] should've been called with a summary!")
		t.FailNow()
	}

	/* four, five and six should've been called */
	if !strings.Contains(sut, "[four]") {
		t.Logf("Task [four] should've been called with a summary!")
		t.FailNow()
	}

	if !strings.Contains(sut, "[five]") {
		t.Logf("Task [four] should've been called with a summary!")
		t.FailNow()
	}

	if !strings.Contains(sut, "[six]") {
		t.Logf("Task [six] should've been called with a summary!")
		t.FailNow()
	}

}

func TestEvery(t *testing.T) {
	cmd := exec.Command("alfred", "thirteen")
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	time.AfterFunc(4*time.Second, func() {
		cmd.Process.Kill()

		buf := new(bytes.Buffer)
		buf.ReadFrom(stdout)
		s := buf.String()

		/* Now, count the number of times we see [thirteen] */
		ran := strings.Split(s, "[thirteen]")

		if len(ran) <= 2 {
			t.Logf("The task [thirteen] did not run as many times as we had hoped")
			t.FailNow()
		}
	})

	<-time.After(5 * time.Second)
}

func TestArgumentsOkAndDefaults(t *testing.T) {
	cmds := []string{"fourteen .", "fifteen"}
	for _, cmd := range cmds {
		sut, _ := run("alfred "+cmd, t)

		/* Verify our task ran */
		if !strings.Contains(sut, "[fourteen]") && !strings.Contains(sut, "[fifteen]") {
			t.Logf("Task should have run succesfully")
			t.FailNow()
		}
		/* Verfify our command had run succesfully */
		if !strings.Contains(sut, "emptydir") {
			t.Logf("Listing of the current directory should show an alfred.yml file")
			t.FailNow()
		}
	}
}

func TestArgumentsFailure(t *testing.T) {
	sut, err := run("alfred fourteen", t)

	if err == nil {
		t.Logf("Missing arguments should cause an invalid exit code")
		t.FailNow()
	}

	/* Verfify our error message */
	if !strings.Contains(sut, "[fourteen:error] Missing argument(s).") {
		t.Logf("The error message is invalid")
		t.FailNow()
	}
}

func TestRetryLogic(t *testing.T) {
	sut, _ := run("alfred twentyseven", t)

	/* Verfify our retry logic via error message */
	if strings.Count(sut, "/step27-idonotexist") != 3 {
		t.Logf("Expected Retry logic 3 times")
		t.FailNow()
	}
}

func TestCommands(t *testing.T) {
	run("alfred twentynine", t)
	if _, err := os.Stat("/tmp/commands.txt"); err == nil {
		t.Logf("/tmp/commands.txt should _NOT_ exist")
		t.FailNow()
	}
}
func TestTest(t *testing.T) {
	run("alfred thirty", t)
	if _, err := os.Stat("/tmp/test.txt"); err == nil {
		t.Logf("/tmp/test.txt should _NOT_ exist")
		t.FailNow()
	}
}
func TestCleanedArgs(t *testing.T) {
	sut, _ := run("alfred thirtyone", t)
	if !strings.Contains(sut, "bingowashisnameo") {
		t.Logf("Expecting bingowashisnameo as a cleanedarg!")
		t.FailNow()
	}
}

func TestTaskMultiArguments(t *testing.T) {
	sut, _ := run("alfred thirty.eight", t)

	if !strings.Contains(sut, "-kc-merrill-") {
		t.Logf("Expecting -kc-merrill- in command output")
		t.FailNow()
	}

	if !strings.Contains(sut, "-bruce-wayne-") {
		t.Logf("Expecting -bruce-wayne- in command output")
		t.FailNow()
	}

	if !strings.Contains(sut, "-clark-kent-") {
		t.Logf("Expecting -clark-kent- in command output")
		t.FailNow()
	}
}

func TestRegisteringVariables(t *testing.T) {
	sut, _ := run("alfred fourty.one", t)

	if !strings.Contains(sut, "The variable is 5678") {
		t.Logf("Expecting 5678 to be returned")
		t.FailNow()
	}
}

func TestExample(t *testing.T) {
	sut, _ := run("alfred", t)
	if len(sut) <= 0 {
		t.Logf("Test example failed ... meaning no data :(")
		t.FailNow()
	}
}
