package task

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestTaskParser_Load(t *testing.T) {
	l := taskParser{"../examples"}
	ts, err := l.Parse()

	assert.Equal(t, nil, err)
	assert.Equal(t, "github.com/jingweno/gotask/examples", ts.ImportPath)
	assert.Equal(t, 2, len(ts.Tasks))
}

func TestTaskParser_filterTaskFiles(t *testing.T) {
	files := []string{"file.go", "file_task.go", "task.go"}
	taskFiles := filterTaskFiles(files)

	assert.Equal(t, 1, len(taskFiles))
	assert.Equal(t, "file_task.go", taskFiles[0])
}

func TestTaskParser_parseTasks(t *testing.T) {
	tasks, _ := parseTasks("../examples/say_hello_task.go")

	assert.Equal(t, 1, len(tasks))
	assert.Equal(t, "say-hello", tasks[0].Name)
	assert.Equal(t, "TaskSayHello", tasks[0].ActionName)
	assert.Equal(t, "Say hello to current user", tasks[0].Usage)
	assert.Equal(t, "Print out hello to current user", tasks[0].Description)
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
