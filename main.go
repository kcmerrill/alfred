package main

import (
	"fmt"
	"strconv"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	#https://blog.stevenocchipinti.com/2013/06/removing-previously-printed-lines.html/
	#https://godoc.org/golang.org/x/crypto/ssh/terminal#GetSize
	fd := int(os.Stdout.Fd())
	fmt.Println(terminal.GetSize(fd))
	for x := 0; x < 100; x++ {
		fmt.Print(strconv.Itoa(x) + "\r")
		<-time.After(time.Second)
	}
	fmt.Println("finished!")
}
