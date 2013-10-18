package tasking

import (
	"fmt"
	"github.com/codegangsta/cli"
	"io"
	"os"
)

func Run(taskSet *TaskSet, args []string) {
	runner := runner{TaskSet: taskSet}
	runner.Run(args)
}

type runner struct {
	TaskSet *TaskSet
}

func (r *runner) Run(args []string) {
	cmds := convertToCommands(r.TaskSet.Tasks)
	app := cli.NewApp()
	app.Name = r.TaskSet.Name
	app.Commands = cmds
	app.Run(args)
	return
}

func convertToCommands(tasks []Task) (cmds []cli.Command) {
	for _, t := range tasks {
		task := t
		cmd := cli.Command{
			Name:        task.Name,
			Usage:       task.Usage,
			Description: task.Description,
			Action: func(c *cli.Context) {
				runTask(task, c.Args())
			},
		}

		cmds = append(cmds, cmd)
	}

	return
}

func runTask(task Task, args []string) {
	t := &T{Args: args}
	task.Action(t)
	var writer io.Writer
	if t.failed {
		writer = os.Stderr
	} else {
		writer = os.Stdout
	}

	for _, out := range t.output {
		fmt.Fprintf(writer, "%v", out)
	}
}
