package task

import (
	"fmt"
)

type TaskSet struct {
	ImportPath string
	Tasks      []Task
}

func (ts *TaskSet) HasTasks() bool {
	return len(ts.Tasks) > 0
}

type Task struct {
	Name        string
	Usage       string
	Description string
	ActionName  string
	Action      func(*T)
}

type T struct {
	Args   []string
	output []string
	failed bool
}

func (t *T) fail() {
	t.failed = true
}

func (t *T) Log(args ...interface{}) {
	t.log(fmt.Sprintln(args...))
}

func (t *T) Logf(format string, args ...interface{}) {
	t.log(fmt.Sprintf(format, args...))
}

func (t *T) Error(args ...interface{}) {
	t.log(fmt.Sprintln(args...))
	t.fail()
}

func (t *T) Errorf(format string, args ...interface{}) {
	t.log(fmt.Sprintf(format, args...))
	t.fail()
}

func (t *T) log(s string) {
	t.output = append(t.output, s)
}
