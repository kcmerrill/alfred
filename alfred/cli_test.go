package alfred

import (
	"fmt"
	"testing"
)

func TestCliParseFileAndTaskURL(t *testing.T) {
	c := &CLI{}
	file, task := c.ParseFileAndTask("http://github.com/kcmerrill/alfred/alfred.yml:bingowashisnameo", "alfred")

	if file != "http://github.com/kcmerrill/alfred/alfred.yml" {
		t.Logf("The github url was not parsed properly")
		t.Logf(file)
		t.FailNow()
	}

	if task != "bingowashisnameo" {
		t.Logf("The task was not parsed properly")
		t.FailNow()
	}

	_, task = c.ParseFileAndTask("http://github.com/kcmerrill/alfred/alfred.yml", "alfred")
	if task != "_list" {
		t.Logf("The task was not parsed properly")
		t.FailNow()
	}
}

func TestCliParseFileAndTaskGithub(t *testing.T) {
	c := &CLI{}
	file, task := c.ParseFileAndTask("kcmerrill/alfred:taskname", "alfred")

	if file != "https://raw.githubusercontent.com/kcmerrill/alfred/master/alfred.yml" {
		t.Logf("The github project url was not parsed properly")
		t.Logf(file)
		t.FailNow()
	}

	if task != "taskname" {
		t.Logf("The task was not parsed properly")
		t.FailNow()
	}

	// lets try again, but with no task name
	_, task = c.ParseFileAndTask("kcmerrill/alfred", "alfred")
	if task != "_list" {
		t.Logf("The task was not parsed properly")
		t.FailNow()
	}
}

func TestCliParseFileAndTaskLocal(t *testing.T) {
	c := &CLI{}
	file, task := c.ParseFileAndTask("sometask", "alfred")

	if file != "_local" {
		t.Logf("This is a local task, file = _local")
		t.Logf(file)
		t.FailNow()
	}

	if task != "sometask" {
		t.Logf("The task was not parsed properly")
		t.FailNow()
	}
}

func TestCLIParse(t *testing.T) {
	c := NewCLI([]string{})
	c.Parse([]string{"alfred", "http://github.com/kcmerrill/alfred/alfred.yml:my.cool.task", "arg1", "arg2"})

	if c.name != "alfred" {
		t.Logf("The name of the application should be alfred")
		t.FailNow()
	}

	if c.file != "http://github.com/kcmerrill/alfred/alfred.yml" {
		t.Logf("The name of the file should be github.com(ish)")
		t.FailNow()
	}

	if c.task != "my.cool.task" {
		t.Logf("The name of the task parsed should be 'my.cool.task")
		t.FailNow()
	}

	if len(c.args) != 2 {
		t.Logf("The args should be arg1 and arg2")
		t.FailNow()
	}

	c = NewCLI([]string{})
	c.Parse([]string{"alfred"})
	if c.file != "_local" {
		t.Logf("Simply calling the application should yield a file of _local")
		t.FailNow()
	}

	if c.name != "alfred" {
		t.Logf("The name of the application should be 'alfred'")
		t.FailNow()
	}

	if c.task != "_list" {
		t.Logf("The default task should be _list")
		t.FailNow()
	}

	if len(c.args) != 0 {
		t.Logf("No arguments should be present")
		fmt.Println(c.args)
		t.FailNow()
	}
}
