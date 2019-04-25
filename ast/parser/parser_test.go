package parser

import (
	"fmt"
	"io/ioutil"
	"testing"

	"gotest.tools/assert"

	"github.com/Soontao/pdi-util/ast/lexer"
)

func TestParser(t *testing.T) {
	c, _ := ioutil.ReadFile("../basic_bo.bo")
	s := lexer.NewLexer(c)
	p := NewParser()
	res, err := p.Parse(s)
	fmt.Println(res)
	assert.NilError(t, err)
}
