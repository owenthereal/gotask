# gotask [![Build Status](https://travis-ci.org/jingweno/gotask.png?branch=master)](https://travis-ci.org/jingweno/gotask)

A convention-over-configuration build tool in Go.

## Motivation

To write build tasks on a Go project in Go instead of Make, Rake or insert your build tool here.

## Overview

`gotask` is a simple build tool designed for Go.
It provides a convention-over-configuration way of writing build tasks in Go.
`gotask` is heavily inspired by [`go test`](http://golang.org/pkg/testing).

## Defining a Task

Similar to defining a Go test, create a file called `TASK_NAME_task.go` and name the task function in the
format of

```go
// +build gotask

package main

import "github.com/jingweno/gotask/tasking"

// Usage
//
// Description
func TaskXxx(t *tasking.T) {
  ...
}
```

where `Xxx` can be any alphanumeric string (but the first letter must not be in [a-z]) and serves to identify the task routine.
The comments for the task function will be automatically parsed as the task's usage and help description:
The first block of the comment is the usage and the rest is the description.
The `// +build gotask` [build tags](http://golang.org/pkg/go/build/#Context) constraint task functions to `gotask` build.

## Compiling Tasks

`gotask` is able to compile defined tasks into an executable using `go build`.
This is useful when you need to distribute your build executables.
See `gotask -c` for details.

## Installation

```plain
$ go get -u github.com/jingweno/gotask
```

## Examples

On a [Go project](http://golang.org/doc/code.html#Organization), create a file called `say_hello_task.go` with the following content:

```go
// +build gotask

package main

import (
	"fmt"
	"github.com/jingweno/gotask/tasking"
	"os/user"
)

// Say hello to current user
//
// Print out hello to current user
func TaskSayHello(t *tasking.T) {
	user, _ := user.Current()
	fmt.Printf("Hello %s\n", user.Name)
}
```

Make sure the build tag `// +build gotask` is the first line of the file and there's an empty line before package definition.
Without the build tag, task functions will be available to the application build which may not be desired.
By convention, the `gotask` CLI is able to discover the task and dasherize the task name.
Running `gotask -h` displays all the tasks:

```plain
$ gotask -h
NAME:
   gotask - Build tool in Go

USAGE:
   gotask [global options] command [command options] [arguments...]

VERSION:
   0.0.2

COMMANDS:
   say-hello    Say hello to current user
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --compile, -c        compile the task binary to pkg.task but do not run it
   --version            print the version
   --help, -h           show help
```

Noticing the first block of the comments appears as the task usage for
`say-hello`, the rest become the description:

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

To compile the task into an executable named `pkg.task` where pkg is the
last segment of the import path using `go build`, type:

```plain
$ gotask -c
```

More [examples](https://github.com/jingweno/gotask/tree/master/examples) are available.

## License

`gotask` is released under the MIT license. See [LICENSE.md](https://github.com/jingweno/gotask/blob/master/LICENSE.md).
