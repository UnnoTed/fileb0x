package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegexpComments(t *testing.T) {
	j := `{
    // wat 
  }`

	j = regexComments.ReplaceAllString(j, "")

	assert.Equal(t, `{
    
  }`, j)
}

func TestRegexpSafeVarName(t *testing.T) {
	v := `hi_i have a cat, my cat's name is cat, cat is a cat!`

	v = SafeVarName.ReplaceAllString(v, "")

	assert.Equal(t, `hiihaveacatmycatsnameiscatcatisacat`, v)
}
