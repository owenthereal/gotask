package build

import (
	"github.com/jingweno/gotask/tasking"
	"io"
	"text/template"
)

type mainWriter struct {
	TaskSet *tasking.TaskSet
}

func (w *mainWriter) Write(wr io.Writer) (err error) {
	err = taskmainTmpl.Execute(wr, w.TaskSet)
	return
}

var taskmainTmpl = template.Must(template.New("main").Parse(`
package main

import (
  "os"
  "github.com/jingweno/gotask/tasking"
{{if .HasTasks}}
  _task "{{.ImportPath}}"
{{end}}
)

var tasks = []tasking.Task{
{{range .Tasks}}
  {
    Name: {{.Name | printf "%q" }},
    Usage: {{.Usage | printf "%q"}},
    Description: {{.Description | printf "%q"}},
    Action: _task.{{.ActionName}},
    Flags: []tasking.Flag{
      {{range .Flags}}
        {{.DefType "tasking"}},
      {{end}}
    },
  },
{{end}}
}

var taskSet = tasking.TaskSet{
  Name: {{.Name | printf "%q" }},
  Dir: {{.Dir | printf "%q" }},
  PkgObj: {{.PkgObj | printf "%q" }},
  ImportPath: {{.ImportPath | printf "%q" }},
  Tasks: tasks,
}

func main() {
  tasking.Run(&taskSet, os.Args)
}
`))
