package task

import (
	"fmt"
)

func RunTasks(funcs []func(*T)) {
	runner := taskRunner{Funcs: funcs}
	runner.Run()
}

type taskRunner struct {
	Funcs []func(*T)
}

func (r *taskRunner) Run() {
	var tasks []*T
	for _, f := range r.Funcs {
		t := &T{}
		f(t)
		tasks = append(tasks, t)
	}

	for _, task := range tasks {
		fmt.Println(task.name)
	}
}
