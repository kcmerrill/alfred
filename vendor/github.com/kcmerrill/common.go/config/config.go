package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// FindAndCombine will go up directories looking for a file + extension and combine all the files into one []byte{}
func FindAndCombine(currentDir, query, extension string) (string, []byte, error) {
	// Grab the current directory
	combinedContents := []byte{}
	// Just keep going ...
	for {
		// Did we find a bunch of config files?

		dir := ""
		file := query
		if strings.Contains(query, string(os.PathSeparator)) {
			path := strings.SplitN(query, string(os.PathSeparator), 2)
			dir = string(os.PathSeparator) + path[0] + string(os.PathSeparator)
			file = path[1]
		}

		patterns := []string{
			currentDir + "/" + file + "." + extension,
			currentDir + "/." + file + "/*" + file + "." + extension,
			currentDir + "/*" + file + "." + extension,
			currentDir + "/" + file + "/*" + file + "." + extension}

		if dir != "" {
			patterns = []string{
				currentDir + dir + file + "." + extension,
				currentDir + dir + "." + file + "/*" + file + "." + extension,
				currentDir + dir + "*" + file + "." + extension,
				currentDir + dir + file + "/*" + file + "." + extension}
		}

		for _, pattern := range patterns {
			if configFiles, filesErr := filepath.Glob(pattern); filesErr == nil && len(configFiles) > 0 {
				for _, configFile := range configFiles {
					if contents, readErr := ioutil.ReadFile(configFile); readErr == nil {
						// Sweet. We found an config file. Lets save it off and return
						combinedContents = append(combinedContents, []byte("\n\n")...)
						combinedContents = append(combinedContents, contents...)
					}
				}
				return currentDir, combinedContents, nil
			}
		}

		currentDir = path.Dir(currentDir)
		if currentDir == "/" {
			// We've gone too far ...
			break
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
