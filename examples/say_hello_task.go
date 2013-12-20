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
//    -n, --name=NAME
//        say hello to an user with the given NAME
//    -g, --greeting=TYPE
//        say hello using the given TYPE of greeting
//    -v, --verbose
//        run in verbose mode
func TaskSayHello(t *tasking.T) {
	username := t.Flags.String("n")
	if username == "" {
		user, _ := user.Current()
		username = user.Name
	}
	if t.Flags.Bool("v") || t.Flags.Bool("verbose") {
		t.Logf("Hello %s, the time now is %s\n", username, time.Now())
	} else {
		t.Logf("Hello %s\n", username)
	}
}
