// staticlint — кастомный линтер,
// который включает в себя все анализаторы, начинающиеся на SA из https://staticcheck.io/docs/checks/.
// А также ST1012 https://staticcheck.io/docs/checks/#ST1012 и ST1013 https://staticcheck.io/docs/checks/#ST1013.
// Дальше включает себя анализаторы printf, shadow, structtag из модуля golang.org/x/tools/go/analysis/passes.
// И в конце списка — кастомный анализатор OsExitChecker, который проверяет наличие os.Exit в main.go
package main
