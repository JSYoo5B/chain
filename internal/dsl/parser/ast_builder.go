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

	a.currentWorkflow.Prerequisite = ast.PrerequisiteBlock{
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
