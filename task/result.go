package task

import (
	"os/exec"
	"syscall"
)

type Result struct {
	Err      error
	ExitCode int
}

func (e *Result) HasError() bool {
	return e.Err != nil
}

func (e *Result) Error() string {
	return e.Err.Error()
}

func newResult(err error) *Result {
	exitCode := 0
	if err != nil {
		exitCode = 1
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
				err = nil
			}
		}
	}

	return &Result{Err: err, ExitCode: exitCode}
}
