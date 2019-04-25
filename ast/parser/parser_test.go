package ast

import (
	"fmt"
	"io/ioutil"
	"testing"

	"gotest.tools/assert"

	"github.com/Soontao/pdi-util/ast/lexer"

	"github.com/Soontao/pdi-util/ast/parser"
)

func TestParser(t *testing.T) {
	c, _ := ioutil.ReadFile("../basic_bo.bo")
	s := lexer.NewLexer(c)
	p := parser.NewParser()
	res, err := p.Parse(s)
	fmt.Println(res)
	assert.NilError(t, err)
}
