// Code generated from ChainParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // ChainParser
import "github.com/antlr4-go/antlr/v4"

// BaseChainParserListener is a complete listener for a parse tree produced by ChainParser.
type BaseChainParserListener struct{}

var _ ChainParserListener = &BaseChainParserListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseChainParserListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseChainParserListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseChainParserListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseChainParserListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterSourceFile is called when production sourceFile is entered.
func (s *BaseChainParserListener) EnterSourceFile(ctx *SourceFileContext) {}

// ExitSourceFile is called when production sourceFile is exited.
func (s *BaseChainParserListener) ExitSourceFile(ctx *SourceFileContext) {}

// EnterWorkflowDefine is called when production workflowDefine is entered.
func (s *BaseChainParserListener) EnterWorkflowDefine(ctx *WorkflowDefineContext) {}

// ExitWorkflowDefine is called when production workflowDefine is exited.
func (s *BaseChainParserListener) ExitWorkflowDefine(ctx *WorkflowDefineContext) {}

// EnterWorkflowDeclare is called when production workflowDeclare is entered.
func (s *BaseChainParserListener) EnterWorkflowDeclare(ctx *WorkflowDeclareContext) {}

// ExitWorkflowDeclare is called when production workflowDeclare is exited.
func (s *BaseChainParserListener) ExitWorkflowDeclare(ctx *WorkflowDeclareContext) {}

// EnterWorkflowParameters is called when production workflowParameters is entered.
func (s *BaseChainParserListener) EnterWorkflowParameters(ctx *WorkflowParametersContext) {}

// ExitWorkflowParameters is called when production workflowParameters is exited.
func (s *BaseChainParserListener) ExitWorkflowParameters(ctx *WorkflowParametersContext) {}

// EnterWorkflowBlock is called when production workflowBlock is entered.
func (s *BaseChainParserListener) EnterWorkflowBlock(ctx *WorkflowBlockContext) {}

// ExitWorkflowBlock is called when production workflowBlock is exited.
func (s *BaseChainParserListener) ExitWorkflowBlock(ctx *WorkflowBlockContext) {}

// EnterPrerequisteStatements is called when production prerequisteStatements is entered.
func (s *BaseChainParserListener) EnterPrerequisteStatements(ctx *PrerequisteStatementsContext) {}

// ExitPrerequisteStatements is called when production prerequisteStatements is exited.
func (s *BaseChainParserListener) ExitPrerequisteStatements(ctx *PrerequisteStatementsContext) {}

// EnterPrerequisiteStmt is called when production prerequisiteStmt is entered.
func (s *BaseChainParserListener) EnterPrerequisiteStmt(ctx *PrerequisiteStmtContext) {}

// ExitPrerequisiteStmt is called when production prerequisiteStmt is exited.
func (s *BaseChainParserListener) ExitPrerequisiteStmt(ctx *PrerequisiteStmtContext) {}

// EnterGolangEmbedStatement is called when production golangEmbedStatement is entered.
func (s *BaseChainParserListener) EnterGolangEmbedStatement(ctx *GolangEmbedStatementContext) {}

// ExitGolangEmbedStatement is called when production golangEmbedStatement is exited.
func (s *BaseChainParserListener) ExitGolangEmbedStatement(ctx *GolangEmbedStatementContext) {}

// EnterNodesStatements is called when production nodesStatements is entered.
func (s *BaseChainParserListener) EnterNodesStatements(ctx *NodesStatementsContext) {}

// ExitNodesStatements is called when production nodesStatements is exited.
func (s *BaseChainParserListener) ExitNodesStatements(ctx *NodesStatementsContext) {}

// EnterNodeStmt is called when production nodeStmt is entered.
func (s *BaseChainParserListener) EnterNodeStmt(ctx *NodeStmtContext) {}

// ExitNodeStmt is called when production nodeStmt is exited.
func (s *BaseChainParserListener) ExitNodeStmt(ctx *NodeStmtContext) {}

// EnterDirectionStatements is called when production directionStatements is entered.
func (s *BaseChainParserListener) EnterDirectionStatements(ctx *DirectionStatementsContext) {}

// ExitDirectionStatements is called when production directionStatements is exited.
func (s *BaseChainParserListener) ExitDirectionStatements(ctx *DirectionStatementsContext) {}

// EnterSuccessStatements is called when production successStatements is entered.
func (s *BaseChainParserListener) EnterSuccessStatements(ctx *SuccessStatementsContext) {}

// ExitSuccessStatements is called when production successStatements is exited.
func (s *BaseChainParserListener) ExitSuccessStatements(ctx *SuccessStatementsContext) {}

// EnterErrorStatements is called when production errorStatements is entered.
func (s *BaseChainParserListener) EnterErrorStatements(ctx *ErrorStatementsContext) {}

// ExitErrorStatements is called when production errorStatements is exited.
func (s *BaseChainParserListener) ExitErrorStatements(ctx *ErrorStatementsContext) {}

