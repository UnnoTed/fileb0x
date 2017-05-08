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
	var s string
	s = buildSafeVarName(`a/safe/variable.name`)
	assert.Equal(t, `ASafeVariableName`, s)

	s = buildSafeVarName(`a/safe/variable.html`)
	assert.Equal(t, `ASafeVariableHTML`, s)

	s = buildSafeVarName(`a/safe/variable.json`)
	assert.Equal(t, `ASafeVariableJSON`, s)

	s = buildSafeVarName(`a/safe/variable.url`)
	assert.Equal(t, `ASafeVariableURL`, s)

	s = buildSafeVarName(`a/sql/variable.name`)
	assert.Equal(t, `ASQLVariableName`, s)

	s = buildSafeVarName(`a/sql/_variable.name`)
	assert.Equal(t, `ASQLVariableName2`, s)
}
