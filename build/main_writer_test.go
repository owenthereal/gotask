package build

import (
	"bytes"
	"github.com/bmizerany/assert"
	"github.com/jingweno/gotask/tasking"
	"strings"
	"testing"
)

func TestMainWriter_Write(t *testing.T) {
	var out bytes.Buffer
	b := mainWriter{
		&tasking.TaskSet{
			ImportPath: "github.com/jingweno/gotask/examples",
			Tasks:      []tasking.Task{{Name: "HelloWorld", ActionName: "TaskHelloWorld", Usage: "Say Hello world", Description: "Print out Hello World"}},
		},
	}
	b.Write(&out)

	assert.Tf(t, strings.Contains(out.String(), `_task "github.com/jingweno/gotask/examples"`), "%v", out.String())
	assert.Tf(t, strings.Contains(out.String(), `{Name: "HelloWorld", Usage: "Say Hello world", Description: "Print out Hello World", Action: _task.TaskHelloWorld}`), "%v", out.String())
}
