// Code generated from ChainParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // ChainParser
import "github.com/antlr4-go/antlr/v4"

type BaseChainParserVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseChainParserVisitor) VisitSourceFile(ctx *SourceFileContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitWorkflowDecl(ctx *WorkflowDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitWorkflowConstruct(ctx *WorkflowConstructContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitWorkflowSpec(ctx *WorkflowSpecContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitWorkflowBlock(ctx *WorkflowBlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitWorkflowStatementList(ctx *WorkflowStatementListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitPrerequisteStatements(ctx *PrerequisteStatementsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitPrerequisiteStmt(ctx *PrerequisiteStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitGolangEmbedStatement(ctx *GolangEmbedStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitNodesStatements(ctx *NodesStatementsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitNodeStmt(ctx *NodeStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitDirectionStatements(ctx *DirectionStatementsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitSuccessStatements(ctx *SuccessStatementsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitErrorStatements(ctx *ErrorStatementsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitAbortStatements(ctx *AbortStatementsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitBranchesStatements(ctx *BranchesStatementsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitDirectionStmt(ctx *DirectionStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitNodeName(ctx *NodeNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitBranchStmt(ctx *BranchStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitBranchDirection(ctx *BranchDirectionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitChain_eos(ctx *Chain_eosContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitPackageClause(ctx *PackageClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitPackageName(ctx *PackageNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitIdentifier(ctx *IdentifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitImportDecl(ctx *ImportDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitImportSpec(ctx *ImportSpecContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitImportPath(ctx *ImportPathContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitDeclaration(ctx *DeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitConstDecl(ctx *ConstDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitConstSpec(ctx *ConstSpecContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitIdentifierList(ctx *IdentifierListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitExpressionList(ctx *ExpressionListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitTypeDecl(ctx *TypeDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitTypeSpec(ctx *TypeSpecContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitAliasDecl(ctx *AliasDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitTypeDef(ctx *TypeDefContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitTypeParameters(ctx *TypeParametersContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitTypeParameterDecl(ctx *TypeParameterDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitTypeElement(ctx *TypeElementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitTypeTerm(ctx *TypeTermContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitFunctionDecl(ctx *FunctionDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitMethodDecl(ctx *MethodDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitReceiver(ctx *ReceiverContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitVarDecl(ctx *VarDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitVarSpec(ctx *VarSpecContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitBlock(ctx *BlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitStatementList(ctx *StatementListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitStatement(ctx *StatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitSimpleStmt(ctx *SimpleStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitExpressionStmt(ctx *ExpressionStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitSendStmt(ctx *SendStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitIncDecStmt(ctx *IncDecStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitAssignment(ctx *AssignmentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitAssign_op(ctx *Assign_opContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitShortVarDecl(ctx *ShortVarDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitLabeledStmt(ctx *LabeledStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitReturnStmt(ctx *ReturnStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitBreakStmt(ctx *BreakStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitContinueStmt(ctx *ContinueStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitGotoStmt(ctx *GotoStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitFallthroughStmt(ctx *FallthroughStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitDeferStmt(ctx *DeferStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitIfStmt(ctx *IfStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitSwitchStmt(ctx *SwitchStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitExprSwitchStmt(ctx *ExprSwitchStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitExprCaseClause(ctx *ExprCaseClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitExprSwitchCase(ctx *ExprSwitchCaseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitTypeSwitchStmt(ctx *TypeSwitchStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitTypeSwitchGuard(ctx *TypeSwitchGuardContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitTypeCaseClause(ctx *TypeCaseClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitTypeSwitchCase(ctx *TypeSwitchCaseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitTypeList(ctx *TypeListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitSelectStmt(ctx *SelectStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitCommClause(ctx *CommClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitCommCase(ctx *CommCaseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitRecvStmt(ctx *RecvStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitForStmt(ctx *ForStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitCondition(ctx *ConditionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitForClause(ctx *ForClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitRangeClause(ctx *RangeClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitGoStmt(ctx *GoStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitType_(ctx *Type_Context) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitTypeArgs(ctx *TypeArgsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitTypeName(ctx *TypeNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitTypeLit(ctx *TypeLitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitArrayType(ctx *ArrayTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitArrayLength(ctx *ArrayLengthContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitElementType(ctx *ElementTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitPointerType(ctx *PointerTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitInterfaceType(ctx *InterfaceTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitSliceType(ctx *SliceTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitMapType(ctx *MapTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitChannelType(ctx *ChannelTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitMethodSpec(ctx *MethodSpecContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitFunctionType(ctx *FunctionTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitSignature(ctx *SignatureContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitResult(ctx *ResultContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitParameters(ctx *ParametersContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitParameterDecl(ctx *ParameterDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitExpression(ctx *ExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitPrimaryExpr(ctx *PrimaryExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitConversion(ctx *ConversionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitOperand(ctx *OperandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitLiteral(ctx *LiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitBasicLit(ctx *BasicLitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitInteger(ctx *IntegerContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitOperandName(ctx *OperandNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitQualifiedIdent(ctx *QualifiedIdentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitCompositeLit(ctx *CompositeLitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitLiteralType(ctx *LiteralTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitLiteralValue(ctx *LiteralValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitElementList(ctx *ElementListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitKeyedElement(ctx *KeyedElementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitKey(ctx *KeyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitElement(ctx *ElementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitStructType(ctx *StructTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitFieldDecl(ctx *FieldDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitString_(ctx *String_Context) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitEmbeddedField(ctx *EmbeddedFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitFunctionLit(ctx *FunctionLitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitIndex(ctx *IndexContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitSlice_(ctx *Slice_Context) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitTypeAssertion(ctx *TypeAssertionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitArguments(ctx *ArgumentsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitMethodExpr(ctx *MethodExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseChainParserVisitor) VisitEos(ctx *EosContext) interface{} {
	return v.VisitChildren(ctx)
}