// EnterAbortStatements is called when production abortStatements is entered.
func (s *BaseChainParserListener) EnterAbortStatements(ctx *AbortStatementsContext) {}

// ExitAbortStatements is called when production abortStatements is exited.
func (s *BaseChainParserListener) ExitAbortStatements(ctx *AbortStatementsContext) {}

// EnterBranchesStatements is called when production branchesStatements is entered.
func (s *BaseChainParserListener) EnterBranchesStatements(ctx *BranchesStatementsContext) {}

// ExitBranchesStatements is called when production branchesStatements is exited.
func (s *BaseChainParserListener) ExitBranchesStatements(ctx *BranchesStatementsContext) {}

// EnterDirectionStmt is called when production directionStmt is entered.
func (s *BaseChainParserListener) EnterDirectionStmt(ctx *DirectionStmtContext) {}

// ExitDirectionStmt is called when production directionStmt is exited.
func (s *BaseChainParserListener) ExitDirectionStmt(ctx *DirectionStmtContext) {}

// EnterNodeName is called when production nodeName is entered.
func (s *BaseChainParserListener) EnterNodeName(ctx *NodeNameContext) {}

// ExitNodeName is called when production nodeName is exited.
func (s *BaseChainParserListener) ExitNodeName(ctx *NodeNameContext) {}

// EnterBranchStmt is called when production branchStmt is entered.
func (s *BaseChainParserListener) EnterBranchStmt(ctx *BranchStmtContext) {}

// ExitBranchStmt is called when production branchStmt is exited.
func (s *BaseChainParserListener) ExitBranchStmt(ctx *BranchStmtContext) {}

// EnterBranchDirection is called when production branchDirection is entered.
func (s *BaseChainParserListener) EnterBranchDirection(ctx *BranchDirectionContext) {}

// ExitBranchDirection is called when production branchDirection is exited.
func (s *BaseChainParserListener) ExitBranchDirection(ctx *BranchDirectionContext) {}

// EnterChain_eos is called when production chain_eos is entered.
func (s *BaseChainParserListener) EnterChain_eos(ctx *Chain_eosContext) {}

// ExitChain_eos is called when production chain_eos is exited.
func (s *BaseChainParserListener) ExitChain_eos(ctx *Chain_eosContext) {}

// EnterPackageClause is called when production packageClause is entered.
func (s *BaseChainParserListener) EnterPackageClause(ctx *PackageClauseContext) {}

// ExitPackageClause is called when production packageClause is exited.
func (s *BaseChainParserListener) ExitPackageClause(ctx *PackageClauseContext) {}

// EnterPackageName is called when production packageName is entered.
func (s *BaseChainParserListener) EnterPackageName(ctx *PackageNameContext) {}

// ExitPackageName is called when production packageName is exited.
func (s *BaseChainParserListener) ExitPackageName(ctx *PackageNameContext) {}

// EnterIdentifier is called when production identifier is entered.
func (s *BaseChainParserListener) EnterIdentifier(ctx *IdentifierContext) {}

// ExitIdentifier is called when production identifier is exited.
func (s *BaseChainParserListener) ExitIdentifier(ctx *IdentifierContext) {}

// EnterImportDecl is called when production importDecl is entered.
func (s *BaseChainParserListener) EnterImportDecl(ctx *ImportDeclContext) {}

// ExitImportDecl is called when production importDecl is exited.
func (s *BaseChainParserListener) ExitImportDecl(ctx *ImportDeclContext) {}

// EnterImportSpec is called when production importSpec is entered.
func (s *BaseChainParserListener) EnterImportSpec(ctx *ImportSpecContext) {}

// ExitImportSpec is called when production importSpec is exited.
func (s *BaseChainParserListener) ExitImportSpec(ctx *ImportSpecContext) {}

// EnterImportPath is called when production importPath is entered.
func (s *BaseChainParserListener) EnterImportPath(ctx *ImportPathContext) {}

// ExitImportPath is called when production importPath is exited.
func (s *BaseChainParserListener) ExitImportPath(ctx *ImportPathContext) {}

// EnterDeclaration is called when production declaration is entered.
func (s *BaseChainParserListener) EnterDeclaration(ctx *DeclarationContext) {}

// ExitDeclaration is called when production declaration is exited.
func (s *BaseChainParserListener) ExitDeclaration(ctx *DeclarationContext) {}

// EnterConstDecl is called when production constDecl is entered.
func (s *BaseChainParserListener) EnterConstDecl(ctx *ConstDeclContext) {}

// ExitConstDecl is called when production constDecl is exited.
func (s *BaseChainParserListener) ExitConstDecl(ctx *ConstDeclContext) {}

// EnterConstSpec is called when production constSpec is entered.
func (s *BaseChainParserListener) EnterConstSpec(ctx *ConstSpecContext) {}

// ExitConstSpec is called when production constSpec is exited.
func (s *BaseChainParserListener) ExitConstSpec(ctx *ConstSpecContext) {}

