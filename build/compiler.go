package build

import (
	"fmt"
	"github.com/jingweno/gotask/tasking"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type compiler struct {
	sourceDir string
	workDir   string
	TaskSet   *tasking.TaskSet
}

func (c *compiler) Compile(outfile string) (execFile string, err error) {
	file, err := c.writeTaskMain(c.workDir, c.TaskSet)
	if err != nil {
		return
	}

	execFile, err = c.compileTaskMain(c.sourceDir, file, outfile)
	return
}

func (c *compiler) writeTaskMain(work string, taskSet *tasking.TaskSet) (file string, err error) {
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

func (c *compiler) compileTaskMain(sourceDir, mainFile, outfile string) (exec string, err error) {
	taskDir := filepath.Dir(mainFile)

	err = os.Chdir(taskDir)
	if err != nil {
		return
	}

	// TODO: consider caching build
	compileCmd := []string{"go", "build", "--tags", "gotask"}
	if outfile != "" {
		exec = outfile
		compileCmd = append(compileCmd, "-o", outfile)
	}

	err = execCmd(compileCmd...)
	if err != nil {
		return
	}

	err = os.Chdir(sourceDir)
	if err != nil {
		return
	}

	// return if exec file has been assigned
	if exec != "" {
		return
	}

	// find exec file if it's not there
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

	err = fmt.Errorf("can't locate build executable for task main %s", mainFile)
	return
}

func withTempDir(f func(string) error) (err error) {
	temp, err := ioutil.TempDir("", "go-task")
	if err != nil {
		return
	}
	defer func() {
		os.RemoveAll(temp)
	}()

	err = f(temp)
	return
}

func Run(sourceDir string, args []string) (err error) {
	parser := NewParser()
	taskSet, err := parser.Parse(sourceDir)
	if err != nil {
		return
	}

	err = withTempDir(func(work string) (err error) {
		compiler := compiler{sourceDir: sourceDir, workDir: work, TaskSet: taskSet}
		execFile, err := compiler.Compile("")
		if err != nil {
			return
		}

		runner := runner{execFile}
		err = runner.Run(args)
		return
	})

	return
}

func Compile(sourceDir string, outfile string) (err error) {
	parser := NewParser()
	taskSet, err := parser.Parse(sourceDir)
	if err != nil {
		return
	}

	err = withTempDir(func(work string) (err error) {
		compiler := compiler{sourceDir: sourceDir, workDir: work, TaskSet: taskSet}
		_, err = compiler.Compile(outfile)
		return
	})

	return
}
