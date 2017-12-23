package alfred

import (
	"github.com/mgutz/ansi"
)

// Context contains the state of a task
type Context struct {
	TaskName string
	TaskFile string
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
	Task        string
	Warning     string
	Args        string
	Command     string
	Reset       string

	// color codes
	Grey   string
	Orange string
	Green  string
}

// InitialContext will return an empty context
func InitialContext(args []string) *Context {
	return &Context{
		TaskName: "n/a",
		Args:     args,
		Register: make(map[string]string),
		Log:      make([]string, 0),
		Ok:       true, // innocent until proven guilty
		Text: TextConfig{
			// TODO: I don't like this, let me chew on this a bit more
			Success:     ansi.ColorCode("green"),
			SuccessIcon: "✔",
			Failure:     ansi.ColorCode("9"),
			FailureIcon: "✘",
			Task:        ansi.ColorCode("33"),
			Warning:     ansi.ColorCode("185"),
			Command:     ansi.ColorCode("reset"),
			Args:        ansi.ColorCode("198"),
			Reset:       ansi.ColorCode("reset"),

			// Color codes
			Grey:   ansi.ColorCode("238"),
			Orange: ansi.ColorCode("202"),
			Green:  ansi.ColorCode("green"),
		}}
}
