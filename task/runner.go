package task

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Runner struct {
	Args []string
}

func (r *Runner) Run() (err error) {
	source, err := os.Getwd()
	if err != nil {
		return
	}

	loader := taskParser{source}
	funcs, err := loader.Parse()
	if err != nil {
		return
	}

	if !funcs.HasTasks() {
		err = fmt.Errorf("%s\t[no task files]", funcs.ImportPath)
		return
	}

	// create temp work dir
	work, err := ioutil.TempDir("", "go-task")
	if err != nil {
		return
	}
	defer func() {
		os.RemoveAll(work)
	}()

	file, err := writeTaskMain(work, funcs)
	execCmd("go", "run", file)

	return
}

func writeTaskMain(work string, funcs *taskFuncs) (file string, err error) {
	// create task dir
	taskDir := filepath.Join(work, filepath.FromSlash(funcs.ImportPath))
	err = os.MkdirAll(taskDir, 0777)
	if err != nil {
		return
	}

	// create main.go
	file = filepath.Join(taskDir, "main.go")
	f, err := os.Create(file)
	if err != nil {
		return
	}
	defer f.Close()

	// write to main.go
	w := mainWriter{funcs}
	err = w.Write(f)

	return
}
