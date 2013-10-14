package task

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type compiler struct {
	currentDir string
	workDir    string
	TaskSet    *TaskSet
}

func (c *compiler) Compile(outDir string) (execFile string, err error) {
	file, err := writeTaskMain(c.workDir, c.TaskSet)
	if err != nil {
		return
	}

	execFile, err = compileTaskMain(c.currentDir, file)
	return
}

func writeTaskMain(work string, taskSet *TaskSet) (file string, err error) {
	// create task dir
	taskDir := filepath.Join(work, filepath.FromSlash(taskSet.ImportPath))
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
	w := mainWriter{taskSet}
	err = w.Write(f)

	return
}

func compileTaskMain(sourceDir, mainFile string) (exec string, err error) {
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

func compileAndRun(args []string) (err error) {
	source, err := os.Getwd()
	if err != nil {
		return
	}

	parser := taskParser{source}
	taskSet, err := parser.Parse()
	if err != nil {
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

	compiler := compiler{currentDir: source, workDir: work, TaskSet: taskSet}
	execFile, err := compiler.Compile("")
	if err != nil {
		return
	}

	runner := runner{execFile}
	err = runner.Run(args)
	return
}
