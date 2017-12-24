package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Find will find a file going up the directory tree one at a time stopping when it finds the file
func Find(filename string) (string, []byte, error) {
	dir, err := os.Getwd()

	if err != nil {
		return "./", nil, fmt.Errorf("Unable to get working directory")
	}

	for {
		if _, err := os.Stat(dir + "/" + filename); err == nil {
			if contents, err := ioutil.ReadFile(dir + "/" + filename); err == nil {
				return dir + "/", contents, nil
			}
			return dir + "/", nil, nil
		}

		if dir == "/" {
			break
		}
		dir = filepath.Dir(dir)
	}
	return "./", nil, fmt.Errorf("Unable to find config")
}
