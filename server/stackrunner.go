package server

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"regexp"
	"strings"
)

// Piper is a basic pipe logger
type Piper struct {
	Data bytes.Buffer
}

func (ppr *Piper) Write(p []byte) (int, error) {
	n, err := fmt.Fprintf(&(ppr.Data), "%s", p)
	return n, err
}

// UnfoldedErrorLines is a cumulative error return type
type UnfoldedErrorLines struct {
	CurrentSet []string
	List       [][]string
}

func changeDir(path string) error {
	return os.Chdir(path)
}

// Filter filters a list of types using a predicate
func Filter(lst []string, pred func(string) bool) []string {
	res := make([]string, 0)
	for _, e := range lst {
		if pred(e) {
			res = append(res, e)
		}
	}
	return res
}

func emptyUnfoldedErrLines() UnfoldedErrorLines {
	cs := make([]string, 0)
	ls := make([][]string, 0)
	uerrLines := UnfoldedErrorLines{cs, ls}
	return uerrLines
}

func unfoldedErrLines(filteredErrLines []string) (UnfoldedErrorLines, error) {
	re, err := regexp.Compile("^.+\\[warn\\] $")
	uerrLines := emptyUnfoldedErrLines()
	if err == nil {
		for _, l := range filteredErrLines {
			if re.MatchString(l) {
				if len(uerrLines.CurrentSet) > 0 {
					uerrLines.List = append(uerrLines.List, uerrLines.CurrentSet)
					uerrLines.CurrentSet = make([]string, 0)
				}
			} else {
				uerrLines.CurrentSet = append(uerrLines.CurrentSet, l)
			}
		}
	}
	return uerrLines, err

}

func parseStackRunnerErr(stderr string) (UnfoldedErrorLines, error) {
	re, err := regexp.Compile(".+\\[warn\\].+")
	uerrLines := emptyUnfoldedErrLines()
	if err == nil {
		errLinesString := strings.Replace(stderr, "\r", "", -1)
		errLines := strings.Split(errLinesString, "\n")
		filteredErrLines := Filter(errLines, func(s string) bool {
			return re.MatchString(s)
		})
		uerrLines, err = unfoldedErrLines(filteredErrLines)
	}
	return uerrLines, err
}

// RunStack runs the stack build command to validate the project
func RunStack(path string, args []string) (UnfoldedErrorLines, error) {
	uerrLines := emptyUnfoldedErrLines()
	err := changeDir(path)
	if err == nil {
		cmd := exec.Command("stack", "build", "--fast", "--verbose")
		pprstdout := Piper{}
		pprstderr := Piper{}
		cmd.Stdout = &pprstdout
		cmd.Stderr = &pprstderr
		err = cmd.Run()
		if err != nil {
			stderr := pprstderr.Data.String()
			// fmt.Println(stderr)
			// stdout := pprstdout.Data.String()
			uerrLines, err = parseStackRunnerErr(stderr)
		}
	}
	return uerrLines, err
}
