package osexit

import (
	"golang.org/x/tools/go/analysis/analysistest"
	"testing"
)

func TestOsExitChecker(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), OsExitChecker, "./...")
}
