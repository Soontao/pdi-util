package parser

import (
	"io/ioutil"
	"testing"

	"gotest.tools/assert"

	"github.com/Soontao/pdi-util/ast/lexer"
)

func TestParserBO(t *testing.T) {
	c, _ := ioutil.ReadFile("../test_data/sample_bo.bo")
	s := lexer.NewLexer(c)
	p := NewParser()
	_, err := p.Parse(s)
	assert.NilError(t, err)
}

func TestParserABSL(t *testing.T) {
	c, _ := ioutil.ReadFile("../test_data/sample_absl.absl")
	s := lexer.NewLexer(c)
	p := NewParser()
	_, err := p.Parse(s)
	assert.NilError(t, err)
}
