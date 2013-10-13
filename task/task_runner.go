package task

import (
	"fmt"
	"os"
)

func RunTasks(tasks []Task) {
	runner := taskRunner{Tasks: tasks, Args: os.Args[1:]}
	runner.Run()
}

type taskRunner struct {
	Tasks []Task
	Args  []string
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
