package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// ReadFile wraps os.Open.
// Panics on error.
func ReadFile(name string) *os.File {
	pwd, _ := os.Getwd()
	path := filepath.Join(pwd, name)
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	return file
}

// Print to console.
func Print(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}
