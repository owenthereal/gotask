package main

import (
	"github.com/jingweno/gotask/task"
	"os"
)

func main() {
	app := task.NewApp()
	app.Run(os.Args)
}
