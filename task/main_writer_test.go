package task

import (
	"bytes"
	"github.com/bmizerany/assert"
	"strings"
	"testing"
)

func TestMainWriter_Write(t *testing.T) {
	var out bytes.Buffer
	b := mainWriter{
		&taskFuncs{
			ImportPath: "github.com/jingweno/gotask/examples",
			Funcs:      []taskFunc{taskFunc{Name: "TaskHelloWorld", Usage: "Say Hello world", Description: "Print out Hello World"}},
		},
	}
	b.Write(&out)

	assert.Tf(t, strings.Contains(out.String(), `_task "github.com/jingweno/gotask/examples"`), "%v", out.String())
	assert.Tf(t, strings.Contains(out.String(), `{"HelloWorld", _task.TaskHelloWorld}`), "%v", out.String())
}
