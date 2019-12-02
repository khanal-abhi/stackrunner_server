package main

import (
	"encoding/json"
	"fmt"
	"os"
	"stackrunner_server/server"
	"strings"
)

func main() {
	args := os.Args
	if l := len(args); l < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <path> [args...]\n", args[0])
	} else {
		buildErrors, err := server.RunStack(args[1], args[2:])
		if err != nil {
			if strings.Contains(err.Error(), "executable file not found in $PATH") {
				be := server.BuildError{
					File:    "",
					Line:    -1,
					Column:  -1,
					Details: nil, Extras: strings.Trim(strings.Replace(err.Error(), "exec:", "", 1), " "),
				}
				dt, err := json.Marshal([]server.BuildError{be})
				if err == nil {
					fmt.Println(string(dt))
				}
			}
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
