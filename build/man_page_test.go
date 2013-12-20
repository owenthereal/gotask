package build

import (
	"github.com/bmizerany/assert"
	"github.com/jingweno/gotask/task"
	"testing"
)

func TestManPageParse_Parse(t *testing.T) {
	doc := `NAME
    say-hello - Say hello to current user

DESCRIPTION
    Print out hello to current user
    one more line

OPTIONS
    -v, --verbose
        Run in verbose mode
    -n, --name=NAME
        Say hello to the given NAME
    -g, --greeting=TYPE
        Say hello with a custom TYPE of greeting

`
	p := &manPageParser{doc}
	mp, err := p.Parse()

	assert.Equal(t, nil, err)
	assert.Equal(t, "say-hello", mp.Name)
	assert.Equal(t, "Say hello to current user", mp.Usage)
	assert.Equal(t, "Print out hello to current user\n   one more line", mp.Description)
	assert.Equal(t, 3, len(mp.Flags))

	firstFlag, ok := mp.Flags[0].(task.BoolFlag)
	assert.Tf(t, ok, "Can't convert flag to task.BoolFlag")
	assert.Equal(t, "v, verbose", firstFlag.Name)
	assert.Equal(t, "Run in verbose mode", firstFlag.Usage)

	stringFlag, ok := mp.Flags[1].(task.StringFlag)
	assert.Tf(t, ok, "Can't convert flag to task.StringFlag")
	assert.Equal(t, "n, name", stringFlag.Name)
	assert.Equal(t, "Say hello to the given NAME", stringFlag.Usage)

	stringFlag, ok = mp.Flags[2].(task.StringFlag)
	assert.Tf(t, ok, "Can't convert flag to task.StringFlag")
	assert.Equal(t, "g, greeting", stringFlag.Name)
	assert.Equal(t, "Say hello with a custom TYPE of greeting", stringFlag.Usage)

	doc = `Name
    say-hello - Say hello to current user

Description
    Print out hello to current user
`
	p = &manPageParser{doc}
	mp, err = p.Parse()

	assert.Equal(t, nil, err)
	assert.Equal(t, "", mp.Name)
	assert.Equal(t, "", mp.Usage)
	assert.Equal(t, "", mp.Description)
	assert.Equal(t, 0, len(mp.Flags))
}
