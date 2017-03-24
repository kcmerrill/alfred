package alfred

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Config reads in the user's configuration file stored in $HOME/.alfred/config.yml
func (a *Alfred) Config() {
	a.config.Remote = make(map[string]string)

	// os.User is really terrible. Also doesn't work with cross compiling
	configFile := os.Getenv("HOME") + "/.alfred/config.yml"

	if _, statError := os.Stat(configFile); statError == nil {
		if contents, readError := ioutil.ReadFile(configFile); readError == nil {
			if err := yaml.Unmarshal([]byte(contents), &a.config); err != nil {
				say("error", "loading config")
				say("config", configFile)
			}
		}
	}
}
