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
			Tasks: []tasking.Task{
				{
					Name:        "HelloWorld",
					ActionName:  "TaskHelloWorld",
					Usage:       "Say Hello world",
					Description: "Print out Hello World",
					Flags: []tasking.Flag{
						tasking.BoolFlag{Name: "-v --verbose", Usage: "Run in verbose mode"},
					},
				},
			},
		},
	}
	b.Write(&out)

	assert.Tf(t, strings.Contains(out.String(), `_task "github.com/jingweno/gotask/examples"`), "%v", out.String())
	assert.Tf(t, strings.Contains(out.String(), `Name: "HelloWorld"`), "%v", out.String())
	assert.Tf(t, strings.Contains(out.String(), `Usage: "Say Hello world"`), "%v", out.String())
	assert.Tf(t, strings.Contains(out.String(), `Description: "Print out Hello World`), "%v", out.String())
	assert.Tf(t, strings.Contains(out.String(), `tasking.BoolFlag{Name: "-v --verbose", Usage: "Run in verbose mode"}`), "%v", out.String())
}
