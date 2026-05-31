SHELL := /bin/bash

PROJECT_DIR=$(shell pwd)
GRAMMAR_DIR=$(PROJECT_DIR)/internal/dsl/grammar
PARSER_DIR=$(PROJECT_DIR)/internal/dsl/parser

.PHONY: parser

parser: $(PARSER_DIR)/chain_parser.go $(PARSER_DIR)/golang/go_parser.go

$(PARSER_DIR)/common_lexer.go: $(GRAMMAR_DIR)/CommonLexer.g4
	@echo "Generate Lexer code for chain dsl"
	@cd $(GRAMMAR_DIR) && \
		antlr -Dlanguage=Go -o $(PARSER_DIR) -listener -package parser CommonLexer.g4

$(PARSER_DIR)/chain_parser.go: $(GRAMMAR_DIR)/ChainParser.g4 $(GRAMMAR_DIR)/GoParser.g4 $(PARSER_DIR)/common_lexer.go
	@echo "Generate Parser code for chain dsl"
	@cd $(GRAMMAR_DIR) && \
		antlr -Dlanguage=Go -o $(PARSER_DIR) -lib $(PARSER_DIR) -listener -package parser ChainParser.g4 && \
	cd $(PARSER_DIR) && \
		awk '{gsub("this", "p"); print}' chain_parser.go > chain_parser_substitute.go && \
		mv chain_parser_substitute.go chain_parser.go

$(PARSER_DIR)/golang/common_lexer.go: $(GRAMMAR_DIR)/CommonLexer.g4
	@echo "Generate Lexer code for embedding golang"
	@cd $(GRAMMAR_DIR) && \
		antlr -Dlanguage=Go -o $(PARSER_DIR)/golang -listener -package golang CommonLexer.g4

$(PARSER_DIR)/golang/go_parser.go: $(GRAMMAR_DIR)/GoParser.g4 $(PARSER_DIR)/golang/common_lexer.go
	@echo "Generate Parser code for embedding golang"
	@cd $(GRAMMAR_DIR) && \
		antlr -Dlanguage=Go -o $(PARSER_DIR)/golang -lib $(PARSER_DIR)/golang -listener -package golang GoParser.g4 && \
	cd $(PARSER_DIR)/golang && \
		awk '{gsub("this", "p"); print}' go_parser.go > go_parser_substitute.go && \
		mv go_parser_substitute.go go_parser.go
