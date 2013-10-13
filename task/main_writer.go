package task

import (
	"io"
	"text/template"
)

type mainWriter struct {
	Funcs *taskFuncs
}

func (b *mainWriter) Write(wr io.Writer) (err error) {
	err = taskmainTmpl.Execute(wr, b.Funcs)
	return
}

var taskmainTmpl = template.Must(template.New("main").Parse(`
package main

import (
	"github.com/jingweno/gotask/task"
{{if .HasTasks}}
	_task "{{.ImportPath}}"
{{end}}
)

var taskFuncs = []func(*task.T){
{{range .Funcs}}
	_task.{{.Name}},
{{end}}
}

func main() {
	task.RunTasks(taskFuncs)
}
`))