// EnterIdentifierList is called when production identifierList is entered.
func (s *BaseChainParserListener) EnterIdentifierList(ctx *IdentifierListContext) {}

// ExitIdentifierList is called when production identifierList is exited.
func (s *BaseChainParserListener) ExitIdentifierList(ctx *IdentifierListContext) {}

// EnterExpressionList is called when production expressionList is entered.
func (s *BaseChainParserListener) EnterExpressionList(ctx *ExpressionListContext) {}

// ExitExpressionList is called when production expressionList is exited.
func (s *BaseChainParserListener) ExitExpressionList(ctx *ExpressionListContext) {}

// EnterTypeDecl is called when production typeDecl is entered.
func (s *BaseChainParserListener) EnterTypeDecl(ctx *TypeDeclContext) {}

// ExitTypeDecl is called when production typeDecl is exited.
func (s *BaseChainParserListener) ExitTypeDecl(ctx *TypeDeclContext) {}

// EnterTypeSpec is called when production typeSpec is entered.
func (s *BaseChainParserListener) EnterTypeSpec(ctx *TypeSpecContext) {}

// ExitTypeSpec is called when production typeSpec is exited.
func (s *BaseChainParserListener) ExitTypeSpec(ctx *TypeSpecContext) {}

// EnterAliasDecl is called when production aliasDecl is entered.
func (s *BaseChainParserListener) EnterAliasDecl(ctx *AliasDeclContext) {}

// ExitAliasDecl is called when production aliasDecl is exited.
func (s *BaseChainParserListener) ExitAliasDecl(ctx *AliasDeclContext) {}

// EnterTypeDef is called when production typeDef is entered.
func (s *BaseChainParserListener) EnterTypeDef(ctx *TypeDefContext) {}

// ExitTypeDef is called when production typeDef is exited.
func (s *BaseChainParserListener) ExitTypeDef(ctx *TypeDefContext) {}

// EnterTypeParameters is called when production typeParameters is entered.
func (s *BaseChainParserListener) EnterTypeParameters(ctx *TypeParametersContext) {}

// ExitTypeParameters is called when production typeParameters is exited.
func (s *BaseChainParserListener) ExitTypeParameters(ctx *TypeParametersContext) {}

// EnterTypeParameterDecl is called when production typeParameterDecl is entered.
func (s *BaseChainParserListener) EnterTypeParameterDecl(ctx *TypeParameterDeclContext) {}

// ExitTypeParameterDecl is called when production typeParameterDecl is exited.
func (s *BaseChainParserListener) ExitTypeParameterDecl(ctx *TypeParameterDeclContext) {}

// EnterTypeElement is called when production typeElement is entered.
func (s *BaseChainParserListener) EnterTypeElement(ctx *TypeElementContext) {}

// ExitTypeElement is called when production typeElement is exited.
func (s *BaseChainParserListener) ExitTypeElement(ctx *TypeElementContext) {}

// EnterTypeTerm is called when production typeTerm is entered.
func (s *BaseChainParserListener) EnterTypeTerm(ctx *TypeTermContext) {}

// ExitTypeTerm is called when production typeTerm is exited.
func (s *BaseChainParserListener) ExitTypeTerm(ctx *TypeTermContext) {}

// EnterFunctionDecl is called when production functionDecl is entered.
func (s *BaseChainParserListener) EnterFunctionDecl(ctx *FunctionDeclContext) {}

// ExitFunctionDecl is called when production functionDecl is exited.
func (s *BaseChainParserListener) ExitFunctionDecl(ctx *FunctionDeclContext) {}

// EnterMethodDecl is called when production methodDecl is entered.
func (s *BaseChainParserListener) EnterMethodDecl(ctx *MethodDeclContext) {}

// ExitMethodDecl is called when production methodDecl is exited.
func (s *BaseChainParserListener) ExitMethodDecl(ctx *MethodDeclContext) {}

// EnterReceiver is called when production receiver is entered.
func (s *BaseChainParserListener) EnterReceiver(ctx *ReceiverContext) {}

// ExitReceiver is called when production receiver is exited.
func (s *BaseChainParserListener) ExitReceiver(ctx *ReceiverContext) {}

// EnterVarDecl is called when production varDecl is entered.
func (s *BaseChainParserListener) EnterVarDecl(ctx *VarDeclContext) {}

// ExitVarDecl is called when production varDecl is exited.
func (s *BaseChainParserListener) ExitVarDecl(ctx *VarDeclContext) {}

// EnterVarSpec is called when production varSpec is entered.
func (s *BaseChainParserListener) EnterVarSpec(ctx *VarSpecContext) {}

// ExitVarSpec is called when production varSpec is exited.
func (s *BaseChainParserListener) ExitVarSpec(ctx *VarSpecContext) {}

// EnterBlock is called when production block is entered.
func (s *BaseChainParserListener) EnterBlock(ctx *BlockContext) {}

