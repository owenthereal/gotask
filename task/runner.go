package task

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Run(args []string) (result *Result) {
	runner := Runner{Args: args}
	err := runner.Run()
	result = newResult(err)
	return
}

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
	if err != nil {
		return
	}

	exec, err := buildTaskMain(source, file)
	if err != nil {
		return
	}

	err = runTaskMain(exec, r.Args)
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

func buildTaskMain(sourceDir, mainFile string) (exec string, err error) {
	taskDir := filepath.Dir(mainFile)

	err = os.Chdir(taskDir)
	if err != nil {
		return
	}

	// TODO: consider caching build
	err = execCmd("go", "build")
	if err != nil {
		return
	}

	err = os.Chdir(sourceDir)
	if err != nil {
		return
	}

	files, err := ioutil.ReadDir(taskDir)
	if err != nil {
		return
	}

	execPrefix := filepath.Base(taskDir)
	for _, file := range files {
		if strings.HasPrefix(file.Name(), execPrefix) {
			exec = filepath.Join(taskDir, file.Name())
			return
		}
	}

	err = fmt.Errorf("can't build task main %s", mainFile)
	return
}

func runTaskMain(exec string, args []string) (err error) {
	cmd := []string{exec}
	cmd = append(cmd, args...)
	err = execCmd(cmd...)
	return
}
