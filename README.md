# gotask

Build tool in Go.

## Overview

`gotask` is a simple build tool designed for Go.
It provides a convention-over-configuration way of writing build tasks in Go.
`gotask` is heavily inspired by [`go test`](http://golang.org/pkg/testing).

## Defining a Task

Similar to writing a Go test, create a file name `TASK_NAME_task.go` and name the task function in the
format of

```go
// Usage
//
// Description
func TestXxx(*task.T) {
  ...
}
```

where `Xxx` can be any alphanumeric string (but the first letter must not be in [a-z]) and serves to identify the task routine.
The comments for the task function will be automatically parsed as the task's usage and help description:
The first block of the comment is the usage and the rest is description.

## Installation

```plain
$ go get -u github.com/jingweno/gotask
```

## Examples

Create a file called `say_hello_task.go` and paste in the following content:

```go
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
```

By convention, `gotask` is able to discover the task and dasherize the name:

```plain
$ gotask -h
NAME:
   gotask - Build tool in Go

USAGE:
   gotask [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   say-hello    Say hello to current user
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --compile, -c        compile the task binary to pkg.task but do not
run it
   --version            print the version
   --help, -h           show help
```

The first block of the comment appear as the task usage for
`say-hello`. The rest become the description:

```plain
$ gotask say-hello -h
NAME:
   say-hello - Say hello to current user

USAGE:
   command say-hello [command options] [arguments...]

DESCRIPTION:
   Print out hello to current user

OPTIONS:
```

To execute the task, type:

```plain
$ gotask say-hello
Hello Owen Ou
```
