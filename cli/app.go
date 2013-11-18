package cli

import (
	"github.com/codegangsta/cli"
	"github.com/jingweno/gotask/build"
	"log"
	"os"
)

var debugFlag = cli.BoolFlag{"debug", "run in debug mode"}

func NewApp() *cli.App {
	cmds, err := parseCommands()
	if err != nil {
		log.Fatal(err)
	}

	app := cli.NewApp()
	app.Name = "gotask"
	app.Usage = "Build tool in Go"
	app.Version = Version
	app.Commands = cmds
	app.Flags = []cli.Flag{
		newFlag{Usage: "create an example task file named pkg_task.go"},
		compileFlag{Usage: "compile the task binary to pkg.task but do not run it"},
		debugFlag,
	}
	app.Action = func(c *cli.Context) {
		if c.Bool("n") || c.Bool("new") {
			err := generateNewTask()
			if err != nil {
				log.Fatal(err)
			}

			return
		}

		if c.Bool("c") || c.Bool("compile") {
			err := compileTasks(c.Bool("debug"))
			if err != nil {
				log.Fatal(err)
			}

			return
		}

		if len(c.Args()) == 0 {
			cli.ShowAppHelp(c)
		}
	}

	return app
}

func parseCommands() (cmds []cli.Command, err error) {
	source, err := os.Getwd()
	if err != nil {
		return
	}

	parser := build.NewParser()
	taskSet, err := parser.Parse(source)
	if err != nil {
		return
	}

	for _, t := range taskSet.Tasks {
		task := t
		cmd := cli.Command{
			Name:        task.Name,
			Usage:       task.Usage,
			Description: task.Description,
			Flags: []cli.Flag{
				debugFlag,
			},
			Action: func(c *cli.Context) {
				args := []string{task.Name}
				args = append(args, c.Args()...)
				err := runTasks(args, c.Bool("debug"))
				if err != nil {
					log.Fatal(err)
				}
			},
		}

		cmds = append(cmds, cmd)
	}

	return
}

func runTasks(args []string, isDebug bool) (err error) {
	sourceDir, err := os.Getwd()
	if err != nil {
		return
	}

	err = build.Run(sourceDir, args, isDebug)
	return
}
