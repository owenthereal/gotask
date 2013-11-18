package build

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestManPageParse_Parse(t *testing.T) {
	doc := `NAME
    say-hello - Say hello to current user

DESCRIPTION
    Print out hello to current user
    one more line

`
	p := &manPageParser{doc}
	mp, err := p.Parse()

	assert.Equal(t, nil, err)
	assert.Equal(t, "say-hello", mp.Name)
	assert.Equal(t, "Say hello to current user", mp.Usage)
	assert.Equal(t, "Print out hello to current user\n   one more line", mp.Description)

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
}
