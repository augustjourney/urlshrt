package osexit

import (
	"golang.org/x/tools/go/analysis"
)

var OsExitChecker = &analysis.Analyzer{
	Name: "osexitchecker",
	Doc:  "check for os exit",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	packageName := pass.Pkg.Name()
	if packageName != "main" {
		return nil, nil
	}

	for _, file := range pass.Files {
		if file.Name.String() != "main" {
			continue
		}

	}

	return nil, nil
}
