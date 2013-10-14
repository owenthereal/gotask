package examples

import (
	"fmt"
	"github.com/jingweno/gotask/task"
	"os/user"
)

// Say hello to current user
//
// Print out hello to current user
func TaskSayHello(t *task.T) {
	user, _ := user.Current()
	fmt.Printf("Hello %s\n", user.Name)
}
