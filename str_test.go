package bindvar

import (
	"testing"

	assertpkg "github.com/stretchr/testify/assert"
)

func TestBindStr(t *testing.T) {
	assert := assertpkg.New(t)

	s := NewString()

	var s0 string

	BindStr(&s0, s)

	s.Set("updating value")
	assert.Equal("updating value", s0)

	s.Set("some new value")
	assert.Equal("some new value", s0)

	Unbind(&s0)
	assert.Equal(0, len(defaultBindVar.fns))
	assert.Equal(0, len(defaultBindVar.m[s.Ptr()]))
}
