package bindvar_test

import (
	"testing"

	"github.com/noypi/bindvar"
	assertpkg "github.com/stretchr/testify/assert"
)

func TestBind(t *testing.T) {
	assert := assertpkg.New(t)

	n := bindvar.NewInt()
	s := bindvar.NewString()

	var s0 string

	bindvar.BindStrFmt(&s0, "n=%d, s=%s", n, s)

	assert.Equal("n=0, s=", s0)

	n.Set(100)
	assert.Equal("n=100, s=", s0)

	s.Set("giant")
	assert.Equal("n=100, s=giant", s0)
}

func TestBindStr(t *testing.T) {
	assert := assertpkg.New(t)

	s := bindvar.NewString()
	var s0 string

	bindvar.BindStr(&s0, "s0 is=", s)

	assert.Equal("s0 is=", s0)
	s.Set("giant")
	assert.Equal("s0 is=giant", s0)

	s2 := bindvar.NewString()
	s3 := bindvar.NewString()
	s.Bind("setting s2=", s2)
	s.BindFmt("s3 %s", s3)

	s2.Set("s2 is")
	assert.Equal("s3 ", s.String())
	assert.Equal("s0 is=s3 ", s0)

	s3.Set("s3 is")
	assert.Equal("s3 s3 is", s.String())
	assert.Equal("s0 is=s3 s3 is", s0)

}
