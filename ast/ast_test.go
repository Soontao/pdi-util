package ast

import (
	"io/ioutil"
	"testing"

	"gotest.tools/assert"
)

func TestParserBO(t *testing.T) {
	c, _ := ioutil.ReadFile("./test_data/sample_bo.bo")
	n, err := ParseAST(c)
	assert.NilError(t, err)
	assert.Assert(t, n != nil)
	bo := n.GetNode("BODefinition")
	assert.Assert(t, bo != nil)
	boName := bo.GetNode("BOName")
	assert.Assert(t, boName != nil)
	id := boName.GetNode("Identifier")
	assert.Assert(t, id != nil)
	assert.Assert(t, id.GetString("ID") == "BO_FDNTest")
	elements := bo.GetNodeList("Elements")
	assert.Assert(t, elements != nil)
}

func TestParserABSL(t *testing.T) {
	c, _ := ioutil.ReadFile("./test_data/sample_absl.absl")
	_, err := ParseAST(c)
	assert.NilError(t, err)
}
