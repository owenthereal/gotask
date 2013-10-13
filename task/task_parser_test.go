package task

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestTaskParser_Load(t *testing.T) {
	l := taskParser{"../examples"}
	funcs, err := l.Parse()

	assert.Equal(t, nil, err)
	assert.Equal(t, "github.com/jingweno/gotask/examples", funcs.ImportPath)
	assert.Equal(t, 1, len(funcs.Funcs))
}

func TestTaskParser_filterTaskFiles(t *testing.T) {
	files := []string{"file.go", "file_task.go", "task.go"}
	taskFiles := filterTaskFiles(files)

	assert.Equal(t, 1, len(taskFiles))
	assert.Equal(t, "file_task.go", taskFiles[0])
}

func TestTaskParser_parseTaskFuncs(t *testing.T) {
	funcs, _ := parseTaskFuncs("../examples/example_task.go")

	assert.Equal(t, 1, len(funcs))
	assert.Equal(t, "TaskHelloWorld", funcs[0].Name)
	assert.Equal(t, "Say hello world", funcs[0].Usage)
	assert.Equal(t, "Print out hello world", funcs[0].Description)
}

func TestParseUsageAndDesc(t *testing.T) {
	doc := `Usage

Desc

Desc2
`
	usage, desc, err := parseUsageAndDesc(doc)

	assert.Equal(t, nil, err)
	assert.Equal(t, "Usage", usage)
	assert.Equal(t, "Desc\n\nDesc2", desc)
}
