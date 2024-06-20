package main

import (
	"os"
)

func main() {
	os.Exit(1) // want "should not use os.Exit in main"
}
