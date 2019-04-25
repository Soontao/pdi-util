package ast

import "github.com/Soontao/pdi-util/ast/token"

// NewImportAsDeclaration
func NewImportAsDeclaration(n, i interface{}) (interface{}, error) {
	return &ImportAsDeclaration{n.(*Namespace), i.(*Identifier)}, nil
}

// NewNamespace
func NewNamespace(id, sub interface{}) (interface{}, error) {
	if sub != nil {
		return &Namespace{id.(*Identifier), sub.(*Namespace)}, nil
	}
	return &Namespace{id.(*Identifier), nil}, nil
}

// NewIdentifier
func NewIdentifier(i interface{}) (interface{}, error) {
	return &Identifier{string((i.(*token.Token)).Lit)}, nil
}
