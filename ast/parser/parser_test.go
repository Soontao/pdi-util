package parser

import (
	"io/ioutil"
	"testing"

	"gotest.tools/assert"

	"github.com/Soontao/pdi-util/ast/lexer"
)

func TestParser(t *testing.T) {
	c, _ := ioutil.ReadFile("../basic_bo.bo")
	s := lexer.NewLexer(c)
	p := NewParser()
	_, err := p.Parse(s)
	assert.NilError(t, err)
}
