package main

import (
	"flag"
	"github.com/kcmerrill/alfred/alfred"
)

func main() {
	dir := flag.String("dir", ".", "Directory where your alfred files are stored")
	serve := flag.Bool("serve", false, "Start alfred's webserver to share alfred files")
	port := flag.String("port", "8080", "Alfred's webserver port")
	flag.Parse()

	/* Giddy up! */
	if *serve {
		alfred.Serve(*dir, *port)
	} else {
		alfred.New()
	}
}
