package handler

import (
	"go/format"
)

func ParseAndGen(funcdata string) string {
	m := Parse(funcdata)
	genRst, _ := m.Gen()
	genRstFormat, _ := format.Source([]byte(genRst))
	return string(genRstFormat)
}
