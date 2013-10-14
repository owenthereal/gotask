package task

import (
	"fmt"
	"io"
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
			t := &T{Args: args}
			task.F(t)
			var writer io.Writer
			if t.failed {
				writer = os.Stderr
			} else {
				writer = os.Stdout
			}

			for _, out := range t.output {
				fmt.Fprintf(writer, "%v", out)
			}

			return
		}
	}

	err = fmt.Errorf("'%s' is not a task", name)
	return
}
