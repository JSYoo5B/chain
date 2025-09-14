package ast

type WorkflowStatement struct {
	WorkflowDeclaration
	WorkflowDefinition

	CodeLocation
}

type WorkflowDeclaration struct {
	ConstructorName   CodeLocation
	ConstructorParams CodeLocation

	WorkflowName CodeLocation
	WorkflowType CodeLocation

	CodeLocation
}

type WorkflowDefinition struct {
	PrerequisiteBlock
	NodesBlock

	Branches  []BranchStatement
	Successes []DirectionStatement
	Failures  []DirectionStatement
	Aborts    []DirectionStatement

	CodeLocation
}

type PrerequisiteBlock struct {
	Code string
	CodeLocation
}

type NodesBlock struct {
	Nodes []WorkflowNode
	CodeLocation
}

type WorkflowNode struct {
	Name string
	CodeLocation
}

type DirectionStatement struct {
	FromNode string
	ToNode   string
	CodeLocation
}

type BranchStatement struct {
	FromNode  string
	Condition string
	ToNode    string
	CodeLocation
}
