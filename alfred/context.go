package alfred

import (
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/mgutz/ansi"
)

// Context contains the state of a task
type Context struct {
	TaskName    string
	TaskFile    string
	Stdin       string
	Started     time.Time
	Log         map[string]*os.File
	Args        []string
	AllArgs     string
	Register    map[string]string
	Ok          bool
	Text        TextConfig
	Silent      bool
	Status      string
	Component   string
	Vars        map[string]string
	Lock        *sync.Mutex
	Out         io.Writer
	DryRun      bool
	Interactive bool
}

// TextConfig contains configuration needed to display text
type TextConfig struct {
	DisableFormatting    bool
	Success              string
	SuccessIcon          string
	Failure              string
	FailureIcon          string
	Task                 string
	Warning              string
	Args                 string
	Command              string
	Reset                string
	TerminalNewLineReset string
	TerminalNewLine      string

	// color codes
	Grey   string
	Orange string
	Green  string
}

// SetVar will the vars map with a given value
func (c *Context) SetVar(key, value string) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	c.Vars[key] = value
}

// GetVar will return a value with a given key, else returns default
func (c *Context) GetVar(key, defaults string) string {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	if v, exists := c.Vars[key]; exists {
		return v
	}
	return defaults
}

// InitialContext will return an empty context
func InitialContext(args []string) *Context {
	return &Context{
		TaskName: "n/a",
		Args:     args,
		AllArgs:  strings.Join(args, " "),
		Register: make(map[string]string),
		Log:      make(map[string]*os.File, 0),
		Ok:       true, // innocent until proven guilty
		Started:  time.Now(),
		Status:   "",
		Vars:     make(map[string]string, 0),
		Lock:     &sync.Mutex{},
		Out:      os.Stdout,
		Text: TextConfig{
			// TODO: I don't like this, let me chew on this a bit more
			Success:              ansi.ColorCode("green"),
			SuccessIcon:          "✔",
			Failure:              ansi.ColorCode("9"),
			FailureIcon:          "✘",
			Task:                 ansi.ColorCode("33"),
			Warning:              ansi.ColorCode("185"),
			Command:              ansi.ColorCode("reset"),
			Args:                 ansi.ColorCode("162"),
			Reset:                ansi.ColorCode("reset"),
			TerminalNewLineReset: "\033[1000D\033[K",
			TerminalNewLine:      "\033[1000D",

			// Color codes
			Grey:   ansi.ColorCode("238"),
			Orange: ansi.ColorCode("202"),
			Green:  ansi.ColorCode("green"),
		}}
}

func copyContex(context *Context, args []string) Context {
	context.Lock.Lock()
	defer context.Lock.Unlock()

	// silly maps, pointers are for kids
	c := *context

	regs := make(map[string]string)
	for k, v := range c.Register {
		regs[k] = v
	}
	c.Register = regs

	a := make([]string, 0)
	for x := 0; x < len(c.Args); x++ {
		a = append(a, c.Args[x])
	}
	c.Args = a

	for idx, v := range args {
		if idx < len(c.Args) {
			c.Args[idx] = v
		} else {
			c.Args = append(c.Args, v)
		}
	}

	logs := make(map[string]*os.File)
	for k, v := range c.Log {
		logs[k] = v
	}
	c.Log = logs

	return c
}
