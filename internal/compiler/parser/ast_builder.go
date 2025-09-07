package parser

import (
	"github.com/JSYoo5B/chain/internal/compiler/ast"
	"github.com/antlr4-go/antlr/v4"
	"strconv"
)

type AstBuilder struct {
	*BaseChainParserListener
	tokenStream antlr.TokenStream

	Result          ast.ChainCode
	currentWorkflow *ast.WorkflowStatement
}

func NewAstBuilder(stream antlr.TokenStream) *AstBuilder {
	return &AstBuilder{
		BaseChainParserListener: &BaseChainParserListener{},
		tokenStream:             stream,
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
	paramText := a.tokenStream.GetTextFromTokens(paramStart, paramEnd)

	typeStart, typeEnd := ctx.GetWorkflowType().GetStart(), ctx.GetWorkflowType().GetStop()
	typeText := a.tokenStream.GetTextFromTokens(typeStart, typeEnd)

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
