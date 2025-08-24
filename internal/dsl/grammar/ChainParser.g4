


// $antlr-format alignTrailingComments true, columnLimit 150, minEmptyLines 1, maxEmptyLinesToKeep 1, reflowComments false, useTab false
// $antlr-format allowShortRulesOnASingleLine false, allowShortBlocksOnASingleLine true, alignSemicolons hanging, alignColons hanging

parser grammar ChainParser;

options {
    tokenVocab = CommonLexer;
    superClass = ChainParserBase;
}

import GoParser;

sourceFile
    : packageClause EOS (importDecl eos)* (workflowDecl eos)* EOF
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
    : nodesStatements EOS (directionStatements)*
    ;

nodesStatements
    : NODES COLON (nodeStmt eos)*
    ;

nodeStmt
    : shortVarDecl
    ;

directionStatements
    : successStatements
    | errorStatements
    | abortStatements
    | branchesStatements
    ;

successStatements
    : SUCCESS COLON (directionStmt eos)*
    ;

errorStatements
    : ERROR COLON (directionStmt eos)*
    ;

abortStatements
    : ABORT COLON (directionStmt eos)*
    ;

branchesStatements
    : BRANCHES COLON (branchStmt eos)*
    ;

directionStmt
    : nodeName (direction nodeName)+
    ;

nodeName
    : IDENTIFIER
    | END
    ;

direction
    : L_TO_R
    | R_TO_L
    ;

branchStmt
    : nodeName (branchDirection nodeName)*
    ;

branchDirection
    : MINUS string_ MINUS GREATER
    ;