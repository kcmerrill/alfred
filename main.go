package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mgutz/ansi"

	. "github.com/kcmerrill/alfred/alfred"
)

var (
	Version = "Development"
	Commit  = ""
)

func main() {
	version := flag.Bool("version", false, "Alfred's version number")
	disableColors := flag.Bool("no-colors", false, "Disable colors")
	disableFormatting := flag.Bool("no-formatting", false, "Show only raw command output")
	debug := flag.Bool("debug", false, "Only show commands to be run")
	log := flag.String("log", "", "Log all tasks to <file>")
	flag.Parse()

	/* Giddy up! */
	if *version {
		fmt.Println()
		fmt.Println("Alfred - Even Batman needs a little help.")
		if Version != "Development" {
			fmt.Print("v", Version)
			fmt.Println("@" + Commit[0:9])
		} else {
			fmt.Println(Version)
		}
		fmt.Println()
		fmt.Println("---")
		fmt.Println("Made with " + ansi.ColorCode("red") + "<3" + ansi.ColorCode("reset") + " by kcmerrill")
		fmt.Println()
		return
	}

	tasks := make(map[string]Task)
	task, args := CLI(flag.Args())
	context := InitialContext(args)

	if *disableColors {
		context.Text = TextConfig{}
	}

	if *log != "" {
		Log(*log, context)
	}

	if *disableFormatting {
		context.Text.DisableFormatting = *disableFormatting
	}

	if *debug {
		context.Text = TextConfig{}
		context.Text.DisableFormatting = true
		context.Debug = true
	}

	// anything from stdin?
	stdinFileInfo, _ := os.Stdin.Stat()
	if stdinFileInfo.Mode()&os.ModeNamedPipe != 0 {
		stdinContent, _ := ioutil.ReadAll(os.Stdin)
		context.Stdin = strings.TrimSpace(string(stdinContent))
	}

	// don't do this if they are :listing
	if len(os.Args) >= 2 {
		NewTask("__init", context, tasks)
	}

	// start the task
	NewTask(task, context, tasks)

	// don't do this if they are :listing
	if len(os.Args) >= 2 {
		NewTask("__exit", context, tasks)
	}
}
