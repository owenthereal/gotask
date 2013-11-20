package tasking

import (
	"flag"
	"fmt"
	"github.com/codegangsta/cli"
	"strings"
)

type TaskSet struct {
	Name       string
	Dir        string
	PkgObj     string
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
	Flags       []Flag
	ActionName  string
	Action      func(*T)
}

func (t *Task) toCLIFlags() (flags []cli.Flag) {
	for _, flag := range t.Flags {
		flags = append(flags, flag)
	}

	return
}

type Flag interface {
	fmt.Stringer
	Apply(*flag.FlagSet)
}

type BoolFlag struct {
	Name  string
	Usage string
}

func (f BoolFlag) String() string {
	return fmt.Sprintf("%s\t%v", f.Name, f.Usage)
}

func (f BoolFlag) Apply(set *flag.FlagSet) {
	for _, name := range f.splitName() {
		set.Bool(name, false, f.Usage)
	}
}

func (f BoolFlag) splitName() (names []string) {
	for _, name := range strings.Split(f.Name, ",") {
		names = append(names, strings.TrimSpace(name))
	}

	return
}
