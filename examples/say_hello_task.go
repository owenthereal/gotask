// +build gotask

package examples

import (
	"github.com/jingweno/gotask/tasking"
	"os/user"
	"time"
)

// NAME
//    say-hello - Say hello to current user
//
// DESCRIPTION
//    Print out hello to current user
//
// OPTIONS
//    --verbose, -v
//        run in verbose mode
func TaskSayHello(t *tasking.T) {
	user, _ := user.Current()
	if t.Flags.Bool("v") || t.Flags.Bool("verbose") {
		t.Logf("Hello %s, the time now is %s\n", user.Name, time.Now())
	} else {
		t.Logf("Hello %s\n", user.Name)
	}
}
