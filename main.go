package main

import (
	"fmt"
	"github.com/jingweno/gotask/task"
	"os"
)

func main() {
	runner := task.Runner{Args: os.Args[1:]}
	err := runner.Run()
	exitCode := 0
	if err != nil {
		fmt.Println(err)
		exitCode = 1
	}

	os.Exit(exitCode)
}
