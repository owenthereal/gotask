package task

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"
)

var taskFileSet = token.NewFileSet()

type taskFuncs struct {
	ImportPath string
	Funcs      []taskFunc
}

func (t *taskFuncs) HasTasks() bool {
	return len(t.Funcs) > 0
}

type taskFunc struct {
	Name string
}

type taskLoader struct {
	Dir string
}

func (l *taskLoader) Load() (funcs *taskFuncs, err error) {
	dir, err := expandDir(l.Dir)
	if err != nil {
		return
	}

	p, e := build.ImportDir(dir, 0)
	taskFiles := append(p.GoFiles, p.IgnoredGoFiles...)
	taskFiles = append(taskFiles, p.CgoFiles...)
	if e != nil {
		// task files may be ignored for build
		if _, ok := e.(*build.NoGoError); !ok || len(taskFiles) == 0 {
			err = e
			return
		}
	}

	fs, err := loadTaskFuncs(dir, taskFiles)
	if err != nil {
		return
	}

	funcs = &taskFuncs{ImportPath: p.ImportPath, Funcs: fs}

	return
}

func expandDir(dir string) (expanded string, err error) {
	expanded, err = filepath.Abs(dir)
	if err != nil {
		return
	}

	if !isFileExist(dir) {
		err = fmt.Errorf("Directory %s does not exist", dir)
		return
	}

	return
}

func loadTaskFuncs(dir string, files []string) (taskFuncs []taskFunc, err error) {
	taskFiles := filterTaskFiles(files)
	for _, taskFile := range taskFiles {
		names, e := parseTaskNames(filepath.Join(dir, taskFile))
		if e != nil {
			err = e
			return
		}

		for _, name := range names {
			taskFuncs = append(taskFuncs, taskFunc{Name: name})
		}
	}

	return
}

func filterTaskFiles(files []string) (taskFiles []string) {
	for _, f := range files {
		if isTaskFile(f, "_task.go") {
			taskFiles = append(taskFiles, f)
		}
	}

	return
}

func parseTaskNames(filename string) (names []string, err error) {
	f, err := parser.ParseFile(taskFileSet, filename, nil, 0)
	if err != nil {
		return
	}

	for _, d := range f.Decls {
		n, ok := d.(*ast.FuncDecl)
		if !ok {
			continue
		}

		if n.Recv != nil {
			continue
		}

		name := n.Name.String()
		if isTask(name, "Task") {
			names = append(names, name)
		}
	}

	return
}

func isTaskFile(name, suffix string) bool {
	if strings.HasSuffix(name, suffix) {
		return true
	}

	return false
}

func isTask(name, prefix string) bool {
	if !strings.HasPrefix(name, prefix) {
		return false
	}
	if len(name) == len(prefix) { // "Task" is ok
		return true
	}

	rune, _ := utf8.DecodeRuneInString(name[len(prefix):])
	return !unicode.IsLower(rune)
}