// ExitBlock is called when production block is exited.
func (s *BaseChainParserListener) ExitBlock(ctx *BlockContext) {}

// EnterStatementList is called when production statementList is entered.
func (s *BaseChainParserListener) EnterStatementList(ctx *StatementListContext) {}

// ExitStatementList is called when production statementList is exited.
func (s *BaseChainParserListener) ExitStatementList(ctx *StatementListContext) {}

// EnterStatement is called when production statement is entered.
func (s *BaseChainParserListener) EnterStatement(ctx *StatementContext) {}

// ExitStatement is called when production statement is exited.
func (s *BaseChainParserListener) ExitStatement(ctx *StatementContext) {}

// EnterSimpleStmt is called when production simpleStmt is entered.
func (s *BaseChainParserListener) EnterSimpleStmt(ctx *SimpleStmtContext) {}

// ExitSimpleStmt is called when production simpleStmt is exited.
func (s *BaseChainParserListener) ExitSimpleStmt(ctx *SimpleStmtContext) {}

// EnterExpressionStmt is called when production expressionStmt is entered.
func (s *BaseChainParserListener) EnterExpressionStmt(ctx *ExpressionStmtContext) {}

// ExitExpressionStmt is called when production expressionStmt is exited.
func (s *BaseChainParserListener) ExitExpressionStmt(ctx *ExpressionStmtContext) {}

// EnterSendStmt is called when production sendStmt is entered.
func (s *BaseChainParserListener) EnterSendStmt(ctx *SendStmtContext) {}

// ExitSendStmt is called when production sendStmt is exited.
func (s *BaseChainParserListener) ExitSendStmt(ctx *SendStmtContext) {}

// EnterIncDecStmt is called when production incDecStmt is entered.
func (s *BaseChainParserListener) EnterIncDecStmt(ctx *IncDecStmtContext) {}

// ExitIncDecStmt is called when production incDecStmt is exited.
func (s *BaseChainParserListener) ExitIncDecStmt(ctx *IncDecStmtContext) {}

// EnterAssignment is called when production assignment is entered.
func (s *BaseChainParserListener) EnterAssignment(ctx *AssignmentContext) {}

// ExitAssignment is called when production assignment is exited.
func (s *BaseChainParserListener) ExitAssignment(ctx *AssignmentContext) {}

// EnterAssign_op is called when production assign_op is entered.
func (s *BaseChainParserListener) EnterAssign_op(ctx *Assign_opContext) {}

// ExitAssign_op is called when production assign_op is exited.
func (s *BaseChainParserListener) ExitAssign_op(ctx *Assign_opContext) {}

// EnterShortVarDecl is called when production shortVarDecl is entered.
func (s *BaseChainParserListener) EnterShortVarDecl(ctx *ShortVarDeclContext) {}

// ExitShortVarDecl is called when production shortVarDecl is exited.
func (s *BaseChainParserListener) ExitShortVarDecl(ctx *ShortVarDeclContext) {}

// EnterLabeledStmt is called when production labeledStmt is entered.
func (s *BaseChainParserListener) EnterLabeledStmt(ctx *LabeledStmtContext) {}

// ExitLabeledStmt is called when production labeledStmt is exited.
func (s *BaseChainParserListener) ExitLabeledStmt(ctx *LabeledStmtContext) {}

// EnterReturnStmt is called when production returnStmt is entered.
func (s *BaseChainParserListener) EnterReturnStmt(ctx *ReturnStmtContext) {}

// ExitReturnStmt is called when production returnStmt is exited.
func (s *BaseChainParserListener) ExitReturnStmt(ctx *ReturnStmtContext) {}

// EnterBreakStmt is called when production breakStmt is entered.
func (s *BaseChainParserListener) EnterBreakStmt(ctx *BreakStmtContext) {}

// ExitBreakStmt is called when production breakStmt is exited.
func (s *BaseChainParserListener) ExitBreakStmt(ctx *BreakStmtContext) {}

// EnterContinueStmt is called when production continueStmt is entered.
func (s *BaseChainParserListener) EnterContinueStmt(ctx *ContinueStmtContext) {}

// ExitContinueStmt is called when production continueStmt is exited.
func (s *BaseChainParserListener) ExitContinueStmt(ctx *ContinueStmtContext) {}

// EnterGotoStmt is called when production gotoStmt is entered.
func (s *BaseChainParserListener) EnterGotoStmt(ctx *GotoStmtContext) {}

// ExitGotoStmt is called when production gotoStmt is exited.
func (s *BaseChainParserListener) ExitGotoStmt(ctx *GotoStmtContext) {}

// EnterFallthroughStmt is called when production fallthroughStmt is entered.
func (s *BaseChainParserListener) EnterFallthroughStmt(ctx *FallthroughStmtContext) {}

// ExitFallthroughStmt is called when production fallthroughStmt is exited.
func (s *BaseChainParserListener) ExitFallthroughStmt(ctx *FallthroughStmtContext) {}

