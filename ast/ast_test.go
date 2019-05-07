package ast

import (
	"io/ioutil"
	"testing"

	"gotest.tools/assert"
)

func TestParserBO(t *testing.T) {
	c, _ := ioutil.ReadFile("./test_data/sample_bo.bo")
	_, err := ParseAST(c)
	assert.NilError(t, err)
}

func TestParserABSL(t *testing.T) {
	c, _ := ioutil.ReadFile("./test_data/sample_absl.absl")
	_, err := ParseAST(c)
	assert.NilError(t, err)
}
