package ast

type ChainCode struct {
	Package   Package
	Imports   []Import
	Workflows []WorkflowStatement
}

type Package struct {
	Name string
	CodeLocation
}

type Import struct {
	Alias string
	Path  string
	CodeLocation
}
