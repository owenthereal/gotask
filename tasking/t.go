package tasking

import (
	"fmt"
	"github.com/kballard/go-shellquote"
	"os"
	"strings"
	"sync"
)

type T struct {
	mu     sync.RWMutex
	Args   []string
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

func (t *T) fail() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.failed = true
}

func (t *T) Failed() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.failed
}

func (t *T) Log(args ...interface{}) {
	fmt.Println(args...)
}

func (t *T) Logf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (t *T) Error(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
	t.fail()
}

func (t *T) Errorf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	t.fail()
}
