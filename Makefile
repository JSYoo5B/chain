SHELL := /bin/bash

PROJECT_DIR=$(shell pwd)
GRAMMAR_DIR=$(PROJECT_DIR)/internal/dsl/grammar
ANTLR_CODEGEN_DIR=$(PROJECT_DIR)/internal/dsl/parser

.PHONY: antlr

antlr:
	cd $(GRAMMAR_DIR) && \
	antlr -Dlanguage=Go -o ../parser -visitor -listener -package parser CommonLexer.g4 && \
	antlr -Dlanguage=Go -o ../parser -lib ../parser -visitor -listener -package parser GoParser.g4 && \
	cd $(ANTLR_CODEGEN_DIR) && \
	awk '{gsub("this", "p"); print}' go_parser.go > go_parser_substitute.go && \
	mv go_parser_substitute.go go_parser.go && \
	cd $(PROJECT_DIR)

