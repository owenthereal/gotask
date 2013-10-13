package main

import (
	"fmt"
	"github.com/jingweno/gotask/task"
	"os"
)

func main() {
	result := task.Run(os.Args[1:])
	if result.HasError() {
		fmt.Fprintln(os.Stderr, result)
	}
	os.Exit(result.ExitCode)
}
