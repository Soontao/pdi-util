// Code generated by gocc; DO NOT EDIT.

package parser

import "github.com/Soontao/pdi-util/ast"

type (
	//TODO: change type and variable names to be consistent with other tables
	ProdTab      [numProductions]ProdTabEntry
	ProdTabEntry struct {
		String     string
		Id         string
		NTType     int
		Index      int
		NumSymbols int
		ReduceFunc func([]Attrib) (Attrib, error)
	}
	Attrib interface {
	}
)

var productionsTable = ProdTab{
	ProdTabEntry{
		String: `S' : StatementList	<<  >>`,
		Id:         "S'",
		NTType:     0,
		Index:      0,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `StatementList : Statement terminator RepeatTerminator StatementList	<<  >>`,
		Id:         "StatementList",
		NTType:     1,
		Index:      1,
		NumSymbols: 4,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `StatementList : Statement RepeatTerminator	<<  >>`,
		Id:         "StatementList",
		NTType:     1,
		Index:      2,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `RepeatTerminator : terminator RepeatTerminator	<<  >>`,
		Id:         "RepeatTerminator",
		NTType:     2,
		Index:      3,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `RepeatTerminator : empty	<<  >>`,
		Id:         "RepeatTerminator",
		NTType:     2,
		Index:      4,
		NumSymbols: 0,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return nil, nil
		},
	},
	ProdTabEntry{
		String: `Statement : ImportAsDeclaration	<<  >>`,
		Id:         "Statement",
		NTType:     3,
		Index:      5,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Statement : ImportNormalDeclaration	<<  >>`,
		Id:         "Statement",
		NTType:     3,
		Index:      6,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `ImportAsDeclaration : "import" Namespace "as" Identifier	<< ast.NewImportAsDeclaration(X[1], X[3]) >>`,
		Id:         "ImportAsDeclaration",
		NTType:     4,
		Index:      7,
		NumSymbols: 4,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewImportAsDeclaration(X[1], X[3])
		},
	},
	ProdTabEntry{
		String: `ImportNormalDeclaration : "import" Namespace	<<  >>`,
		Id:         "ImportNormalDeclaration",
		NTType:     5,
		Index:      8,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Namespace : Identifier	<< ast.NewNamespace(X[0], nil) >>`,
		Id:         "Namespace",
		NTType:     6,
		Index:      9,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewNamespace(X[0], nil)
		},
	},
	ProdTabEntry{
		String: `Namespace : Identifier "." Namespace	<< ast.NewNamespace(X[0], X[2]) >>`,
		Id:         "Namespace",
		NTType:     6,
		Index:      10,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewNamespace(X[0], X[2])
		},
	},
	ProdTabEntry{
		String: `Identifier : identifier	<< ast.NewIdentifier(X[0]) >>`,
		Id:         "Identifier",
		NTType:     7,
		Index:      11,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewIdentifier(X[0])
		},
	},
}
