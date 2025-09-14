SHELL := /bin/bash

PROJECT_DIR=$(shell pwd)
GRAMMAR_DIR=$(PROJECT_DIR)/internal/dsl/grammar
PARSER_DIR=$(PROJECT_DIR)/internal/dsl/parser

.PHONY: parser

parser: internal/dsl/parser/chain_parser.go

internal/dsl/parser/common_lexer.go: internal/dsl/grammar/CommonLexer.g4
	@echo "Generate Lexer code..."
	@cd $(GRAMMAR_DIR) && \
		antlr -Dlanguage=Go -o $(PARSER_DIR) -listener -package parser CommonLexer.g4

internal/dsl/parser/chain_parser.go: internal/dsl/grammar/ChainParser.g4 internal/dsl/grammar/GoParser.g4 internal/dsl/parser/common_lexer.go
	@echo "Generate Parser code..."
	@cd $(GRAMMAR_DIR) && \
		antlr -Dlanguage=Go -o $(PARSER_DIR) -lib $(PARSER_DIR) -listener -package parser ChainParser.g4 && \
	cd $(PARSER_DIR) && \
		awk '{gsub("this", "p"); print}' chain_parser.go > chain_parser_substitute.go && \
		mv chain_parser_substitute.go chain_parser.go