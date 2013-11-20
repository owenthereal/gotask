# gotask [![Build Status](https://travis-ci.org/jingweno/gotask.png?branch=master)](https://travis-ci.org/jingweno/gotask)

A convention-over-configuration build tool in Go.

## Motivation

To write build tasks on a Go project in Go instead of Make, Rake or insert your build tool here.

## Overview

`gotask` is a simple build tool designed for Go.
It provides a convention-over-configuration way of writing build tasks in Go.
`gotask` is heavily inspired by [`go test`](http://golang.org/pkg/testing).

## Installation

```plain
$ go get -u github.com/jingweno/gotask
```

## Defining a Task

Similar to defining a Go test, create a file called `TASK_NAME_task.go` and name the task function in the
format of

```go
// +build gotask

package main

import "github.com/jingweno/gotask/tasking"

// NAME
//    The name of the task - a one-line description of what it does
//
// DESCRIPTION
//    A textual description of the task function
//
// OPTIONS
//    Definition of what command line options it takes
func TaskXxx(t *tasking.T) {
  ...
}
```

where `Xxx` can be any alphanumeric string (but the first letter must not be in [a-z]) and serves to identify the task name.

### Task Name

Without declaring the [task name in the comments](https://github.com/jingweno/gotask#comments-as-man-page),
`gotask` will dasherize the `Xxx` part of the task function name and use it as the task name.

### Comments as Man Page

The comments for the task function are parsed as the task's man page by following the [man page layout](http://en.wikipedia.org/wiki/Man_page#Layout):
Section NAME contains the name of the task and a one-line description of what it does, separated by a "-".
Section DESCRIPTION contains the textual description of the task function.
Section OPTIONS contains the definition of the command line flags it takes.

### Build Tags

The `// +build gotask` [build tag](http://golang.org/pkg/go/build/#Context) constraints task functions to `gotask` build only.
Without the build tag, task functions will be available to application build which may not be desired.

## Compiling Tasks

`gotask` is able to compile defined tasks into an executable using `go build`.
This is useful when you need to distribute your build executables.
See `gotask -c` for details.

## Task Scaffold

`gotask` is able to generate a task scaffold named `pkg_task.go` with the `--generate` or `-g` flag.
`pkg` is the name of the package where `gotask` is run:

```plain
// in a folder where package example is defined
$ gotask -g
create example_task.go
```

## Examples

On a [Go project](http://golang.org/doc/code.html#Organization), create a file called `say_hello_task.go` with the following content:

```go
// +build gotask

package main

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
```

Make sure the build tag `// +build gotask` is the first line of the file and there's an empty line before package definition.
The comments of the task should be in the format of the [man page layout](http://en.wikipedia.org/wiki/Man_page#Layout).
Running `gotask -h` displays all the tasks:

```plain
$ gotask -h
NAME:
   gotask - Build tool in Go

USAGE:
   gotask [global options] command [command options] [arguments...]

VERSION:
   0.8.0

OMMANDS:
   say-hello    Say hello to current user
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --generate, -g       generate a task scaffold named pkg_task.go
   --compile, -c        compile the task binary to pkg.task but do not run it
   --debug              run in debug mode
   --version            print the version
   --help, -h           show help
```

Noticing section NAME of the comments appears as the task name and usage for
`say-hello`, section DESCRIPTION becomes the description, section OPTIONS becomes the options:

```plain
$ gotask say-hello -h
NAME:
   say-hello - Say hello to current user

USAGE:
   command say-hello [command options] [arguments...]

DESCRIPTION:
   Print out hello to current user

OPTIONS:
   --verbose, -v        run in verbose mode
   --debug              run in debug mode
```

To execute the task, type:

```plain
$ gotask say-hello
Hello Owen Ou
```

To execute the task in verbose mode, type:

```plain
$ gotask say-hello -v
Hello Owen Ou, the time now is 2013-11-20 15:32:00.73771438 -0800 PST
```

To compile the task into an executable named `pkg.task` where pkg is the
last segment of the import path using `go build`, type:

```plain
$ gotask -c
```

More [examples](https://github.com/jingweno/gotask/tree/master/examples) are available.

## License

`gotask` is released under the MIT license. See [LICENSE.md](https://github.com/jingweno/gotask/blob/master/LICENSE.md).