// EnterDeferStmt is called when production deferStmt is entered.
func (s *BaseChainParserListener) EnterDeferStmt(ctx *DeferStmtContext) {}

// ExitDeferStmt is called when production deferStmt is exited.
func (s *BaseChainParserListener) ExitDeferStmt(ctx *DeferStmtContext) {}

// EnterIfStmt is called when production ifStmt is entered.
func (s *BaseChainParserListener) EnterIfStmt(ctx *IfStmtContext) {}

// ExitIfStmt is called when production ifStmt is exited.
func (s *BaseChainParserListener) ExitIfStmt(ctx *IfStmtContext) {}

// EnterSwitchStmt is called when production switchStmt is entered.
func (s *BaseChainParserListener) EnterSwitchStmt(ctx *SwitchStmtContext) {}

// ExitSwitchStmt is called when production switchStmt is exited.
func (s *BaseChainParserListener) ExitSwitchStmt(ctx *SwitchStmtContext) {}

// EnterExprSwitchStmt is called when production exprSwitchStmt is entered.
func (s *BaseChainParserListener) EnterExprSwitchStmt(ctx *ExprSwitchStmtContext) {}

// ExitExprSwitchStmt is called when production exprSwitchStmt is exited.
func (s *BaseChainParserListener) ExitExprSwitchStmt(ctx *ExprSwitchStmtContext) {}

// EnterExprCaseClause is called when production exprCaseClause is entered.
func (s *BaseChainParserListener) EnterExprCaseClause(ctx *ExprCaseClauseContext) {}

// ExitExprCaseClause is called when production exprCaseClause is exited.
func (s *BaseChainParserListener) ExitExprCaseClause(ctx *ExprCaseClauseContext) {}

// EnterExprSwitchCase is called when production exprSwitchCase is entered.
func (s *BaseChainParserListener) EnterExprSwitchCase(ctx *ExprSwitchCaseContext) {}

// ExitExprSwitchCase is called when production exprSwitchCase is exited.
func (s *BaseChainParserListener) ExitExprSwitchCase(ctx *ExprSwitchCaseContext) {}

// EnterTypeSwitchStmt is called when production typeSwitchStmt is entered.
func (s *BaseChainParserListener) EnterTypeSwitchStmt(ctx *TypeSwitchStmtContext) {}

// ExitTypeSwitchStmt is called when production typeSwitchStmt is exited.
func (s *BaseChainParserListener) ExitTypeSwitchStmt(ctx *TypeSwitchStmtContext) {}

// EnterTypeSwitchGuard is called when production typeSwitchGuard is entered.
func (s *BaseChainParserListener) EnterTypeSwitchGuard(ctx *TypeSwitchGuardContext) {}

// ExitTypeSwitchGuard is called when production typeSwitchGuard is exited.
func (s *BaseChainParserListener) ExitTypeSwitchGuard(ctx *TypeSwitchGuardContext) {}

// EnterTypeCaseClause is called when production typeCaseClause is entered.
func (s *BaseChainParserListener) EnterTypeCaseClause(ctx *TypeCaseClauseContext) {}

// ExitTypeCaseClause is called when production typeCaseClause is exited.
func (s *BaseChainParserListener) ExitTypeCaseClause(ctx *TypeCaseClauseContext) {}

// EnterTypeSwitchCase is called when production typeSwitchCase is entered.
func (s *BaseChainParserListener) EnterTypeSwitchCase(ctx *TypeSwitchCaseContext) {}

// ExitTypeSwitchCase is called when production typeSwitchCase is exited.
func (s *BaseChainParserListener) ExitTypeSwitchCase(ctx *TypeSwitchCaseContext) {}

// EnterTypeList is called when production typeList is entered.
func (s *BaseChainParserListener) EnterTypeList(ctx *TypeListContext) {}

// ExitTypeList is called when production typeList is exited.
func (s *BaseChainParserListener) ExitTypeList(ctx *TypeListContext) {}

// EnterSelectStmt is called when production selectStmt is entered.
func (s *BaseChainParserListener) EnterSelectStmt(ctx *SelectStmtContext) {}

// ExitSelectStmt is called when production selectStmt is exited.
func (s *BaseChainParserListener) ExitSelectStmt(ctx *SelectStmtContext) {}

// EnterCommClause is called when production commClause is entered.
func (s *BaseChainParserListener) EnterCommClause(ctx *CommClauseContext) {}

// ExitCommClause is called when production commClause is exited.
func (s *BaseChainParserListener) ExitCommClause(ctx *CommClauseContext) {}

// EnterCommCase is called when production commCase is entered.
func (s *BaseChainParserListener) EnterCommCase(ctx *CommCaseContext) {}

// ExitCommCase is called when production commCase is exited.
func (s *BaseChainParserListener) ExitCommCase(ctx *CommCaseContext) {}

// EnterRecvStmt is called when production recvStmt is entered.
func (s *BaseChainParserListener) EnterRecvStmt(ctx *RecvStmtContext) {}

