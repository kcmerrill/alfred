package alfred

import (
	"fmt"
	"os"
	"sync"

	"github.com/kcmerrill/common.go/config"
	"github.com/kcmerrill/common.go/file"
	yaml "gopkg.in/yaml.v2"
)

// Alfred coordinates the tasks
type Alfred struct {
	Lock  *sync.Mutex
	Tasks map[string]Task
	Args  []string
}

// Initialize will create and setup a new version of alfred
func Initialize(alfred *Alfred) *Alfred {
	// make sure we don't stomp on anything ....
	if alfred.Tasks == nil {
		alfred.Tasks = make(map[string]Task)
	}

	// basically whatever is left over
	if alfred.Args == nil {
		alfred.Args = os.Args
	}

	// no biggie to stomp on
	alfred.Lock = &sync.Mutex{}

	// return alfred
	return alfred
}

// Task will start off our entrypoint task
func (a *Alfred) Task(cli CLI) {
	var tasks []byte
	if cli.file != "_local" {
		f, err := file.Get(cli.file)
		if err != nil {
			fmt.Println(translate("{{ .Text.FailureIcon }}{{ .Text.Failure }}Unable to load: "+cli.file+"{{ .Text.Reset }}", emptyContext()))
			os.Exit(42)
		}
		tasks = f
	} else {
		// must be local
		local, err := config.Find("alfred.yml")
		if err != nil {
			fmt.Println(translate("{{ .Text.FailureIcon }}{{ .Text.Failure }}Unable to load: "+cli.file+"{{ .Text.Reset }}", emptyContext()))
			os.Exit(42)
		}
		tasks = local
	}

	err := yaml.Unmarshal(tasks, a.Tasks)
	if err != nil {
		fmt.Println(translate("{{ .Text.FailureIcon }}{{ .Text.Failure }}Unable to unmarshal: "+cli.file+"{{ .Text.Reset }}", emptyContext()))
		fmt.Println(translate("{{ .Text.FailureIcon }}{{ .Text.Failure }}"+err.Error()+"{{ .Text.Reset }}", emptyContext()))
		os.Exit(42)
	}

	// lets get rocking and rolling ...
	NewTask(cli.task, InitialContext(cli.args), a.Tasks)
}
