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
			Funcs:      []taskFunc{taskFunc{"TaskHelloWorld"}},
		},
	}
	b.Write(&out)

	assert.T(t, strings.Contains(out.String(), `_task "github.com/jingweno/gotask/examples"`))
	assert.T(t, strings.Contains(out.String(), `{"TaskHelloWorld", _task.TaskHelloWorld}`))
}
