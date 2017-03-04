package alfred

import (
	"io/ioutil"
	"os"
	"os/user"

	"gopkg.in/yaml.v2"
)

// Config reads in the user's configuration file stored in $HOME/.alfred/config.yml
func (a *Alfred) Config() {
	a.config.Remote = make(map[string]string)
	usr, err := user.Current()
	if err != nil {
		return
	}

	configFile := usr.HomeDir + "/.alfred/config.yml"

	if _, statError := os.Stat(configFile); statError == nil {
		if contents, readError := ioutil.ReadFile(configFile); readError == nil {
			if err = yaml.Unmarshal([]byte(contents), &a.config); err != nil {
				say("error", "loading config")
				say("config", configFile)
			}
		}
	}
}
