package utils

import (
	"path/filepath"
	"strings"
)

// FixPath converts \ and \\ to /
func FixPath(path string) string {
	a := filepath.Clean(path)
	b := strings.Replace(a, `\`, "/", -1)
	c := strings.Replace(b, `\\`, "/", -1)
	return c
}

// FixName converts [/ to _](1), [  to -](2) and [, to __](3)
func FixName(path string) string {
	a := FixPath(path)
	b := strings.Replace(a, "/", "_", -1)     // / to _
	c := strings.Replace(b, " ", "-", -1)     // {space} to -
	return strings.Replace(c, ",", "___", -1) // , to __
}
