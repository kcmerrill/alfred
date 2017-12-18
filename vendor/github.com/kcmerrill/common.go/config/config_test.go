package config

import (
	"os"
	"testing"
)

func TestFindConfig(t *testing.T) {
	c, _ := os.Getwd()
	if os.Chdir("t/a/b/c") != nil {
		t.Errorf("Unable to switch to directory 'c'")
	}
	if contents, err := Find("a.yml"); err == nil {
		if string(contents) != "a\n" {
			t.Errorf("Contents of a.yml != 'a'")
		}
	} else {
		t.Errorf("Unable to find a.yml")
	}
	os.Chdir(c)
}
