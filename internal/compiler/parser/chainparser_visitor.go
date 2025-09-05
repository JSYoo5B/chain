// Code generated from ChainParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // ChainParser
import "github.com/antlr4-go/antlr/v4"

// A complete Visitor for a parse tree produced by ChainParser.
type ChainParserVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by ChainParser#sourceFile.
	VisitSourceFile(ctx *SourceFileContext) interface{}

	// Visit a parse tree produced by ChainParser#workflowDecl.
	VisitWorkflowDecl(ctx *WorkflowDeclContext) interface{}

	// Visit a parse tree produced by ChainParser#workflowConstruct.
	VisitWorkflowConstruct(ctx *WorkflowConstructContext) interface{}

	// Visit a parse tree produced by ChainParser#workflowSpec.
	VisitWorkflowSpec(ctx *WorkflowSpecContext) interface{}

	// Visit a parse tree produced by ChainParser#workflowBlock.
	VisitWorkflowBlock(ctx *WorkflowBlockContext) interface{}

	// Visit a parse tree produced by ChainParser#workflowStatementList.
	VisitWorkflowStatementList(ctx *WorkflowStatementListContext) interface{}

	// Visit a parse tree produced by ChainParser#prerequisteStatements.
	VisitPrerequisteStatements(ctx *PrerequisteStatementsContext) interface{}

	// Visit a parse tree produced by ChainParser#prerequisiteStmt.
	VisitPrerequisiteStmt(ctx *PrerequisiteStmtContext) interface{}

	// Visit a parse tree produced by ChainParser#golangEmbedStatement.
	VisitGolangEmbedStatement(ctx *GolangEmbedStatementContext) interface{}

	// Visit a parse tree produced by ChainParser#nodesStatements.
	VisitNodesStatements(ctx *NodesStatementsContext) interface{}

	// Visit a parse tree produced by ChainParser#nodeStmt.
	VisitNodeStmt(ctx *NodeStmtContext) interface{}

	// Visit a parse tree produced by ChainParser#directionStatements.
	VisitDirectionStatements(ctx *DirectionStatementsContext) interface{}

	// Visit a parse tree produced by ChainParser#successStatements.
	VisitSuccessStatements(ctx *SuccessStatementsContext) interface{}

	// Visit a parse tree produced by ChainParser#errorStatements.
	VisitErrorStatements(ctx *ErrorStatementsContext) interface{}

	// Visit a parse tree produced by ChainParser#abortStatements.
	VisitAbortStatements(ctx *AbortStatementsContext) interface{}

	// Visit a parse tree produced by ChainParser#branchesStatements.
	VisitBranchesStatements(ctx *BranchesStatementsContext) interface{}

	// Visit a parse tree produced by ChainParser#directionStmt.
	VisitDirectionStmt(ctx *DirectionStmtContext) interface{}

	// Visit a parse tree produced by ChainParser#nodeName.
	VisitNodeName(ctx *NodeNameContext) interface{}

	// Visit a parse tree produced by ChainParser#branchStmt.
	VisitBranchStmt(ctx *BranchStmtContext) interface{}

	// Visit a parse tree produced by ChainParser#branchDirection.
	VisitBranchDirection(ctx *BranchDirectionContext) interface{}

	// Visit a parse tree produced by ChainParser#chain_eos.
	VisitChain_eos(ctx *Chain_eosContext) interface{}

	// Visit a parse tree produced by ChainParser#packageClause.
	VisitPackageClause(ctx *PackageClauseContext) interface{}

	// Visit a parse tree produced by ChainParser#packageName.
	VisitPackageName(ctx *PackageNameContext) interface{}

	// Visit a parse tree produced by ChainParser#identifier.
	VisitIdentifier(ctx *IdentifierContext) interface{}

	// Visit a parse tree produced by ChainParser#importDecl.
	VisitImportDecl(ctx *ImportDeclContext) interface{}

	// Visit a parse tree produced by ChainParser#importSpec.
	VisitImportSpec(ctx *ImportSpecContext) interface{}

	// Visit a parse tree produced by ChainParser#importPath.
	VisitImportPath(ctx *ImportPathContext) interface{}

	// Visit a parse tree produced by ChainParser#declaration.
	VisitDeclaration(ctx *DeclarationContext) interface{}

	// Visit a parse tree produced by ChainParser#constDecl.
	VisitConstDecl(ctx *ConstDeclContext) interface{}

	// Visit a parse tree produced by ChainParser#constSpec.
	VisitConstSpec(ctx *ConstSpecContext) interface{}

	// Visit a parse tree produced by ChainParser#identifierList.
	VisitIdentifierList(ctx *IdentifierListContext) interface{}

	// Visit a parse tree produced by ChainParser#expressionList.
	VisitExpressionList(ctx *ExpressionListContext) interface{}

	// Visit a parse tree produced by ChainParser#typeDecl.
	VisitTypeDecl(ctx *TypeDeclContext) interface{}

	// Visit a parse tree produced by ChainParser#typeSpec.
	VisitTypeSpec(ctx *TypeSpecContext) interface{}

	// Visit a parse tree produced by ChainParser#aliasDecl.
	VisitAliasDecl(ctx *AliasDeclContext) interface{}

	// Visit a parse tree produced by ChainParser#typeDef.
	VisitTypeDef(ctx *TypeDefContext) interface{}

	// Visit a parse tree produced by ChainParser#typeParameters.
	VisitTypeParameters(ctx *TypeParametersContext) interface{}

	// Visit a parse tree produced by ChainParser#typeParameterDecl.
	VisitTypeParameterDecl(ctx *TypeParameterDeclContext) interface{}

	// Visit a parse tree produced by ChainParser#typeElement.
	VisitTypeElement(ctx *TypeElementContext) interface{}

	// Visit a parse tree produced by ChainParser#typeTerm.
	VisitTypeTerm(ctx *TypeTermContext) interface{}

	// Visit a parse tree produced by ChainParser#functionDecl.
	VisitFunctionDecl(ctx *FunctionDeclContext) interface{}

	// Visit a parse tree produced by ChainParser#methodDecl.
	VisitMethodDecl(ctx *MethodDeclContext) interface{}

	// Visit a parse tree produced by ChainParser#receiver.
	VisitReceiver(ctx *ReceiverContext) interface{}

	// Visit a parse tree produced by ChainParser#varDecl.
	VisitVarDecl(ctx *VarDeclContext) interface{}

	// Visit a parse tree produced by ChainParser#varSpec.
	VisitVarSpec(ctx *VarSpecContext) interface{}

	// Visit a parse tree produced by ChainParser#block.
	VisitBlock(ctx *BlockContext) interface{}

	// Visit a parse tree produced by ChainParser#statementList.
	VisitStatementList(ctx *StatementListContext) interface{}

	// Visit a parse tree produced by ChainParser#statement.
	VisitStatement(ctx *StatementContext) interface{}

	// Visit a parse tree produced by ChainParser#simpleStmt.
	VisitSimpleStmt(ctx *SimpleStmtContext) interface{}

	// Visit a parse tree produced by ChainParser#expressionStmt.
	VisitExpressionStmt(ctx *ExpressionStmtContext) interface{}

	// Visit a parse tree produced by ChainParser#sendStmt.
	VisitSendStmt(ctx *SendStmtContext) interface{}

	// Visit a parse tree produced by ChainParser#incDecStmt.
	VisitIncDecStmt(ctx *IncDecStmtContext) interface{}

	// Visit a parse tree produced by ChainParser#assignment.
	VisitAssignment(ctx *AssignmentContext) interface{}

	// Visit a parse tree produced by ChainParser#assign_op.
	VisitAssign_op(ctx *Assign_opContext) interface{}

	// Visit a parse tree produced by ChainParser#shortVarDecl.
	VisitShortVarDecl(ctx *ShortVarDeclContext) interface{}

	// Visit a parse tree produced by ChainParser#labeledStmt.
	VisitLabeledStmt(ctx *LabeledStmtContext) interface{}

	// Visit a parse tree produced by ChainParser#returnStmt.
	VisitReturnStmt(ctx *ReturnStmtContext) interface{}

	// Visit a parse tree produced by ChainParser#breakStmt.
	VisitBreakStmt(ctx *BreakStmtContext) interface{}

	// Visit a parse tree produced by ChainParser#continueStmt.
	VisitContinueStmt(ctx *ContinueStmtContext) interface{}

	// Visit a parse tree produced by ChainParser#gotoStmt.
	VisitGotoStmt(ctx *GotoStmtContext) interface{}

	// Visit a parse tree produced by ChainParser#fallthroughStmt.
	VisitFallthroughStmt(ctx *FallthroughStmtContext) interface{}

	// Visit a parse tree produced by ChainParser#deferStmt.
	VisitDeferStmt(ctx *DeferStmtContext) interface{}

	// Visit a parse tree produced by ChainParser#ifStmt.
	VisitIfStmt(ctx *IfStmtContext) interface{}

	// Visit a parse tree produced by ChainParser#switchStmt.
	VisitSwitchStmt(ctx *SwitchStmtContext) interface{}

	// Visit a parse tree produced by ChainParser#exprSwitchStmt.
	VisitExprSwitchStmt(ctx *ExprSwitchStmtContext) interface{}

	// Visit a parse tree produced by ChainParser#exprCaseClause.
	VisitExprCaseClause(ctx *ExprCaseClauseContext) interface{}

	// Visit a parse tree produced by ChainParser#exprSwitchCase.
	VisitExprSwitchCase(ctx *ExprSwitchCaseContext) interface{}

	// Visit a parse tree produced by ChainParser#typeSwitchStmt.
	VisitTypeSwitchStmt(ctx *TypeSwitchStmtContext) interface{}

	// Visit a parse tree produced by ChainParser#typeSwitchGuard.
	VisitTypeSwitchGuard(ctx *TypeSwitchGuardContext) interface{}

	// Visit a parse tree produced by ChainParser#typeCaseClause.
	VisitTypeCaseClause(ctx *TypeCaseClauseContext) interface{}

	// Visit a parse tree produced by ChainParser#typeSwitchCase.
	VisitTypeSwitchCase(ctx *TypeSwitchCaseContext) interface{}

	// Visit a parse tree produced by ChainParser#typeList.
	VisitTypeList(ctx *TypeListContext) interface{}

	// Visit a parse tree produced by ChainParser#selectStmt.
	VisitSelectStmt(ctx *SelectStmtContext) interface{}

	// Visit a parse tree produced by ChainParser#commClause.
	VisitCommClause(ctx *CommClauseContext) interface{}

	// Visit a parse tree produced by ChainParser#commCase.
	VisitCommCase(ctx *CommCaseContext) interface{}

	// Visit a parse tree produced by ChainParser#recvStmt.
	VisitRecvStmt(ctx *RecvStmtContext) interface{}

	// Visit a parse tree produced by ChainParser#forStmt.
	VisitForStmt(ctx *ForStmtContext) interface{}

	// Visit a parse tree produced by ChainParser#condition.
	VisitCondition(ctx *ConditionContext) interface{}

	// Visit a parse tree produced by ChainParser#forClause.
	VisitForClause(ctx *ForClauseContext) interface{}

	// Visit a parse tree produced by ChainParser#rangeClause.
	VisitRangeClause(ctx *RangeClauseContext) interface{}

	// Visit a parse tree produced by ChainParser#goStmt.
	VisitGoStmt(ctx *GoStmtContext) interface{}

	// Visit a parse tree produced by ChainParser#type_.
	VisitType_(ctx *Type_Context) interface{}

	// Visit a parse tree produced by ChainParser#typeArgs.
	VisitTypeArgs(ctx *TypeArgsContext) interface{}

	// Visit a parse tree produced by ChainParser#typeName.
	VisitTypeName(ctx *TypeNameContext) interface{}

	// Visit a parse tree produced by ChainParser#typeLit.
	VisitTypeLit(ctx *TypeLitContext) interface{}

	// Visit a parse tree produced by ChainParser#arrayType.
	VisitArrayType(ctx *ArrayTypeContext) interface{}

	// Visit a parse tree produced by ChainParser#arrayLength.
	VisitArrayLength(ctx *ArrayLengthContext) interface{}

	// Visit a parse tree produced by ChainParser#elementType.
	VisitElementType(ctx *ElementTypeContext) interface{}

	// Visit a parse tree produced by ChainParser#pointerType.
	VisitPointerType(ctx *PointerTypeContext) interface{}

	// Visit a parse tree produced by ChainParser#interfaceType.
	VisitInterfaceType(ctx *InterfaceTypeContext) interface{}

	// Visit a parse tree produced by ChainParser#sliceType.
	VisitSliceType(ctx *SliceTypeContext) interface{}

	// Visit a parse tree produced by ChainParser#mapType.
	VisitMapType(ctx *MapTypeContext) interface{}

	// Visit a parse tree produced by ChainParser#channelType.
	VisitChannelType(ctx *ChannelTypeContext) interface{}

	// Visit a parse tree produced by ChainParser#methodSpec.
	VisitMethodSpec(ctx *MethodSpecContext) interface{}

	// Visit a parse tree produced by ChainParser#functionType.
	VisitFunctionType(ctx *FunctionTypeContext) interface{}

	// Visit a parse tree produced by ChainParser#signature.
	VisitSignature(ctx *SignatureContext) interface{}

	// Visit a parse tree produced by ChainParser#result.
	VisitResult(ctx *ResultContext) interface{}

	// Visit a parse tree produced by ChainParser#parameters.
	VisitParameters(ctx *ParametersContext) interface{}

	// Visit a parse tree produced by ChainParser#parameterDecl.
	VisitParameterDecl(ctx *ParameterDeclContext) interface{}

	// Visit a parse tree produced by ChainParser#expression.
	VisitExpression(ctx *ExpressionContext) interface{}

	// Visit a parse tree produced by ChainParser#primaryExpr.
	VisitPrimaryExpr(ctx *PrimaryExprContext) interface{}

	// Visit a parse tree produced by ChainParser#conversion.
	VisitConversion(ctx *ConversionContext) interface{}

	// Visit a parse tree produced by ChainParser#operand.
	VisitOperand(ctx *OperandContext) interface{}

	// Visit a parse tree produced by ChainParser#literal.
	VisitLiteral(ctx *LiteralContext) interface{}

	// Visit a parse tree produced by ChainParser#basicLit.
	VisitBasicLit(ctx *BasicLitContext) interface{}

	// Visit a parse tree produced by ChainParser#integer.
	VisitInteger(ctx *IntegerContext) interface{}

	// Visit a parse tree produced by ChainParser#operandName.
	VisitOperandName(ctx *OperandNameContext) interface{}

	// Visit a parse tree produced by ChainParser#qualifiedIdent.
	VisitQualifiedIdent(ctx *QualifiedIdentContext) interface{}

	// Visit a parse tree produced by ChainParser#compositeLit.
	VisitCompositeLit(ctx *CompositeLitContext) interface{}

	// Visit a parse tree produced by ChainParser#literalType.
	VisitLiteralType(ctx *LiteralTypeContext) interface{}

	// Visit a parse tree produced by ChainParser#literalValue.
	VisitLiteralValue(ctx *LiteralValueContext) interface{}

	// Visit a parse tree produced by ChainParser#elementList.
	VisitElementList(ctx *ElementListContext) interface{}

	// Visit a parse tree produced by ChainParser#keyedElement.
	VisitKeyedElement(ctx *KeyedElementContext) interface{}

	// Visit a parse tree produced by ChainParser#key.
	VisitKey(ctx *KeyContext) interface{}

	// Visit a parse tree produced by ChainParser#element.
	VisitElement(ctx *ElementContext) interface{}

	// Visit a parse tree produced by ChainParser#structType.
	VisitStructType(ctx *StructTypeContext) interface{}

	// Visit a parse tree produced by ChainParser#fieldDecl.
	VisitFieldDecl(ctx *FieldDeclContext) interface{}

	// Visit a parse tree produced by ChainParser#string_.
	VisitString_(ctx *String_Context) interface{}

	// Visit a parse tree produced by ChainParser#embeddedField.
	VisitEmbeddedField(ctx *EmbeddedFieldContext) interface{}

	// Visit a parse tree produced by ChainParser#functionLit.
	VisitFunctionLit(ctx *FunctionLitContext) interface{}

	// Visit a parse tree produced by ChainParser#index.
	VisitIndex(ctx *IndexContext) interface{}

	// Visit a parse tree produced by ChainParser#slice_.
	VisitSlice_(ctx *Slice_Context) interface{}

	// Visit a parse tree produced by ChainParser#typeAssertion.
	VisitTypeAssertion(ctx *TypeAssertionContext) interface{}

	// Visit a parse tree produced by ChainParser#arguments.
	VisitArguments(ctx *ArgumentsContext) interface{}

	// Visit a parse tree produced by ChainParser#methodExpr.
	VisitMethodExpr(ctx *MethodExprContext) interface{}

	// Visit a parse tree produced by ChainParser#eos.
	VisitEos(ctx *EosContext) interface{}
}
