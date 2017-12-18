package file

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Get will return a url in []bytes
func Get(url string) ([]byte, error) {
	resp, err := http.Get(url)

	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf("Unable to find file: " + url)
	}

	defer resp.Body.Close()
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("Unable to read the body of the file: " + url)
	}
	return body, nil
}