// ExitRecvStmt is called when production recvStmt is exited.
func (s *BaseChainParserListener) ExitRecvStmt(ctx *RecvStmtContext) {}

// EnterForStmt is called when production forStmt is entered.
func (s *BaseChainParserListener) EnterForStmt(ctx *ForStmtContext) {}

// ExitForStmt is called when production forStmt is exited.
func (s *BaseChainParserListener) ExitForStmt(ctx *ForStmtContext) {}

// EnterCondition is called when production condition is entered.
func (s *BaseChainParserListener) EnterCondition(ctx *ConditionContext) {}

// ExitCondition is called when production condition is exited.
func (s *BaseChainParserListener) ExitCondition(ctx *ConditionContext) {}

// EnterForClause is called when production forClause is entered.
func (s *BaseChainParserListener) EnterForClause(ctx *ForClauseContext) {}

// ExitForClause is called when production forClause is exited.
func (s *BaseChainParserListener) ExitForClause(ctx *ForClauseContext) {}

// EnterRangeClause is called when production rangeClause is entered.
func (s *BaseChainParserListener) EnterRangeClause(ctx *RangeClauseContext) {}

// ExitRangeClause is called when production rangeClause is exited.
func (s *BaseChainParserListener) ExitRangeClause(ctx *RangeClauseContext) {}

// EnterGoStmt is called when production goStmt is entered.
func (s *BaseChainParserListener) EnterGoStmt(ctx *GoStmtContext) {}

// ExitGoStmt is called when production goStmt is exited.
func (s *BaseChainParserListener) ExitGoStmt(ctx *GoStmtContext) {}

// EnterType_ is called when production type_ is entered.
func (s *BaseChainParserListener) EnterType_(ctx *Type_Context) {}

// ExitType_ is called when production type_ is exited.
func (s *BaseChainParserListener) ExitType_(ctx *Type_Context) {}

// EnterTypeArgs is called when production typeArgs is entered.
func (s *BaseChainParserListener) EnterTypeArgs(ctx *TypeArgsContext) {}

// ExitTypeArgs is called when production typeArgs is exited.
func (s *BaseChainParserListener) ExitTypeArgs(ctx *TypeArgsContext) {}

// EnterTypeName is called when production typeName is entered.
func (s *BaseChainParserListener) EnterTypeName(ctx *TypeNameContext) {}

// ExitTypeName is called when production typeName is exited.
func (s *BaseChainParserListener) ExitTypeName(ctx *TypeNameContext) {}

// EnterTypeLit is called when production typeLit is entered.
func (s *BaseChainParserListener) EnterTypeLit(ctx *TypeLitContext) {}

// ExitTypeLit is called when production typeLit is exited.
func (s *BaseChainParserListener) ExitTypeLit(ctx *TypeLitContext) {}

// EnterArrayType is called when production arrayType is entered.
func (s *BaseChainParserListener) EnterArrayType(ctx *ArrayTypeContext) {}

// ExitArrayType is called when production arrayType is exited.
func (s *BaseChainParserListener) ExitArrayType(ctx *ArrayTypeContext) {}

// EnterArrayLength is called when production arrayLength is entered.
func (s *BaseChainParserListener) EnterArrayLength(ctx *ArrayLengthContext) {}

// ExitArrayLength is called when production arrayLength is exited.
func (s *BaseChainParserListener) ExitArrayLength(ctx *ArrayLengthContext) {}

// EnterElementType is called when production elementType is entered.
func (s *BaseChainParserListener) EnterElementType(ctx *ElementTypeContext) {}

// ExitElementType is called when production elementType is exited.
func (s *BaseChainParserListener) ExitElementType(ctx *ElementTypeContext) {}

// EnterPointerType is called when production pointerType is entered.
func (s *BaseChainParserListener) EnterPointerType(ctx *PointerTypeContext) {}

// ExitPointerType is called when production pointerType is exited.
func (s *BaseChainParserListener) ExitPointerType(ctx *PointerTypeContext) {}

// EnterInterfaceType is called when production interfaceType is entered.
func (s *BaseChainParserListener) EnterInterfaceType(ctx *InterfaceTypeContext) {}

// ExitInterfaceType is called when production interfaceType is exited.
func (s *BaseChainParserListener) ExitInterfaceType(ctx *InterfaceTypeContext) {}

// EnterSliceType is called when production sliceType is entered.
func (s *BaseChainParserListener) EnterSliceType(ctx *SliceTypeContext) {}

// ExitSliceType is called when production sliceType is exited.
func (s *BaseChainParserListener) ExitSliceType(ctx *SliceTypeContext) {}

// EnterMapType is called when production mapType is entered.
func (s *BaseChainParserListener) EnterMapType(ctx *MapTypeContext) {}

// ExitMapType is called when production mapType is exited.
func (s *BaseChainParserListener) ExitMapType(ctx *MapTypeContext) {}

// EnterChannelType is called when production channelType is entered.
func (s *BaseChainParserListener) EnterChannelType(ctx *ChannelTypeContext) {}

