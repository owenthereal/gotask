package task

import (
	"fmt"
)

func RunTasks(tasks []Task) {
	runner := taskRunner{Tasks: tasks}
	runner.Run()
}

type taskRunner struct {
	Tasks []Task
}

func (r *taskRunner) Run() {
	var tasks []*T
	for _, task := range r.Tasks {
		t := &T{name: task.Name}
		task.F(t)
		tasks = append(tasks, t)
	}

	for _, task := range tasks {
		fmt.Println(task.name)
	}
}
