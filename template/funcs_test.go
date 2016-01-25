package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExported(t *testing.T) {
	SetUnexported(true)
	e := exported("HELLO")
	assert.Equal(t, `hello`, e)

	SetUnexported(false)
	e = exported("HELLO")
	assert.Equal(t, `HELLO`, e)
}

func TestExportedTitle(t *testing.T) {
	SetUnexported(true)
	e := exportedTitle("HELLO")
	assert.Equal(t, `hELLO`, e)

	SetUnexported(false)
	e = exportedTitle("hello")
	assert.Equal(t, `Hello`, e)
}

func TestVarName(t *testing.T) {
	ns := `s@af{{e}} Var!@)*iable& Na!*(@$!@)me`
	s := buildSafeVarName(ns)
	assert.Equal(t, `safeVariableName`, s)
}
