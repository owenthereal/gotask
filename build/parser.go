package build

import (
	"fmt"
	"github.com/jingweno/gotask/tasking"
	"go/ast"
	"go/build"
	goparser "go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"
	"unicode/utf8"
)

func NewParser() *parser {
	return &parser{}
}

type parser struct{}

func (l *parser) Parse(dir string) (taskSet *tasking.TaskSet, err error) {
	dir, err = expandPath(dir)
	if err != nil {
		return
	}

	importPath, err := findImportPath(os.Getenv("GOPATH"), dir)
	if err != nil {
		return
	}

	p, e := build.Import(importPath, dir, 0)
	taskFiles := append(p.GoFiles, p.IgnoredGoFiles...)
	taskFiles = append(taskFiles, p.CgoFiles...)
	if e != nil {
		// task files may be ignored for build
		if _, ok := e.(*build.NoGoError); !ok || len(taskFiles) == 0 {
			err = e
			return
		}
	}

	tasks, err := loadTasks(dir, taskFiles)
	if err != nil {
		return
	}

	name := p.Name
	if name == "" {
		name = filepath.Base(p.Dir)
	}

	importPath = strings.Replace(p.ImportPath, "\\", "/", -1)

	taskSet = &tasking.TaskSet{Name: name, Dir: p.Dir, PkgObj: p.PkgObj, ImportPath: importPath, Tasks: tasks}

	return
}

func expandPath(path string) (expanded string, err error) {
	expanded, err = filepath.Abs(path)
	if err != nil {
		return
	}

	if !isFileExist(expanded) {
		err = fmt.Errorf("Path %s does not exist", expanded)
		return
	}

	return
}

func findImportPath(gp, dir string) (importPath string, err error) {
	var gopaths []string
	// GOPATHs are separated by ; on Windows
	if runtime.GOOS == "windows" {
		gopaths = strings.Split(gp, ";")
	} else {
		gopaths = strings.Split(gp, ":")
	}

	if len(gopaths) == 0 {
		err = fmt.Errorf("Environment variable GOPATH is not found")
		return
	}

	for _, gopath := range gopaths {
		gopath, e := expandPath(gopath)
		if e != nil {
			continue
		}

		srcPath := filepath.Join(gopath, "src")
		if !strings.HasPrefix(dir, srcPath) {
			continue
		}

		importPath, e = filepath.Rel(srcPath, dir)
		if e == nil && importPath != "" {
			break
		}
	}

	if importPath == "" {
		err = fmt.Errorf("Can't find import path in %s", dir)
	}

	return
}

func loadTasks(dir string, files []string) (tasks []tasking.Task, err error) {
	taskFiles := filterTaskFiles(files)
	for _, taskFile := range taskFiles {
		ts, e := parseTasks(filepath.Join(dir, taskFile))
		if e != nil {
			err = e
			return
		}

		tasks = append(tasks, ts...)
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

func parseTasks(filename string) (tasks []tasking.Task, err error) {
	taskFileSet := token.NewFileSet()
	f, err := goparser.ParseFile(taskFileSet, filename, nil, goparser.ParseComments)
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

		actionName := n.Name.String()
		if isTask(actionName, "Task") {
			p := &manPageParser{n.Doc.Text()}
			mp, e := p.Parse()
			if e != nil {
				continue
			}

			if mp.Name == "" {
				mp.Name = convertActionNameToTaskName(actionName)
			}

			t := tasking.Task{Name: mp.Name, ActionName: actionName, Usage: mp.Usage, Description: mp.Description, Flags: mp.Flags}
			tasks = append(tasks, t)
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

func convertActionNameToTaskName(s string) string {
	n := strings.TrimPrefix(s, "Task")
	return dasherize(n)
}
