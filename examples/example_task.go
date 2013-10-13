package examples

import (
	"fmt"
	"github.com/jingweno/gotask/task"
)

func TaskHelloWorld(t *task.T) {
	t.Name("say-hello").Describe("Say Hello")
	t.Action(func() {
		fmt.Println("Hello world")
	})
}
