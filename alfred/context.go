package alfred

import (
	"os"

	"github.com/mgutz/ansi"
)

// Context contains the state of a task
type Context struct {
	TaskName string
	Log      []string
	Args     []string
	Register map[string]string
	Ok       bool
	Text     TextConfig
	Silent   bool
}

// TextConfig contains configuration needed to display text
type TextConfig struct {
	Success     string
	SuccessIcon string
	Failure     string
	FailureIcon string
	Reset       string
}

// InitialContext will return an empty context
func InitialContext(args []string) *Context {
	return &Context{
		TaskName: "n/a",
		Args:     os.Args[1:],
		Register: make(map[string]string),
		Log:      make([]string, 0),
		Ok:       true, // innocent until proven guilty
		Text: TextConfig{
			Success:     ansi.ColorCode("green"),
			SuccessIcon: "✔",
			Failure:     ansi.ColorCode("red"),
			FailureIcon: "✘",
			Reset:       ansi.ColorCode("reset"),
		}}
}
