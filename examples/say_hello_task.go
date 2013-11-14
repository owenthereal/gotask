// +build gotask

package examples

import (
	"fmt"
	"github.com/jingweno/gotask/tasking"
	"os/user"
)

// NAME
//    say-hello - Say hello to current user
//
// DESCRIPTION
//    Print out hello to current user
func TaskSayHello(t *tasking.T) {
	user, _ := user.Current()
	fmt.Printf("Hello %s\n", user.Name)
}
