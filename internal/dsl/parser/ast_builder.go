package parser

import (
	"github.com/JSYoo5B/chain/internal/dsl/ast"
	"github.com/antlr4-go/antlr/v4"
	"strconv"
	"strings"
)

type AstBuilder struct {
	*BaseChainParserListener
	rawTextStream antlr.TokenStream

	Result          ast.ChainCode
	currentWorkflow *ast.WorkflowStatement
}

func NewAstBuilder(stream antlr.TokenStream) *AstBuilder {
	return &AstBuilder{
		BaseChainParserListener: &BaseChainParserListener{},
		rawTextStream:           stream,
		Result:                  ast.ChainCode{},
		currentWorkflow:         nil,
	}
}

func (a *AstBuilder) EnterPackageClause(ctx *PackageClauseContext) {
	startToken := ctx.GetStart()

	a.Result.Package = ast.Package{
		CodeLocation: newCodeLocationFromToken(startToken),
		Name:         ctx.PackageName().GetText(),
	}
}

func (a *AstBuilder) EnterImportSpec(ctx *ImportSpecContext) {
	startToken := ctx.GetStart()

	var alias, path string
	pathWithQuotes := ctx.ImportPath().GetText()
	path, _ = strconv.Unquote(pathWithQuotes)

	if ctx.DOT() != nil {
		alias = "."
	} else if ctx.PackageName() != nil {
		alias = ctx.PackageName().GetText()
	} else {
		alias = inferPackageNameFromPath(path)
	}

	imp := ast.Import{
		CodeLocation: newCodeLocationFromToken(startToken),
		Alias:        alias,
		Path:         path,
	}

	a.Result.Imports = append(a.Result.Imports, imp)
}

func (a *AstBuilder) EnterWorkflowDecl(ctx *WorkflowDeclContext) {
	startToken := ctx.GetStart()

	workflow := ast.WorkflowStatement{
		CodeLocation: newCodeLocationFromToken(startToken),
	}

	a.currentWorkflow = &workflow
}

func (a *AstBuilder) ExitWorkflowDecl(_ *WorkflowDeclContext) {
	a.Result.Workflows = append(a.Result.Workflows, *a.currentWorkflow)
	a.currentWorkflow = nil
}

func (a *AstBuilder) EnterWorkflowSignature(ctx *WorkflowSignatureContext) {
	if a.currentWorkflow == nil {
		return
	}

	startToken := ctx.GetStart()

	paramStart, paramEnd := ctx.WorkflowParameters().GetStart(), ctx.WorkflowParameters().GetStop()
	paramText := a.rawTextStream.GetTextFromTokens(paramStart, paramEnd)

	typeStart, typeEnd := ctx.GetWorkflowType().GetStart(), ctx.GetWorkflowType().GetStop()
	typeText := a.rawTextStream.GetTextFromTokens(typeStart, typeEnd)

	declare := ast.WorkflowDeclaration{
		ConstructorName: newCodeLocationFromToken(ctx.GetWorkflowConstruct()),
		ConstructorParams: ast.CodeLocation{
			Line:   paramStart.GetLine(),
			Column: paramStart.GetColumn(),
			Text:   paramText,
		},
		WorkflowName: newCodeLocationFromToken(ctx.GetWorkflowName()),
		WorkflowType: ast.CodeLocation{
			Line:   typeStart.GetLine(),
			Column: typeStart.GetColumn(),
			Text:   typeText,
		},
		CodeLocation: newCodeLocationFromToken(startToken),
	}

	a.currentWorkflow.WorkflowDeclaration = declare
}

func (a *AstBuilder) EnterPrerequisiteBlock(ctx *PrerequisiteBlockContext) {
	if a.currentWorkflow == nil {
		return
	}

	start, end := ctx.L_CURLY().GetSymbol(), ctx.R_CURLY().GetSymbol()
	textWithBlock := a.rawTextStream.GetTextFromTokens(start, end)
	text := strings.TrimPrefix(textWithBlock, "{")
	text = strings.TrimSuffix(text, "}")

	a.currentWorkflow.PrerequisiteBlock = ast.PrerequisiteBlock{
		Code: text,
		CodeLocation: ast.CodeLocation{
			Line:   start.GetLine(),
			Column: start.GetColumn(),
			Text:   text,
		},
	}
}

