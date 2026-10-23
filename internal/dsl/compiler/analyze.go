// Package compiler validates DSL workflows and generates Go workflow constructors.
package compiler

import (
	"fmt"
	"go/scanner"

	"github.com/JSYoo5B/chain/internal/dsl/frontend"
)

// Program contains validated execution plans in declaration order.
type Program struct {
	File      *frontend.File
	Workflows []WorkflowPlan
}

type WorkflowPlan struct {
	Syntax *frontend.Workflow
	Nodes  []NodePlan
}

type NodePlan struct {
	Name       frontend.Identifier
	Directions map[string]frontend.Identifier
}

// Analyze checks declared names, direction uniqueness, directed cycles, and weak
// connectivity, matching Workflow.ValidateGraph. The first node remains the
// entry point. Go types and dynamic Action/BranchAction contracts are checked
// by the generated package and runtime, respectively.
func Analyze(file *frontend.File) (*Program, error) {
	program := &Program{File: file}
	var diagnostics scanner.ErrorList
	for _, workflow := range file.Workflows {
		plan, errs := analyzeWorkflow(file, workflow)
		diagnostics = append(diagnostics, errs...)
		program.Workflows = append(program.Workflows, plan)
	}
	if err := diagnostics.Err(); err != nil {
		diagnostics.Sort()
		return nil, diagnostics.Err()
	}
	return program, nil
}

func analyzeWorkflow(file *frontend.File, w *frontend.Workflow) (WorkflowPlan, scanner.ErrorList) {
	plan := WorkflowPlan{Syntax: w}
	var diagnostics scanner.ErrorList
	report := func(name frontend.Identifier, message string) {
		diagnostics.Add(file.Source.Position(name.Span.Start), message)
	}
	indices := make(map[string]int, len(w.Nodes))
	for _, node := range w.Nodes {
		if node.Name == "end" {
			report(node, "end is a termination target, not a workflow node")
		} else if node.Name == "_" {
			report(node, "blank identifier cannot name a workflow node")
		}
		if _, exists := indices[node.Name]; exists {
			report(node, fmt.Sprintf("duplicate node %q", node.Name))
		}
		indices[node.Name] = len(plan.Nodes)
		plan.Nodes = append(plan.Nodes, NodePlan{Name: node, Directions: map[string]frontend.Identifier{}})
	}
	for _, edge := range w.Connections {
		from, exists := indices[edge.From.Name]
		if edge.From.Name == "end" {
			report(edge.From, "end cannot be the source of a connection")
			continue
		} else if !exists {
			report(edge.From, fmt.Sprintf("undeclared source node %q", edge.From.Name))
			continue
		}
		if _, exists := indices[edge.To.Name]; !exists && edge.To.Name != "end" {
			report(edge.To, fmt.Sprintf("undeclared destination node %q", edge.To.Name))
			continue
		}
		if _, exists := plan.Nodes[from].Directions[edge.Direction]; exists {
			report(edge.From, fmt.Sprintf("duplicate direction %q from node %q", edge.Direction, edge.From.Name))
			continue
		}
		plan.Nodes[from].Directions[edge.Direction] = edge.To
	}
	if len(diagnostics) > 0 {
		return plan, diagnostics
	}

	// Include every direction in graph validation, without treating end as a node.
	directed := make([][]frontend.Connection, len(w.Nodes))
	neighbors := make([][]int, len(w.Nodes))
	for _, edge := range w.Connections {
		if edge.To.Name == "end" {
			continue
		}
		from, to := indices[edge.From.Name], indices[edge.To.Name]
		directed[from] = append(directed[from], edge)
		neighbors[from] = append(neighbors[from], to)
		neighbors[to] = append(neighbors[to], from)
	}
	colors := make([]uint8, len(w.Nodes))
	var visit func(int) bool
	visit = func(node int) bool {
		colors[node] = 1
		for _, edge := range directed[node] {
			next := indices[edge.To.Name]
			if colors[next] == 1 {
				report(edge.To, fmt.Sprintf("directed cycle reaches node %q", edge.To.Name))
				return false
			}
			if colors[next] == 0 && !visit(next) {
				return false
			}
		}
		colors[node] = 2
		return true
	}
	for node := range w.Nodes {
		if colors[node] == 0 && !visit(node) {
			return plan, diagnostics
		}
	}
	seen := make([]bool, len(w.Nodes))
	seen[0] = true
	queue := []int{0}
	for head := 0; head < len(queue); head++ {
		for _, next := range neighbors[queue[head]] {
			if !seen[next] {
				seen[next] = true
				queue = append(queue, next)
			}
		}
	}
	for i, node := range w.Nodes {
		if !seen[i] {
			report(node, fmt.Sprintf("node %q is disconnected from entry node %q", node.Name, w.Nodes[0].Name))
		}
	}
	return plan, diagnostics
}
