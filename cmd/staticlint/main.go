package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
	"staticlint/osexit"
	"strings"
)

func main() {
	var checks []*analysis.Analyzer

	for _, v := range staticcheck.Analyzers {
		// добавляем все проверки, которые начинаются на SA
		if strings.HasPrefix(v.Analyzer.Name, "SA") {
			checks = append(checks, v.Analyzer)
		}

		// добавляем проверку ST1012
		// Poorly chosen name for error variable
		if v.Analyzer.Name == "ST1012" {
			checks = append(checks, v.Analyzer)
		}

		// добавляем проверку ST1013
		// Should use constants for HTTP error codes, not magic numbers
		// частая привычка писать цифры, вместо использования констант
		if v.Analyzer.Name == "ST1013" {
			checks = append(checks, v.Analyzer)
		}
	}

	// добавляем дефолтные анализаторы из модуля passes
	checks = append(checks, printf.Analyzer, shadow.Analyzer, structtag.Analyzer)

	// добавляем кастомный анализатор для os.exit
	checks = append(checks, osexit.OsExitChecker)

	// запускаем мультичекер
	multichecker.Main(checks...)
}
