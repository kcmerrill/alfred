package main

import (
	"os"

	. "github.com/kcmerrill/alfred/alfred"
)

func main() {
	Initialize(&Alfred{}).Task(NewCLI(os.Args))
}
