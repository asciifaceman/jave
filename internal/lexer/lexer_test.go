package lexer_test

import (
	"testing"

	"github.com/asciifaceman/jave/internal/lexer"
	"github.com/asciifaceman/jave/internal/token"
)

func TestLex_AssignAndStatementEnd(t *testing.T) {
	src := `allow exact Count 2b=2 5;;`
	tokens, diags := lexer.Lex(src)
	if len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %d", len(diags))
	}

	kinds := kindsOf(tokens)
	want := []token.Kind{
		token.Allow,
		token.TypeExact,
		token.Identifier,
		token.Assign,
		token.Integer,
		token.StmtEnd,
		token.EOF,
	}
	assertKinds(t, kinds, want)
}

func TestLex_HelloWorldShape(t *testing.T) {
	src := `outy seq Foremost<> --> <<nada>> {
    pront("hello, jave");;
    give up;;
}`
	tokens, diags := lexer.Lex(src)
	if len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %d", len(diags))
	}

	if len(tokens) < 10 {
		t.Fatalf("expected enough tokens, got %d", len(tokens))
	}
	if tokens[0].Kind != token.Outy {
		t.Fatalf("expected first token Outy, got %s", tokens[0].Kind)
	}
	if tokens[0].Pos.Line != 1 || tokens[0].Pos.Column != 1 {
		t.Fatalf("unexpected position for first token: %d:%d", tokens[0].Pos.Line, tokens[0].Pos.Column)
	}
}

func TestLex_UnterminatedStringReportsDiagnostic(t *testing.T) {
	src := `pront("hello);;`
	_, diags := lexer.Lex(src)
	if len(diags) == 0 {
		t.Fatal("expected diagnostics, got none")
	}
	if diags[0].Message != "unterminated string literal" {
		t.Fatalf("unexpected diagnostic message: %q", diags[0].Message)
	}
}

func kindsOf(tokens []token.Token) []token.Kind {
	out := make([]token.Kind, 0, len(tokens))
	for _, tok := range tokens {
		out = append(out, tok.Kind)
	}
	return out
}

func assertKinds(t *testing.T, got []token.Kind, want []token.Kind) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("token count mismatch: got %d want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("token kind mismatch at %d: got %s want %s", i, got[i], want[i])
		}
	}
}
