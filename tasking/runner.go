package tasking

import (
	"github.com/codegangsta/cli"
	"os"
)

func Run(taskSet *TaskSet, args []string) {
	runner := runner{TaskSet: taskSet}
	err := runner.Run(args)
	if err != nil {
		os.Exit(1)
	}
}

type runner struct {
	TaskSet *TaskSet
}

func (r *runner) Run(args []string) error {
	cmds := convertToCommands(r.TaskSet.Tasks)
	app := cli.NewApp()
	app.Name = r.TaskSet.Name
	app.Commands = cmds
	return app.Run(args)
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
	if t.Failed() {
		os.Exit(1)
	}
}
