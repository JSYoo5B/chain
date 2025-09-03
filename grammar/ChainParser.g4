// $antlr-format alignTrailingComments true, columnLimit 150, minEmptyLines 1, maxEmptyLinesToKeep 1, reflowComments false, useTab false
// $antlr-format allowShortRulesOnASingleLine false, allowShortBlocksOnASingleLine true, alignSemicolons hanging, alignColons hanging

parser grammar ChainParser;

options {
    tokenVocab = CommonLexer;
    superClass = ChainParserBase;
}

import GoParser;

sourceFile
    : packageClause eos (importDecl eos)* (workflowDecl eos)* EOF
    ;

workflowDecl
    : WORKFLOW workflowConstruct parameters GENERATES? workflowSpec workflowBlock
    ;

workflowConstruct
    : IDENTIFIER
    ;

workflowSpec
    : IDENTIFIER L_BRACKET typeElement R_BRACKET
    ;

workflowBlock
    : L_CURLY workflowStatementList R_CURLY
    ;

workflowStatementList
    : prerequisteStatements nodesStatements directionStatements*
    ;

prerequisteStatements
    : PREREQUISITE L_CURLY prerequisiteStmt R_CURLY EOS*
    ;

prerequisiteStmt
    : ( (SEMI | EOS | /* {this.closingBracket()}? */ ) golangEmbedStatement eos)*
    ;

golangEmbedStatement
    : declaration
    | labeledStmt
    | simpleStmt
    | goStmt
    | breakStmt
    | continueStmt
    | gotoStmt
    | fallthroughStmt
    | block
    | ifStmt
    | switchStmt
    | selectStmt
    | forStmt
    | deferStmt
    ;

nodesStatements
    : NODES COLON (nodeStmt eos)* EOS*
    ;

nodeStmt
    : identifierList
    ;

directionStatements
    : successStatements
    | errorStatements
    | abortStatements
    | branchesStatements
    ;

successStatements
    : SUCCESS COLON (directionStmt chain_eos?)+
    ;

errorStatements
    : ERROR COLON (directionStmt chain_eos?)+
    ;

abortStatements
    : ABORT COLON (directionStmt chain_eos?)+
    ;

branchesStatements
    : BRANCHES COLON (branchStmt chain_eos?)+
    ;

directionStmt
    : nodeName (direction = (L_TO_R | R_TO_L) nodeName)*
    ;

nodeName
    : END
    | IDENTIFIER
    ;

branchStmt
    : IDENTIFIER branchDirection nodeName
    ;

branchDirection
    : MINUS string_ MINUS GREATER
    ;

chain_eos
    : EOS
    ;