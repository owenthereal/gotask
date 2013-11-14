package build

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/jingweno/gotask/tasking"
	"go/ast"
	"go/build"
	goparser "go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"regexp"
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
			m, e := parseManPage(n.Doc.Text())
			if e != nil {
				continue
			}

			name := m["NAME"]
			usage := m["USAGE"]
			desc := m["DESCRIPTION"]
			if name == "" {
				name = convertActionNameToTaskName(actionName)
			}

			t := tasking.Task{Name: name, ActionName: actionName, Usage: usage, Description: desc}
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

func parseManPage(doc string) (result map[string]string, err error) {
	result = make(map[string]string)
	headingRegexp := regexp.MustCompile(`^([A-Z]+)$`)
	reader := bufio.NewReader(bytes.NewReader([]byte(doc)))

	var (
		line    string
		heading string
		content []string
	)
	for err == nil {
		line, err = readLine(reader)

		if headingRegexp.MatchString(line) {
			if heading != line {
				if heading != "" {
					result[heading] = concatHeadingContent(content)
				}

				heading = line
				content = []string{}
			}
		} else {
			if line != "" {
				line = strings.TrimSpace(line)
			}
			content = append(content, line)
		}
	}
	// the last one
	if heading != "" {
		result[heading] = concatHeadingContent(content)
	}

	if err == io.EOF {
		err = nil
	}

	// set NAME and USAGE
	if name, ok := result["NAME"]; ok {
		s := strings.SplitN(name, " - ", 2)
		if len(s) == 1 {
			result["NAME"] = ""
			result["USAGE"] = strings.TrimSpace(s[0])
		} else {
			result["NAME"] = strings.TrimSpace(s[0])
			result["USAGE"] = strings.TrimSpace(s[1])
		}
	}

	return
}

func concatHeadingContent(content []string) string {
	return strings.TrimSpace(strings.Join(content, "\n   "))
}

func readLine(r *bufio.Reader) (string, error) {
	var (
		isPrefix = true
		err      error
		line, ln []byte
	)

	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}

	return string(ln), err
}
