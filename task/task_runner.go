package task

import (
	"fmt"
	"github.com/codegangsta/cli"
	"io"
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
	cmds := convertToCommands(r.Tasks)
	app := cli.NewApp()
	app.Name = "gotask"
	app.Usage = "Build tool in Go"
	app.Version = "0.0.1"
	app.Commands = cmds
	app.Run(r.Args)
	return
}

func convertToCommands(tasks []Task) (cmds []cli.Command) {
	for _, task := range tasks {
		t := task
		cmd := cli.Command{
			Name:        task.Name,
			Usage:       task.Usage,
			Description: task.Description,
			Action: func(c *cli.Context) {
				runTask(t, c.Args())
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
