package alfred

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/user"
)

func (a *Alfred) Config() {
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
