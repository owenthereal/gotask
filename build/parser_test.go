package build

import (
	"fmt"
	"github.com/bmizerany/assert"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestParser_findImportPath(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	if runtime.GOOS == "windows" {
		gopath = fmt.Sprintf("/etc;%s", gopath)
	} else {
		gopath = fmt.Sprintf("/etc:%s", gopath)
	}
	dir, _ := expandPath("../examples")
	importPath, err := findImportPath(gopath, dir)

	assert.Equal(t, nil, err)
	assert.Equal(t, filepath.Join("github.com", "jingweno", "gotask", "examples"), importPath)
}

func TestParser_Load(t *testing.T) {
	p := NewParser()
	ts, err := p.Parse("../examples")

	assert.Equal(t, nil, err)
	assert.Tf(t, strings.HasSuffix(ts.PkgObj, filepath.Join("github.com", "jingweno", "gotask", "examples.a")), "%s", ts.PkgObj)
	assert.Equal(t, "github.com/jingweno/gotask/examples", ts.ImportPath)
	assert.Equal(t, 2, len(ts.Tasks))
}

func TestTaskParser_filterTaskFiles(t *testing.T) {
	files := []string{"file.go", "file_task.go", "task.go"}
	taskFiles := filterTaskFiles(files)

	assert.Equal(t, 1, len(taskFiles))
	assert.Equal(t, "file_task.go", taskFiles[0])
}

func TestParser_parseTasks(t *testing.T) {
	tasks, _ := parseTasks("../examples/say_hello_task.go")

	assert.Equal(t, 1, len(tasks))
	assert.Equal(t, "TaskSayHello", tasks[0].ActionName)
	assert.Equal(t, "say-hello", tasks[0].Name)
	assert.Equal(t, "Say hello to current user", tasks[0].Usage)
	assert.Equal(t, "Print out hello to current user", tasks[0].Description)
}
