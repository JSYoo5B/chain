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
    : WORKFLOW workflowSignature L_CURLY workflowBody R_CURLY
    ;

workflowSignature
    : workflowConstruct=IDENTIFIER
        L_PAREN workflowParameters R_PAREN
        workflowName=IDENTIFIER
        L_BRACKET workflowType=typeElement R_BRACKET
    ;

workflowBody
    : prerequisiteBlock?
        nodesBlock
        (
            successDirectionBlock
            | errorDirectionBlock
            | abortDirectionBlock
            | branchDirectionBlock
        )+
    ;

workflowParameters
    : (parameterDecl (COMMA parameterDecl)* COMMA?)?
    ;

prerequisiteBlock
    : PREREQUISITE L_CURLY prerequisiteStmt R_CURLY EOS*
    ;

nodesBlock
    : NODES COLON nodeName (COMMA nodeName)+ chain_eos
    ;

successDirectionBlock
    : SUCCESS COLON (directionStmt chain_eos?)+
    ;

errorDirectionBlock
    : ERROR COLON (directionStmt chain_eos?)+
    ;

abortDirectionBlock
    : ABORT COLON (directionStmt chain_eos?)+
    ;

branchDirectionBlock
    : BRANCHES COLON (branchStmt chain_eos?)+
    ;

prerequisiteStmt
    : (golangEmbedStatement chain_eos)*
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

nodeName
    : END
    | IDENTIFIER
    ;

directionStmt
    : nodeName (edgeDirection nodeName)+
    ;

edgeDirection
    : L_TO_R | R_TO_L
    ;

branchStmt
    : sourceNode=IDENTIFIER  branchCond=BRANCH_DIR destNode=IDENTIFIER
    ;

chain_eos
    : EOS
    ;