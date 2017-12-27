package config

import (
	"os"
	"strings"
	"testing"
)

func TestFindAndCombine(t *testing.T) {
	c, _ := os.Getwd()
	if os.Chdir("t/alfred/level.one/level.two") != nil {
		t.Fatalf("Test directories do not exist")
	}

	dir, contents, err := FindAndCombine("alfred", "yml")

	if err != nil {
		t.Errorf("Huh. Wasn't expecting an error.")
	}

	if !strings.Contains(string(contents), "a") {
		t.Errorf("Expected a to be in the file")
	}

	if !strings.Contains(string(contents), "b") {
		t.Errorf("Expected b to be in the file")
	}

	if !strings.Contains(string(contents), "c") {
		t.Errorf("Expected c to be in the file")
	}

	if !strings.HasSuffix(dir, "common.go/config/t") {
		t.Errorf("Expecting the directory to contain common.go/config/t/ instead got " + dir)
	}

	os.Chdir(c)
}

func TestFindConfig(t *testing.T) {
	c, _ := os.Getwd()
	if os.Chdir("t/a/b/c") != nil {
		t.Errorf("Unable to switch to directory 'c'")
	}
	if dir, contents, err := Find("a.yml"); err == nil {
		if string(contents) != "a\n" {
			t.Errorf("Contents of a.yml != 'a'")
		}
		if !strings.HasSuffix(dir, "common.go/config/t/a/") {
			t.Errorf("Expecting the directory to contain common.go/config/t/a/ instead got " + dir)
		}
	} else {
		t.Errorf("Unable to find a.yml")
	}
	os.Chdir(c)
}
