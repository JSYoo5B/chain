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
	Prerequisite PrerequisiteBlock
	Nodes        NodeBlock

	Branches  []BranchStatement
	Successes []DirectionStatement
	Errors    []DirectionStatement
	Aborts    []DirectionStatement

	CodeLocation
}

type PrerequisiteBlock struct {
	Code string
	CodeLocation
}

type NodeBlock struct {
	Nodes []WorkflowNode
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
