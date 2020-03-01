package myerror

import (
	"fmt"
	"os"
	"os/user"
	"runtime"

	"github.com/rsmaxwell/page/internal/version"
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

	e.lines = append(e.lines, "page")

	pc, fn, linenumber, _ := runtime.Caller(1)
	e.lines = append(e.lines, fmt.Sprintf("Function: %s", runtime.FuncForPC(pc).Name()))
	e.lines = append(e.lines, fmt.Sprintf("File:Line: %s:%d", fn, linenumber))

	e.lines = append(e.lines, "Version: "+version.Version())

	user, _ := user.Current()
	e.lines = append(e.lines, "Username: "+user.Username)
	e.lines = append(e.lines, "HomeDir: "+user.HomeDir)

	dir, _ := os.Getwd()
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
