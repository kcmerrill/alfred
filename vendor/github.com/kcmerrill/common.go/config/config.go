package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// FindAndCombine will go up directories looking for a file + extension and combine all the files into one []byte{}
func FindAndCombine(file, extension string) (string, []byte, error) {
	// Grab the current directory
	dir, err := os.Getwd()
	combinedContents := []byte{}
	if err == nil {
		// Just keep going ...
		for {
			// Did we find a bunch of config files?
			patterns := []string{
				dir + "/" + file + "." + extension,
				dir + "/." + file + "/*" + file + "." + extension,
				dir + "/" + file + "/*" + file + "." + extension}
			for _, pattern := range patterns {
				if configFiles, filesErr := filepath.Glob(pattern); filesErr == nil && len(configFiles) > 0 {
					for _, configFile := range configFiles {
						if contents, readErr := ioutil.ReadFile(configFile); readErr == nil {
							// Sweet. We found an config file. Lets save it off and return
							combinedContents = append(combinedContents, []byte("\n\n")...)
							combinedContents = append(combinedContents, contents...)
						}
					}
					return dir, combinedContents, nil
				}
			}

			dir = path.Dir(dir)
			if dir == "/" {
				// We've gone too far ...
				break
			}
		}
	}
	// We didn't find anything. /cry
	return "./", []byte{}, fmt.Errorf("Unable to look up directory")
}

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
