package code

import (
	"bytes"
	"fmt"
	"github.com/JSYoo5B/chain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChainCode_GenerateChainCode(t *testing.T) {
	buf := bytes.Buffer{}

	code := ChainCode{
		PackageName: "branch",
		Imports:     []string{`_ "context"`, `_ "fmt"`},
		Workflows: []WorkflowDef{
			WorkflowDef{
				ConstructName:   "basicCollatzFunction",
				ConstructParams: "",
				NodeConstructors: []string{
					`branch, even := checkNext(), half()`,
					`odd1, odd2 := triple(), inc()`,
				},
				WorkflowName: "ShortcutCollatz",
				WorkflowType: "int",
				Nodes: []WorkflowNode{
					WorkflowNode{VarName: "branch", ConstructExpr: "checkNext()"},
					WorkflowNode{VarName: "even", ConstructExpr: "half()"},
					WorkflowNode{VarName: "odd1", ConstructExpr: "triple()"},
					WorkflowNode{VarName: "odd2", ConstructExpr: "inc()"},
				},
				Edges: []WorkflowEdge{
					WorkflowEdge{
						BaseNode: "branch",
						WorkType: "int",
						Plan:     map[string]string{"even": "even", "odd": "odd1"},
					},
					WorkflowEdge{
						BaseNode: "even",
						WorkType: "int",
						Plan:     nil,
					},
					WorkflowEdge{
						BaseNode: "odd1",
						WorkType: "int",
						Plan:     map[string]string{chain.Success: "odd2"},
					},
					WorkflowEdge{
						BaseNode: "odd2",
						WorkType: "int",
						Plan:     nil,
					},
				},
			},
		},
	}

	err := code.GenerateChainCode(&buf)
	assert.NoError(t, err)

	fmt.Println(buf.String())
}
