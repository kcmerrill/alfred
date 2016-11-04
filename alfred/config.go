package alfred

import (
	"io/ioutil"
	"os"
	"os/user"

	"gopkg.in/yaml.v2"
)

/*
Config
reads in the user's local yaml configuration file */
func (a *Alfred) Config() {
	a.config.Remote = make(map[string]string)
	usr, err := user.Current()
	if err != nil {
		return
	}

	configFile := usr.HomeDir + "/.alfred/config.yml"

	if _, stat_err := os.Stat(configFile); stat_err == nil {
		if contents, read_err := ioutil.ReadFile(configFile); read_err == nil {
			if err = yaml.Unmarshal([]byte(contents), &a.config); err != nil {
				say("error", "loading config")
				say("config", configFile)
			}
		}
	}
}
