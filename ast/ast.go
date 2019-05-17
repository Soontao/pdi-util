package ast

import (
	"github.com/Soontao/pdi-util/ast/lexer"
	"github.com/Soontao/pdi-util/ast/parser"
	"github.com/Soontao/pdi-util/ast/types"
)

// ParseAST node
func ParseAST(source []byte) (*types.GrammerNode, error) {
	s := lexer.NewLexer(source)
	p := parser.NewParser()
	n, err := p.Parse(s)
	if err != nil {
		return nil, err
	}
	return n.(*types.GrammerNode), err
}
