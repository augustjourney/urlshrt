package osexit

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

var OsExitChecker = &analysis.Analyzer{
	Name: "osexitchecker",
	Doc:  "check for os exit",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	packageName := pass.Pkg.Name()

	// проверяем только пакет main
	if packageName != "main" {
		return nil, nil
	}

	for _, file := range pass.Files {
		// проверяем только файл main
		if file.Name.String() != "main" {
			continue
		}

		ast.Inspect(file, func(n ast.Node) bool {
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			if fn, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
				if ident, ok := fn.X.(*ast.Ident); ok {
					if ident.Name == "os" && fn.Sel.Name == "Exit" {
						pass.Reportf(callExpr.Pos(), "should not use os.Exit in main")
					}
				}
			}

			return true
		})

	}

	return nil, nil
}
