package alfred

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

func configC(task Task, context *Context, tasks map[string]Task) {
	// load up the config, if any
	if task.Config != "" {
		var yamlFile []byte
		dir, _ := task.dir(context)
		if contents, configFileErr := ioutil.ReadFile(translate(task.Config, context)); configFileErr == nil {
			yamlFile = contents
		} else {
			yamlFile = []byte(translate(task.Config, context))
		}
		c := make(map[string]string)
		if configFileUnMarshalErr := yaml.Unmarshal(yamlFile, &c); configFileUnMarshalErr != nil {
			outFail("config", "{{ .Text.Failure }}"+configFileUnMarshalErr.Error(), context)
			// TODO: Should we bail if we can't unmarshal? Or leave it up to the task? ... Let me chew on this
			task.Exit(context, tasks)
		}
		for key, value := range c {
			context.Vars[key] = evaluate(value, dir)
		}
	}
}
