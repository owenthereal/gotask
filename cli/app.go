package cli

import (
	"flag"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/jingweno/gotask/build"
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
	app.Action = func(c *cli.Context) {
		if c.Bool("c") || c.Bool("compile") {
			err := build.CompileAndRun(c.Args(), true)
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

	for _, task := range taskSet.Tasks {
		t := task
		cmd := cli.Command{
			Name:        task.Name,
			Usage:       task.Usage,
			Description: task.Description,
			Action: func(c *cli.Context) {
				a := []string{t.Name}
				a = append(a, c.Args()...)
				err := build.CompileAndRun(a, false)
				if err != nil {
					log.Fatal(err)
				}
			},
		}
		cmds = append(cmds, cmd)
	}

	return
}
