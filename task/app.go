package task

import (
	"flag"
	"fmt"
	"github.com/codegangsta/cli"
	"log"
	"os"
)

type compileFlag struct {
	Usage string
}

func (f compileFlag) String() string {
	return fmt.Sprintf("--compile, -c\t%v", f.Usage)
}

func (f compileFlag) Apply(set *flag.FlagSet) {
	set.Bool("c", false, f.Usage)
	set.Bool("compile", false, f.Usage)
}

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
		compileFlag{Usage: "compile the task binary to pkg.task but do not run it"},
	}
	//app.Action = func(c *cli.Context) {
	//if len(c.Args()) == 0 {
	//cli.ShowAppHelp(c)
	//}
	//}

	return app
}

func parseCommands() (cmds []cli.Command, err error) {
	source, err := os.Getwd()
	if err != nil {
		return
	}

	parser := taskParser{source}
	taskSet, err := parser.Parse()
	if err != nil {
		return
	}

	for _, task := range taskSet.Tasks {
		t := task
		cmd := cli.Command{
			Name:        task.Name,
			Usage:       task.Usage,
			Description: task.Description,
			Action: func(c *cli.Context) {
				a := []string{t.Name}
				a = append(a, c.Args()...)
				compileAndRun(a)
			},
		}
		cmds = append(cmds, cmd)
	}

	return
}
