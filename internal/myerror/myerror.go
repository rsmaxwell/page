package myerror

import (
	"fmt"
	"os"

	"github.com/rsmaxwell/page/internal/basic/version"
)

// MyError is the error structure
type MyError struct {
	lines []string
}

// New created a new error
func New(line string) MyError {
	e := MyError{}
	l := make([]string, 3)
	e.lines = l

	dir, _ := os.Getwd()
	e.lines = append(e.lines, "page, version: "+version.Version())
	e.lines = append(e.lines, "Current Working Directory: "+dir)

	return e.Add(line)
}

// Add function
func (e MyError) Add(line string) MyError {
	e.lines = append(e.lines, line)
	return e
}

// Handle function
func (e MyError) Handle() {
	for _, line := range e.lines {
		fmt.Println("<p>" + line + "</p>")
	}
}
