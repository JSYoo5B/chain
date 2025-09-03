SHELL := /bin/bash

PROJECT_DIR=$(shell pwd)
GRAMMAR_DIR=$(PROJECT_DIR)/grammar
ANTLR_CODEGEN_DIR=$(PROJECT_DIR)/internal/parser

.PHONY: antlr

antlr:
	cd $(GRAMMAR_DIR) && \
	antlr -Dlanguage=Go -o ../internal/parser -visitor -listener -package parser CommonLexer.g4 && \
	antlr -Dlanguage=Go -o ../internal/parser -lib ../internal/parser -visitor -listener -package parser ChainParser.g4 && \
	cd $(ANTLR_CODEGEN_DIR) && \
	awk '{gsub("this", "p"); print}' chain_parser.go > chain_parser_substitute.go && \
	mv chain_parser_substitute.go chain_parser.go && \
	cd $(PROJECT_DIR)