// ExitChannelType is called when production channelType is exited.
func (s *BaseChainParserListener) ExitChannelType(ctx *ChannelTypeContext) {}

// EnterMethodSpec is called when production methodSpec is entered.
func (s *BaseChainParserListener) EnterMethodSpec(ctx *MethodSpecContext) {}

// ExitMethodSpec is called when production methodSpec is exited.
func (s *BaseChainParserListener) ExitMethodSpec(ctx *MethodSpecContext) {}

// EnterFunctionType is called when production functionType is entered.
func (s *BaseChainParserListener) EnterFunctionType(ctx *FunctionTypeContext) {}

// ExitFunctionType is called when production functionType is exited.
func (s *BaseChainParserListener) ExitFunctionType(ctx *FunctionTypeContext) {}

// EnterSignature is called when production signature is entered.
func (s *BaseChainParserListener) EnterSignature(ctx *SignatureContext) {}

// ExitSignature is called when production signature is exited.
func (s *BaseChainParserListener) ExitSignature(ctx *SignatureContext) {}

// EnterResult is called when production result is entered.
func (s *BaseChainParserListener) EnterResult(ctx *ResultContext) {}

// ExitResult is called when production result is exited.
func (s *BaseChainParserListener) ExitResult(ctx *ResultContext) {}

// EnterParameters is called when production parameters is entered.
func (s *BaseChainParserListener) EnterParameters(ctx *ParametersContext) {}

// ExitParameters is called when production parameters is exited.
func (s *BaseChainParserListener) ExitParameters(ctx *ParametersContext) {}

// EnterParameterDecl is called when production parameterDecl is entered.
func (s *BaseChainParserListener) EnterParameterDecl(ctx *ParameterDeclContext) {}

// ExitParameterDecl is called when production parameterDecl is exited.
func (s *BaseChainParserListener) ExitParameterDecl(ctx *ParameterDeclContext) {}

// EnterExpression is called when production expression is entered.
func (s *BaseChainParserListener) EnterExpression(ctx *ExpressionContext) {}

// ExitExpression is called when production expression is exited.
func (s *BaseChainParserListener) ExitExpression(ctx *ExpressionContext) {}

// EnterPrimaryExpr is called when production primaryExpr is entered.
func (s *BaseChainParserListener) EnterPrimaryExpr(ctx *PrimaryExprContext) {}

// ExitPrimaryExpr is called when production primaryExpr is exited.
func (s *BaseChainParserListener) ExitPrimaryExpr(ctx *PrimaryExprContext) {}

// EnterConversion is called when production conversion is entered.
func (s *BaseChainParserListener) EnterConversion(ctx *ConversionContext) {}

// ExitConversion is called when production conversion is exited.
func (s *BaseChainParserListener) ExitConversion(ctx *ConversionContext) {}

// EnterOperand is called when production operand is entered.
func (s *BaseChainParserListener) EnterOperand(ctx *OperandContext) {}

// ExitOperand is called when production operand is exited.
func (s *BaseChainParserListener) ExitOperand(ctx *OperandContext) {}

// EnterLiteral is called when production literal is entered.
func (s *BaseChainParserListener) EnterLiteral(ctx *LiteralContext) {}

// ExitLiteral is called when production literal is exited.
func (s *BaseChainParserListener) ExitLiteral(ctx *LiteralContext) {}

// EnterBasicLit is called when production basicLit is entered.
func (s *BaseChainParserListener) EnterBasicLit(ctx *BasicLitContext) {}

// ExitBasicLit is called when production basicLit is exited.
func (s *BaseChainParserListener) ExitBasicLit(ctx *BasicLitContext) {}

// EnterInteger is called when production integer is entered.
func (s *BaseChainParserListener) EnterInteger(ctx *IntegerContext) {}

// ExitInteger is called when production integer is exited.
func (s *BaseChainParserListener) ExitInteger(ctx *IntegerContext) {}

// EnterOperandName is called when production operandName is entered.
func (s *BaseChainParserListener) EnterOperandName(ctx *OperandNameContext) {}

// ExitOperandName is called when production operandName is exited.
func (s *BaseChainParserListener) ExitOperandName(ctx *OperandNameContext) {}

// EnterQualifiedIdent is called when production qualifiedIdent is entered.
func (s *BaseChainParserListener) EnterQualifiedIdent(ctx *QualifiedIdentContext) {}

// ExitQualifiedIdent is called when production qualifiedIdent is exited.
func (s *BaseChainParserListener) ExitQualifiedIdent(ctx *QualifiedIdentContext) {}

// EnterCompositeLit is called when production compositeLit is entered.
func (s *BaseChainParserListener) EnterCompositeLit(ctx *CompositeLitContext) {}

// ExitCompositeLit is called when production compositeLit is exited.
func (s *BaseChainParserListener) ExitCompositeLit(ctx *CompositeLitContext) {}

