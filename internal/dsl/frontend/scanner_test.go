package frontend

import (
	"errors"
	"go/scanner"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoRegionPreservesSource(t *testing.T) {
	cases := map[string]string{
		"empty":                          "()",
		"constructor parameters":         "(a, b int, callback func(x, y int) (int, error), options ...Option)",
		"DSL words in parameters":        "(success, end int, nodes map[string]Option)",
		"nested state type":              "[map[string][]pkg.Box[struct{ N int }]]",
		"DSL words in body":              "{\n success := action()\n end := success\n _ = end\n}",
		"nested literals and closures":   "{\n f := func() { _ = struct{ N int }{N: 1} }\n _ = f\n}",
		"quoted delimiters":              "{\n text := \"} \\\" nodes: ( ]\"\n rune := '}'\n _ = text\n _ = rune\n}",
		"raw delimiters and DSL text":    "{\n text := `} prerequisite { nodes: --> -\"even\"->`\n _ = text\n}",
		"comments containing delimiters": "{\n // } ) ] nodes:\n /* { ( [ */\n a := action() /* } */\n}",
		"CRLF and raw string":            "{\r\n text := `first\r\n} nodes:\r\nlast`\r\n _ = text\r\n}",
		"multiline parameters":           "(\n a, b int,\n callback func(int) error,\n options ...Option,\n)",
		// Syntax belongs to go/parser; extraction only establishes boundaries.
		"invalid Go with balanced delimiters": "{ a := }",
	}
	for name, regionText := range cases {
		t.Run(name, func(t *testing.T) {
			prefix := "// 한글 prefix\n입력 "
			source := NewSource("example.chain", []byte(prefix+regionText+" tail"))
			c, err := scan(source)
			require.NoError(t, err)
			require.Equal(t, "입력", c.current().literal)
			c.advance()

			region, err := c.takeGoRegion()
			require.NoError(t, err)
			assert.Equal(t, Span{Start: len(prefix), End: len(prefix) + len(regionText)}, region.Span)
			assert.Equal(t, regionText, source.Text(region.Span))
			assert.Equal(t, regionText[1:len(regionText)-1], source.Text(region.Content))
			assert.Equal(t, "tail", c.current().literal)
		})
	}
}

func TestGoRegionRejectsBrokenBoundaries(t *testing.T) {
	cases := []struct {
		text    string
		offset  int
		message string
	}{
		{"name", 0, "expected '(', '[' or '{'"},
		{"", 0, "expected '(', '[' or '{'"},
		{"}", 0, "expected '(', '[' or '{'"},
		{"{", 1, "expected \"}\""},
		{"([", 2, "expected \"]\""},
		{"([)]", 2, "expected \"]\""},
		{"{\n f := func() {\n", 17, "expected \"}\""},
		{"{\n f := func(}\n", 13, "expected \")\""},
	}
	for _, tc := range cases {
		t.Run(tc.text, func(t *testing.T) {
			c, err := scan(NewSource("broken.chain", []byte(tc.text)))
			require.NoError(t, err)
			_, err = c.takeGoRegion()
			require.ErrorContains(t, err, tc.message)
			var diagnostic *scanner.Error
			require.True(t, errors.As(err, &diagnostic))
			assert.Equal(t, "broken.chain", diagnostic.Pos.Filename)
			assert.Equal(t, tc.offset, diagnostic.Pos.Offset)
			if tc.offset > 0 {
				assert.Contains(t, diagnostic.Msg, "at broken.chain:")
			}
		})
	}
}

func TestScannerRejectsLexicalErrors(t *testing.T) {
	for _, text := range []string{
		"{ text := \"unfinished }",
		"{ text := `unfinished }",
		"{ /* unfinished }",
		"{ text := \"\\q\" }",
		"{ \x00 }",
	} {
		t.Run(text, func(t *testing.T) {
			c, err := scan(NewSource("broken.chain", []byte(text)))
			require.Error(t, err)
			assert.Nil(t, c, "do not expose a successful token stream after lexical errors")
		})
	}
}

func TestScannerKeepsDSLWordsAsIdentifiers(t *testing.T) {
	c, err := scan(NewSource("keywords.chain", []byte("workflow prerequisite nodes branches success failure abort end")))
	require.NoError(t, err)
	var words []string
	for c.current().kind == token.IDENT {
		words = append(words, c.current().literal)
		c.advance()
	}
	assert.Equal(t, []string{"workflow", "prerequisite", "nodes", "branches", "success", "failure", "abort", "end"}, words)
}

func TestScannerKeepsArrowOffsets(t *testing.T) {
	// The DSL layer will combine these tokens only when their source positions
	// are adjacent. The scanner must retain the distinction from spaced tokens.
	c, err := scan(NewSource("arrows.chain", []byte("--> <-- -\"even\"-> -`odd`-> -- >")))
	require.NoError(t, err)
	var kinds []token.Token
	var offsets []int
	for c.current().kind != token.EOF {
		kinds = append(kinds, c.current().kind)
		offsets = append(offsets, c.current().offset)
		c.advance()
	}
	assert.Equal(t, []token.Token{
		token.DEC, token.GTR, token.ARROW, token.SUB,
		token.SUB, token.STRING, token.SUB, token.GTR,
		token.SUB, token.STRING, token.SUB, token.GTR,
		token.DEC, token.GTR,
	}, kinds)
	assert.Equal(t, []int{0, 2, 4, 6, 8, 9, 15, 16, 18, 19, 24, 25, 27, 30}, offsets)
}

func TestScannerDistinguishesSemicolons(t *testing.T) {
	c, err := scan(NewSource("semis.chain", []byte("a;\nb\n")))
	require.NoError(t, err)
	var semicolons []lexeme
	for c.current().kind != token.EOF {
		if c.current().kind == token.SEMICOLON {
			semicolons = append(semicolons, c.current())
		}
		c.advance()
	}
	assert.Equal(t, []lexeme{
		{kind: token.SEMICOLON, offset: 1, literal: ";"},
		{kind: token.SEMICOLON, offset: 4, literal: "\n"},
	}, semicolons)
}

func TestPhysicalPositionsIgnoreLineDirectives(t *testing.T) {
	const text = "//line other.go:900\n{\n]"
	source := NewSource("original.chain", []byte(text))
	c, err := scan(source)
	require.NoError(t, err)
	_, err = c.takeGoRegion()
	require.ErrorContains(t, err, "original.chain:3:1")
	assert.NotContains(t, err.Error(), "other.go")

	_, err = scan(NewSource("original.chain", []byte("//line other.go:900\n\x00")))
	require.ErrorContains(t, err, "original.chain:2:1")
	assert.NotContains(t, err.Error(), "other.go")
}

func TestSourceOwnsItsBytes(t *testing.T) {
	input := []byte("(a int)")
	source := NewSource("copy.chain", input)
	input[1] = 'b'
	assert.Equal(t, "(a int)", source.Text(Span{Start: 0, End: len(input)}))
}
