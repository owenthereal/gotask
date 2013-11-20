package tasking

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/kballard/go-shellquote"
	"os"
	"strings"
	"sync"
)

type Flags struct {
	C *cli.Context
}

func (f Flags) Bool(name string) bool {
	return f.C.Bool(name)
}

type T struct {
	mu     sync.RWMutex
	Args   []string
	Flags  Flags
	output []string
	failed bool
}

// Run the system command. If multiple arguments are given, they're concatenated to one command.
//
// Example:
//   t.Exec("ls -ltr")
//   t.Exec("ls", FILE1, FILE2)
func (t *T) Exec(cmd ...string) (err error) {
	toRun := strings.Join(cmd, " ")
	input, err := shellquote.Split(toRun)
	if err != nil {
		return
	}

	err = execCmd(input)

	return
}

// Fail marks the task as having failed but continues execution.
func (t *T) Fail() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.failed = true
}

// Check if the task has failed
func (t *T) Failed() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.failed
}

// Log formats its arguments using default formatting, analogous to Println.
func (t *T) Log(args ...interface{}) {
	fmt.Println(args...)
}

// Logf formats its arguments according to the format, analogous to Printf.
func (t *T) Logf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

// Error is equivalent to Log followed by Fail.
func (t *T) Error(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
	t.Fail()
}

// Errorf is equivalent to Logf followed by Fail.
func (t *T) Errorf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	t.Fail()
}
