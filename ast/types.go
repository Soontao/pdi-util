package ast

// ImportAsDeclaration type
type ImportAsDeclaration struct {
	Namespace *Namespace
	Alias     *Identifier
}

// Namespace type
type Namespace struct {
	identifier *Identifier
	sub        *Namespace
}

// Identifier type
type Identifier struct {
	value string
}
