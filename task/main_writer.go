package task

import (
	"io"
	"text/template"
)

type mainWriter struct {
	TaskSet *TaskSet
}

func (w *mainWriter) Write(wr io.Writer) (err error) {
	err = taskmainTmpl.Execute(wr, w.TaskSet)
	return
}

var taskmainTmpl = template.Must(template.New("main").Parse(`
package main

import (
  "os"
  "github.com/jingweno/gotask/task"
{{if .HasTasks}}
  _task "{{.ImportPath}}"
{{end}}
)

var tasks = []task.Task{
{{range .Tasks}}
  {Name: "{{.Name}}", Usage: {{.Usage | printf "%q"}}, Description: {{.Description | printf "%q"}}, Action: _task.{{.ActionName}}},
{{end}}
}

func main() {
  result := task.RunTasks(tasks)
  os.Exit(result.ExitCode)
}
`))
