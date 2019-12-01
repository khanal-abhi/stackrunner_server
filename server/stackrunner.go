package server

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"regexp"
	"strings"
)

/**
BEGIN Private declarations
*/

func changeDir(path string) error {
	return os.Chdir(path)
}

func generateStackCommand(eerr error) (*exec.Cmd, *Piper, *Piper) {
	if eerr != nil {
		return nil, nil, nil
	}
	cmd := exec.Command("stack", "build", "--fast", "--verbose")
	pprstdout := Piper{}
	pprstderr := Piper{}
	cmd.Stdout = &pprstdout
	cmd.Stderr = &pprstderr
	return cmd, &pprstdout, &pprstderr
}

/**
END Private declarations
*/

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

// EmptyUnfoldedErrLines generates an empty structure
func EmptyUnfoldedErrLines() UnfoldedErrorLines {
	cs := make([]string, 0)
	ls := make([][]string, 0)
	uerrLines := UnfoldedErrorLines{cs, ls}
	return uerrLines
}

// UnfoldErrLines returns the unfolded lines of errors
func UnfoldErrLines(filteredErrLines []string, eerr error) (*UnfoldedErrorLines, error) {
	if eerr != nil {
		return nil, eerr
	}
	re, err := regexp.Compile("^.+\\[warn\\] $")
	uerrLines := EmptyUnfoldedErrLines()
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
		if len(uerrLines.CurrentSet) > 0 {
			uerrLines.List = append(uerrLines.List, uerrLines.CurrentSet)
			uerrLines.CurrentSet = nil
		}
	}
	return &uerrLines, err

}

// FilterErrLines splits the stderr by newlines and filters them
func FilterErrLines(ls *string, derr error) ([]string, error) {
	if derr != nil {
		return nil, derr
	}
	filteredErrLines := make([]string, 0)
	re, err := regexp.Compile(".+\\[warn\\].+")
	if err == nil {
		errLinesString := strings.Replace(*ls, "\r", "", -1)
		errLines := strings.Split(errLinesString, "\n")
		filteredErrLines = Filter(errLines, func(s string) bool {
			return re.MatchString(s)
		})
	}
	return filteredErrLines, err
}

// RunStack runs the stack build command to validate the project
func RunStack(path string, args []string) (*UnfoldedErrorLines, error) {
	uerrLines := EmptyUnfoldedErrLines()
	err := changeDir(path)
	cmd, _, pprstderr := generateStackCommand(err)
	err1 := cmd.Run()
	if err1 == nil {
		return &uerrLines, err
	}
	var ls *string = nil
	if pprstderr != nil {
		d := pprstderr.Data.String()
		ls = &d
	}
	filteredErrLines, err := FilterErrLines(ls, err)
	unsafeUerrlines, err := UnfoldErrLines(filteredErrLines, err)
	uerrLines = *unsafeUerrlines
	return &uerrLines, err
}
