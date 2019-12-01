package main

import (
	"encoding/json"
	"fmt"
	"os"
	"stackrunner_server/server"
)

func main() {
	args := os.Args
	if l := len(args); l < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <path> [args...]\n", args[0])
	} else {
		buildErrors, err := server.RunStack(args[1], args[2:])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			if len(buildErrors) > 0 {
				bts, err := json.Marshal(buildErrors)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(string(bts))
				}

			} else {
				fmt.Println("[]")
			}
		}
	}
}
