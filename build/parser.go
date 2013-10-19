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
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

func NewParser() *parser {
	return &parser{}
}

type parser struct{}

func (l *parser) Parse(dir string) (taskSet *tasking.TaskSet, err error) {
	dir, err = expandDir(dir)
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

	tasks, err := loadTasks(dir, taskFiles)
	if err != nil {
		return
	}

	name := p.Name
	if name == "" {
		name = filepath.Base(p.Dir)
	}

	taskSet = &tasking.TaskSet{Name: name, Dir: p.Dir, ImportPath: p.ImportPath, Tasks: tasks}

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
			usage, desc, e := parseUsageAndDesc(n.Doc.Text())
			if e != nil {
				continue
			}

			name := convertActionNameToTaskName(actionName)
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

func parseUsageAndDesc(doc string) (usage, desc string, err error) {
	reader := bufio.NewReader(bytes.NewReader([]byte(doc)))
	r := regexp.MustCompile("\\S")
	var usageParts, descParts []string

	line, err := readLine(reader)
	for err == nil {
		if len(descParts) == 0 && r.MatchString(line) {
			usageParts = append(usageParts, line)
		} else {
			descParts = append(descParts, line)
		}

		line, err = readLine(reader)
	}

	if err == io.EOF {
		err = nil
	}

	usage = strings.Join(usageParts, " ")
	usage = strings.TrimSpace(usage)

	desc = strings.Join(descParts, "\n")
	desc = strings.TrimSpace(desc)

	return
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
