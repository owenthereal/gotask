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
		&TaskSet{
			ImportPath: "github.com/jingweno/gotask/examples",
			Tasks:      []Task{{Name: "HelloWorld", ActionName: "TaskHelloWorld", Usage: "Say Hello world", Description: "Print out Hello World"}},
		},
	}
	b.Write(&out)

	assert.Tf(t, strings.Contains(out.String(), `_task "github.com/jingweno/gotask/examples"`), "%v", out.String())
	assert.Tf(t, strings.Contains(out.String(), `{"HelloWorld", "Say Hello world", "Print out Hello World", _task.TaskHelloWorld}`), "%v", out.String())
}
