package server

import (
	"strings"
	"testing"
)

var buildPlanErrorLines = `2019-12-02 08:14:08.881154: [debug] Process finished in 24ms: /home/abhi/.stack/programs/x86_64-linux/ghc-tinfo6-8.6.5/bin/ghc-pkg-8.6.5 --user --no-user-package-db --package-db /home/abhi/Projects/stack/HSRest/.stack-work/install/x86_64-linux-tinfo6/8d9a83e26c843fdf600d1c9d454326c9304938424b9ec126ba388061b5342e58/8.6.5/pkgdb dump --expand-pkgroot
2019-12-02 08:14:08.882044: [debug] Constructing the build plan
2019-12-02 08:14:08.886901: [error] 
Error: While constructing the build plan, the following exceptions were encountered:

In the dependencies for HSRest-0.1.0.0:
    bystring needed, but the stack configuration has no specified version (no
             package with that name found, perhaps there is a typo in a
             package's build-depends or an omission from the stack.yaml packages
             list?)
needed since HSRest is a build target.

Some different approaches to resolving this:


2019-12-02 08:14:08.887601: [error] Plan construction failed.`

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

var unfoldedErrorLines = UnfoldedErrorLines{nil, [][]string{filteredErrLines[1:]}}

var errorLineToParse = "2019-12-01 11:52:31.304038: [warn] /home/abhi/Projects/stack/HSRest/app/Main.hs:6:8: error:"

const errorFile = "/home/abhi/Projects/stack/HSRest/app/Main.hs"
const errorLine = 6
const errorColumn = 8
const errorExtras = ""

func buildErrorTestHelper(be *BuildError, cb func(bool)) {

	fileFail := strings.Compare(be.File, errorFile) != 0
	lineFail := be.Line != errorLine
	columnFail := be.Column != errorColumn
	extrasFail := strings.Compare(be.Extras, errorExtras) != 0

	fail := fileFail || lineFail || columnFail || extrasFail
	cb(fail)
}

func TestContainsBuildPlanErrors(t *testing.T) {
	cbpe, bpe, err := ContainsBuildPlanErrors(&buildPlanErrorLines, nil)
	if err != nil {
		t.FailNow()
	} else if !cbpe {
		t.FailNow()
	} else if bpe == nil {
		t.FailNow()
	}
}

func TestParseErrorLine(t *testing.T) {
	be, err := ParseErrorLine(errorLineToParse, nil)
	if err != nil ||
		be == nil {
		t.FailNow()
	}
	buildErrorTestHelper(be, func(fail bool) {
		if fail {
			t.FailNow()
		}
	})
}

func TestParseUnfoldedErrorLines(t *testing.T) {
	bes, err := ParseUnfoldedErrorLines(&unfoldedErrorLines, nil)
	if err != nil || len(bes) == 0 {
		t.FailNow()
	} else {
		be := bes[0]
		buildErrorTestHelper(&be, func(fail bool) {
			if fail {
				t.FailNow()
			}
		})
	}
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
