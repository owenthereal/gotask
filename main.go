package main

import (
	"github.com/jingweno/gotask/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Run(os.Args)
}
