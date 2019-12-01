package server

import (
	"strings"
	"testing"
)

var stackRunnerError = `2019-12-01 11:52:30.871603: [info] Preprocessing library for HSRest-0.1.0.0..
2019-12-01 11:52:30.871683: [info] Building library for HSRest-0.1.0.0..
2019-12-01 11:52:31.202341: [info] Preprocessing executable 'HSRest-exe' for HSRest-0.1.0.0..
2019-12-01 11:52:31.202696: [info] Building executable 'HSRest-exe' for HSRest-0.1.0.0..
2019-12-01 11:52:31.296355: [info] [2 of 2] Compiling Main
2019-12-01 11:52:31.302871: [warn] 
2019-12-01 11:52:31.304038: [warn] /home/abhi/Projects/stack/HSRest/app/Main.hs:6:8: error:
2019-12-01 11:52:31.304159: [warn]     • Variable not in scope: someFuncs :: IO ()
2019-12-01 11:52:31.304238: [warn]     • Perhaps you meant ‘someFunc’ (imported from Lib)
2019-12-01 11:52:31.304305: [warn]   |
2019-12-01 11:52:31.304378: [warn] 6 | main = someFuncs
2019-12-01 11:52:31.304442: [warn]   |        ^^^^^^^^^
2019-12-01 11:52:31.337118: [warn] 
2019-12-01 11:52:31.337345: [debug] Start: getPackageFiles /home/abhi/Projects/stack/HSRest/HSRest.cabal
2019-12-01 11:52:31.340887: [debug] Finished in 3ms: getPackageFiles /home/abhi/Projects/stack/HSRest/HSRest.cabal
2019-12-01 11:52:31.342969: [error] 
--  While building package HSRest-0.1.0.0 using:
      /home/abhi/.stack/setup-exe-cache/x86_64-linux-tinfo6/Cabal-simple_mPHDZzAJ_2.4.0.1_ghc-8.6.5 --builddir=.stack-work/dist/x86_64-linux-tinfo6/Cabal-2.4.0.1 build lib:HSRest exe:HSRest-exe --ghc-options " -fdiagnostics-color=always"
    Process exited with code: ExitFailure 1
`

var filteredErrLines = []string{
	"2019-12-01 11:52:31.302871: [warn] ",
	"2019-12-01 11:52:31.304038: [warn] /home/abhi/Projects/stack/HSRest/app/Main.hs:6:8: error:",
	"2019-12-01 11:52:31.304159: [warn]     • Variable not in scope: someFuncs :: IO ()",
	"2019-12-01 11:52:31.304238: [warn]     • Perhaps you meant ‘someFunc’ (imported from Lib)",
	"2019-12-01 11:52:31.304305: [warn]   |",
	"2019-12-01 11:52:31.304378: [warn] 6 | main = someFuncs",
	"2019-12-01 11:52:31.304442: [warn]   |        ^^^^^^^^^",
	"2019-12-01 11:52:31.337118: [warn] ",
}

// TestFilter tests the Filter function
func TestFilter(t *testing.T) {
	lst := []string{"apple", "banana", "bagels"}
	res := Filter(lst, func(s string) bool {
		return len(s) > 0 && s[0] == 'b'
	})
	if len(res) != 2 {
		t.Fail()
	} else if strings.Compare(res[0], "banana") != 0 ||
		strings.Compare(res[1], "bagels") != 0 {
		t.Fail()
	}
}

// TestEmptyUnfoldedErrLines test making an empty sturcture
func TestEmptyUnfoldedErrLines(t *testing.T) {
	str := EmptyUnfoldedErrLines()
	if len(str.CurrentSet) != 0 ||
		len(str.List) != 0 {
		t.Fail()
	}
}

// TestFilterErrLines proper splitting and filtering of error lines
func TestFilterErrLines(t *testing.T) {
	parsedStackRunnerErr, err := FilterErrLines(&stackRunnerError, nil)
	if err != nil {
		t.Fail()
	} else if len(parsedStackRunnerErr) != 8 {
		t.Fail()
	}
}

// TestUnfoldErrLines tests unfolding of filtered error lines
func TestUnfoldErrLines(t *testing.T) {
	unfoldedErrorLines, err := UnfoldErrLines(filteredErrLines, nil)
	if err != nil {
		t.Fail()
	} else if unfoldedErrorLines == nil {
		t.Fail()
	} else if len(unfoldedErrorLines.CurrentSet) > 0 ||
		len(unfoldedErrorLines.List) == 0 {
		t.Fail()
	}
}
