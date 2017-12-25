package main

import (
	"flag"
	"fmt"

	"github.com/mgutz/ansi"

	. "github.com/kcmerrill/alfred/alfred"
)

var (
	Version = "Development"
	Commit  = ""
)

func main() {
	version := flag.Bool("version", false, "Alfred's version number")
	disableColors := flag.Bool("disable-colors", false, "Disable colors")
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

	NewTask(task, context, tasks)
}
