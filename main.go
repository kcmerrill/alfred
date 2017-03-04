package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kcmerrill/alfred/alfred"
)

var (
	Version = "dev"
	Commit  = "n/a"
)

func main() {
	dir := flag.String("dir", ".", "Directory where your alfred files are stored")
	serve := flag.Bool("serve", false, "Start alfred's webserver to share alfred files")
	port := flag.String("port", "8080", "Alfred's webserver port")
	version := flag.Bool("version", false, "Alfred's version number")
	flag.Parse()

	/* Giddy up! */
	if *version {
		fmt.Println()
		fmt.Println("Alfred - Because even Batman needs a little help.")
		fmt.Println("---")
		fmt.Println("Version: ", Version)
		fmt.Println("CommitId: ", Commit)
		fmt.Println("---")
		fmt.Println("Made with <3 by http://kcmerrill.com")
		fmt.Println()
		return
	}
	if *serve {
		alfred.Serve(*dir, *port)
		return
	}

	alfred.New(os.Args)
}
