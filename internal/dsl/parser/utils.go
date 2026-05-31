package parser

import (
	"github.com/JSYoo5B/chain/internal/dsl/ast"
	"github.com/antlr4-go/antlr/v4"
	"path/filepath"
	"regexp"
	"strings"
)

func newCodeLocationFromToken(tok antlr.Token) ast.CodeLocation {
	return ast.CodeLocation{
		Line:   tok.GetLine(),
		Column: tok.GetColumn(),
		Text:   tok.GetText(),
	}
}

var versionRegex = regexp.MustCompile(`^v[0-9]+$`)

func inferPackageNameFromPath(path string) string {
	base := filepath.Base(path)

	if versionRegex.MatchString(base) {
		parentDir := filepath.Dir(path)
		return filepath.Base(parentDir)
	}

	parts := strings.Split(base, ".")
	if len(parts) > 1 {
		lastPart := parts[len(parts)-1]
		if versionRegex.MatchString(lastPart) {
			return strings.Join(parts[:len(parts)-1], ".")
		}
	}

	return base
}
