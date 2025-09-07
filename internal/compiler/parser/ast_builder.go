package parser

import (
	"github.com/JSYoo5B/chain/internal/compiler/ast"
	"strconv"
)

type AstBuilder struct {
	*BaseChainParserListener

	Result          ast.ChainCode
	currentWorkflow *ast.WorkflowStatement
}

func NewAstBuilder() *AstBuilder {
	return &AstBuilder{
		BaseChainParserListener: &BaseChainParserListener{},
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
