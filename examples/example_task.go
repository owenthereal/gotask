package examples

import (
	"fmt"
	"github.com/jingweno/gotask/task"
)

func TaskHelloWorld(t *task.T) {
	t.Describe("Hello world").Name("hello")
	fmt.Println("Hello world")
}
