package main

import (
	"fmt"
	"os"
	"stackrunner_server/server"
)

func main() {
	args := os.Args
	if l := len(args); l < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <path> [args...]\n", args[0])
	} else {
		uerrLines, err := server.RunStack(args[1], args[2:])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			fmt.Println(uerrLines)
		}
	}
}
