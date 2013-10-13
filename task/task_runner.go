package task

import (
	"fmt"
	"os"
)

func RunTasks(tasks []Task) *Result {
	runner := taskRunner{Tasks: tasks, Args: os.Args[1:]}
	return runner.Run()
}

type taskRunner struct {
	Tasks []Task
	Args  []string
}

func (r *taskRunner) Run() (result *Result) {
	var tasks []*T
	for _, task := range r.Tasks {
		t := &T{name: task.Name}
		task.F(t)
		tasks = append(tasks, t)
	}

	if len(r.Args) == 0 {
		printUsage(tasks)
		result = newResult(nil)
		return
	}

	e := execTask(tasks, r.Args[0], r.Args[1:])
	result = newResult(e)

	if e == nil {
	} else {
		fmt.Fprintln(os.Stderr, e)
		printUsage(tasks)
	}

	return
}

func printUsage(tasks []*T) {
	for _, task := range tasks {
		fmt.Println(task.name)
	}
}

func execTask(tasks []*T, taskName string, args []string) (err error) {
	for _, task := range tasks {
		if task.name == taskName {
			task.action()
			return
		}
	}

	err = fmt.Errorf("'%s' is not a task", taskName)
	return
}
