package task

import (
	"flag"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/jingweno/gotask/tasking"
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
	Action      func(*tasking.T)
}

func (t *Task) ToCLIFlags() (flags []cli.Flag) {
	for _, flag := range t.Flags {
		flags = append(flags, flag)
	}

	return
}

type Flag interface {
	fmt.Stringer
	Apply(*flag.FlagSet)
	DefType(importAsPkg string) string
}

type BoolFlag struct {
	Name  string
	Usage string
}

func (f BoolFlag) String() string {
	return fmt.Sprintf("%s\t%v", strings.Join(f.splitName(), ", "), f.Usage)
}

func (f BoolFlag) Apply(set *flag.FlagSet) {
	for _, name := range f.splitName() {
		set.Bool(strings.TrimLeft(name, "-"), false, f.Usage)
	}
}

func (f BoolFlag) DefType(importAsPkg string) string {
	return fmt.Sprintf(`%s.BoolFlag{Name: "%s", Usage: "%s"}`, importAsPkg, f.Name, f.Usage)
}

func (f BoolFlag) splitName() (names []string) {
	for _, name := range strings.Split(f.Name, ",") {
		names = append(names, strings.TrimSpace(name))
	}

	return
}
