package tasking

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestBoolFlag_DefType(t *testing.T) {
	f := BoolFlag{Name: "name", Usage: "usage"}
	assert.Equal(t, `tasking.BoolFlag{Name: "name", Usage: "usage"}`, f.DefType("tasking"))
}