// EnterLiteralType is called when production literalType is entered.
func (s *BaseChainParserListener) EnterLiteralType(ctx *LiteralTypeContext) {}

// ExitLiteralType is called when production literalType is exited.
func (s *BaseChainParserListener) ExitLiteralType(ctx *LiteralTypeContext) {}

// EnterLiteralValue is called when production literalValue is entered.
func (s *BaseChainParserListener) EnterLiteralValue(ctx *LiteralValueContext) {}

// ExitLiteralValue is called when production literalValue is exited.
func (s *BaseChainParserListener) ExitLiteralValue(ctx *LiteralValueContext) {}

// EnterElementList is called when production elementList is entered.
func (s *BaseChainParserListener) EnterElementList(ctx *ElementListContext) {}

// ExitElementList is called when production elementList is exited.
func (s *BaseChainParserListener) ExitElementList(ctx *ElementListContext) {}

// EnterKeyedElement is called when production keyedElement is entered.
func (s *BaseChainParserListener) EnterKeyedElement(ctx *KeyedElementContext) {}

// ExitKeyedElement is called when production keyedElement is exited.
func (s *BaseChainParserListener) ExitKeyedElement(ctx *KeyedElementContext) {}

// EnterKey is called when production key is entered.
func (s *BaseChainParserListener) EnterKey(ctx *KeyContext) {}

// ExitKey is called when production key is exited.
func (s *BaseChainParserListener) ExitKey(ctx *KeyContext) {}

// EnterElement is called when production element is entered.
func (s *BaseChainParserListener) EnterElement(ctx *ElementContext) {}

// ExitElement is called when production element is exited.
func (s *BaseChainParserListener) ExitElement(ctx *ElementContext) {}

// EnterStructType is called when production structType is entered.
func (s *BaseChainParserListener) EnterStructType(ctx *StructTypeContext) {}

// ExitStructType is called when production structType is exited.
func (s *BaseChainParserListener) ExitStructType(ctx *StructTypeContext) {}

// EnterFieldDecl is called when production fieldDecl is entered.
func (s *BaseChainParserListener) EnterFieldDecl(ctx *FieldDeclContext) {}

// ExitFieldDecl is called when production fieldDecl is exited.
func (s *BaseChainParserListener) ExitFieldDecl(ctx *FieldDeclContext) {}

// EnterString_ is called when production string_ is entered.
func (s *BaseChainParserListener) EnterString_(ctx *String_Context) {}

// ExitString_ is called when production string_ is exited.
func (s *BaseChainParserListener) ExitString_(ctx *String_Context) {}

// EnterEmbeddedField is called when production embeddedField is entered.
func (s *BaseChainParserListener) EnterEmbeddedField(ctx *EmbeddedFieldContext) {}

// ExitEmbeddedField is called when production embeddedField is exited.
func (s *BaseChainParserListener) ExitEmbeddedField(ctx *EmbeddedFieldContext) {}

// EnterFunctionLit is called when production functionLit is entered.
func (s *BaseChainParserListener) EnterFunctionLit(ctx *FunctionLitContext) {}

// ExitFunctionLit is called when production functionLit is exited.
func (s *BaseChainParserListener) ExitFunctionLit(ctx *FunctionLitContext) {}

// EnterIndex is called when production index is entered.
func (s *BaseChainParserListener) EnterIndex(ctx *IndexContext) {}

// ExitIndex is called when production index is exited.
func (s *BaseChainParserListener) ExitIndex(ctx *IndexContext) {}

// EnterSlice_ is called when production slice_ is entered.
func (s *BaseChainParserListener) EnterSlice_(ctx *Slice_Context) {}

// ExitSlice_ is called when production slice_ is exited.
func (s *BaseChainParserListener) ExitSlice_(ctx *Slice_Context) {}

// EnterTypeAssertion is called when production typeAssertion is entered.
func (s *BaseChainParserListener) EnterTypeAssertion(ctx *TypeAssertionContext) {}

// ExitTypeAssertion is called when production typeAssertion is exited.
func (s *BaseChainParserListener) ExitTypeAssertion(ctx *TypeAssertionContext) {}

// EnterArguments is called when production arguments is entered.
func (s *BaseChainParserListener) EnterArguments(ctx *ArgumentsContext) {}

// ExitArguments is called when production arguments is exited.
func (s *BaseChainParserListener) ExitArguments(ctx *ArgumentsContext) {}

// EnterMethodExpr is called when production methodExpr is entered.
func (s *BaseChainParserListener) EnterMethodExpr(ctx *MethodExprContext) {}

// ExitMethodExpr is called when production methodExpr is exited.
func (s *BaseChainParserListener) ExitMethodExpr(ctx *MethodExprContext) {}

// EnterEos is called when production eos is entered.
func (s *BaseChainParserListener) EnterEos(ctx *EosContext) {}

// ExitEos is called when production eos is exited.
func (s *BaseChainParserListener) ExitEos(ctx *EosContext) {}