func (a *AstBuilder) EnterNodesBlock(ctx *NodesBlockContext) {
	if a.currentWorkflow == nil {
		return
	}

	startToken, allNodeName := ctx.GetStart(), ctx.AllNodeName()
	nodes := make([]ast.WorkflowNode, 0, len(allNodeName))
	for _, nodeName := range allNodeName {
		node := ast.WorkflowNode{
			Name:         nodeName.GetText(),
			CodeLocation: newCodeLocationFromToken(nodeName.GetStart()),
		}
		nodes = append(nodes, node)
	}

	a.currentWorkflow.NodesBlock = ast.NodesBlock{
		Nodes:        nodes,
		CodeLocation: newCodeLocationFromToken(startToken),
	}
}

func (a *AstBuilder) EnterSuccessDirectionBlock(ctx *SuccessDirectionBlockContext) {
	if a.currentWorkflow == nil {
		return
	}

	allDirectionStatements := ctx.AllDirectionStmt()
	for _, directionStmt := range allDirectionStatements {
		a.currentWorkflow.Successes = append(a.currentWorkflow.Successes, parseDirections(directionStmt)...)
	}
}

func (a *AstBuilder) EnterFailureDirectionBlock(ctx *FailureDirectionBlockContext) {
	if a.currentWorkflow == nil {
		return
	}

	allDirectionStatements := ctx.AllDirectionStmt()
	for _, directionStmt := range allDirectionStatements {
		a.currentWorkflow.Failures = append(a.currentWorkflow.Failures, parseDirections(directionStmt)...)
	}
}

func (a *AstBuilder) EnterAbortDirectionBlock(ctx *AbortDirectionBlockContext) {
	if a.currentWorkflow == nil {
		return
	}

	allDirectionStatements := ctx.AllDirectionStmt()
	for _, directionStmt := range allDirectionStatements {
		a.currentWorkflow.Aborts = append(a.currentWorkflow.Aborts, parseDirections(directionStmt)...)
	}
}

func (a *AstBuilder) EnterBranchDirectionBlock(ctx *BranchDirectionBlockContext) {
	if a.currentWorkflow == nil {
		return
	}

	allBranchStatements := ctx.AllBranchStmt()
	for _, branchStmt := range allBranchStatements {
		a.currentWorkflow.Branches = append(a.currentWorkflow.Branches, parseBranch(branchStmt))
	}
}

func parseDirections(directionStatement IDirectionStmtContext) []ast.DirectionStatement {
	nodes, edges := directionStatement.AllNodeName(), directionStatement.AllEdgeDirection()
	directions := make([]ast.DirectionStatement, 0, len(edges))

	for i, edge := range edges {
		var from, to string
		left, right := nodes[i], nodes[i+1]
		if edge.GetText() == "-->" {
			from, to = left.GetText(), right.GetText()
		} else if edge.GetText() == "<--" {
			from, to = right.GetText(), left.GetText()
		}
		direction := ast.DirectionStatement{
			FromNode:     from,
			ToNode:       to,
			CodeLocation: newCodeLocationFromToken(left.GetStart()),
		}
		directions = append(directions, direction)
	}

	return directions
}

func parseBranch(branchStatement IBranchStmtContext) ast.BranchStatement {
	start := branchStatement.GetStart()
	from, to := branchStatement.GetSourceNode().GetText(), branchStatement.GetDestNode().GetText()
	conditionLiteral := branchStatement.GetBranchCond().GetText()
	condition, _ := strconv.Unquote(conditionLiteral[1 : len(conditionLiteral)-2])

	return ast.BranchStatement{
		FromNode:     from,
		Condition:    condition,
		ToNode:       to,
		CodeLocation: newCodeLocationFromToken(start),
	}
}
