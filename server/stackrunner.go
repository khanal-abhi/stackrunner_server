package server

import (
	"bytes"
	"encoding/json"
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

func generateStackCommand(eerr error) (*exec.Cmd, *Piper, *Piper, error) {
	if eerr != nil {
		return nil, nil, nil, nil
	}
	_, err := exec.LookPath("stack")
	if err != nil {
		return nil, nil, nil, err
	}
	cmd := exec.Command("stack", "build", "--fast", "--verbose")
	pprstdout := Piper{}
	pprstderr := Piper{}
	cmd.Stdout = &pprstdout
	cmd.Stderr = &pprstderr
	return cmd, &pprstdout, &pprstderr, err
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

// BuildError defines the structure for build errors
type BuildError struct {
	File    string   `json:"file"`
	Line    int      `json:"line"`
	Column  int      `json:"column"`
	Details []string `json:"details"`
	Extras  string   `json:"extras"`
}

// ContainsBuildPlanErrors checks to see if the lines contain build plan errors
func ContainsBuildPlanErrors(lines *string, eerr error) (bool, *string, error) {
	var serr *string = nil
	if eerr != nil {
		return false, serr, eerr
	}
	var err error = nil
	tmp := `Error: While constructing the build plan, the following exceptions were encountered:`
	match := strings.Contains(*lines, tmp)
	if match {
		ss := strings.Split(*lines, tmp)
		l := len(ss)
		if l > 1 {
			serr = &ss[l-1]
		}
	}
	return match, serr, err
}

// ParseErrorLine parses the error line into the error struct
func ParseErrorLine(line string, eerr error) (*BuildError, error) {
	if eerr != nil {
		return nil, eerr
	}
	re, err := regexp.Compile("^.+\\[warn\\] (\\/.+):(\\d+):(\\d+)(: error:.*)")
	if err != nil {
		return nil, err
	}
	b := []byte(`
	{	
		"file": "$1",
		"line": $2,
		"column": $3,
		"details": [],
		"extras": "$4"
	}
	`)
	data := re.ReplaceAll([]byte(line), b)
	be := BuildError{}
	err = json.Unmarshal(data, &be)
	if len(strings.Trim(be.Extras, " ")) > 0 {
		be.Extras = strings.Replace(be.Extras, ": error:", "", -1)
		be.Extras = strings.Trim(be.Extras, " ")
	}
	return &be, err
}

// ParseUnfoldedErrorLines parses unfolded error lines to list of BuildErrors
func ParseUnfoldedErrorLines(unfoldedErrorLines *UnfoldedErrorLines, eerr error) ([]BuildError, error) {
	if eerr != nil {
		return nil, eerr
	}
	bes := make([]BuildError, 0)
	var err error = nil
	for _, errset := range unfoldedErrorLines.List {
		l := len(errset)
		if l > 0 {
			be, err := ParseErrorLine(errset[0], nil)
			if err == nil {
				for i, e := range errset {
					if i > 0 {
						if be.Details == nil {
							be.Details = make([]string, 0)
						}
						be.Details = append(be.Details, e)
					}
				}
				bes = append(bes, *be)
			}
		}
	}
	return bes, err
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
func RunStack(path string, args []string) ([]BuildError, error) {
	err := changeDir(path)
	cmd, _, pprstderr, err := generateStackCommand(err)
	if cmd == nil {
		return nil, err
	}
	err1 := cmd.Run()
	if err1 == nil {
		return nil, err
	}
	var ls *string = nil
	if pprstderr != nil {
		d := pprstderr.Data.String()
		ls = &d
	}
	cbp, bpe, err := ContainsBuildPlanErrors(ls, err)
	if cbp {
		be := BuildError{path + "/package.yaml", 0, 0, strings.Split(*bpe, "\n"), ""}
		return []BuildError{be}, err
	}
	filteredErrLines, err := FilterErrLines(ls, err)
	unsafeUerrlines, err := UnfoldErrLines(filteredErrLines, err)
	uerrLines := *unsafeUerrlines
	buildErrors, err := ParseUnfoldedErrorLines(&uerrLines, err)
	return buildErrors, err
}
