package template

import (
	"regexp"
	"strings"
	"text/template"

	"github.com/UnnoTed/fileb0x/config"
)

var r = regexp.MustCompile(`[^a-zA-Z0-9]`)

var funcsTemplate = template.FuncMap{
	"exported":         exported,
	"exportedTitle":    exportedTitle,
	"buildSafeVarName": buildSafeVarName,
}

var unexported bool

// SetUnexported variables, functions and types
func SetUnexported(e bool) {
	unexported = e
}

func exported(field string) string {
	if !unexported {
		return strings.ToUpper(field)
	}

	return strings.ToLower(field)
}

func exportedTitle(field string) string {
	if !unexported {
		return strings.Title(field)
	}

	return strings.ToLower(field[0:1]) + field[1:]
}

func buildSafeVarName(path string) string {
	return config.SafeVarName.ReplaceAllString(path, "")
}
