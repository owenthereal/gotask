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
	if len(r.Args) == 0 {
		printUsage(r.Tasks)
		result = newResult(nil)
		return
	}

	name := r.Args[0]
	args := r.Args[1:]
	err := execTask(r.Tasks, name, args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		printUsage(r.Tasks)
	}

	result = newResult(err)
	return
}

func printUsage(tasks []Task) {
	for _, task := range tasks {
		fmt.Printf("%s\t%s\n", task.Name, task.Usage)
	}
}

func execTask(tasks []Task, name string, args []string) (err error) {
	for _, task := range tasks {
		if name == task.Name {
			t := &T{}
			task.F(t)
			if t.Err != "" {
				fmt.Fprintln(os.Stderr, t.Err)
			}

			return
		}
	}

	err = fmt.Errorf("'%s' is not a task", name)
	return
}
