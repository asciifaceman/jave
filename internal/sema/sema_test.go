package sema_test

import (
	"testing"

	"github.com/asciifaceman/jave/internal/lexer"
	"github.com/asciifaceman/jave/internal/parser"
	"github.com/asciifaceman/jave/internal/sema"
)

func TestAnalyze_ValidForemost(t *testing.T) {
	src := `outy seq Foremost<> --> <<nada>> {
    give up;;
}`
	diags := analyzeSrc(t, src)
	if len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %d", len(diags))
	}
}

func TestAnalyze_MissingForemost(t *testing.T) {
	src := `outy seq Add<> --> <<exact>> {
    give 1 up;;
}`
	diags := analyzeSrc(t, src)
	if len(diags) == 0 {
		t.Fatal("expected diagnostics, got none")
	}
}

func TestAnalyze_NonNadaMustReturnValue(t *testing.T) {
	src := `outy seq Foremost<> --> <<nada>> {
    give up;;
}

outy seq Add<> --> <<exact>> {
    give up;;
}`
	diags := analyzeSrc(t, src)
	if len(diags) == 0 {
		t.Fatal("expected diagnostics, got none")
	}
}

func TestAnalyze_UndefinedIdentifier(t *testing.T) {
	src := `outy seq Foremost<> --> <<nada>> {
    pront(UnknownName);;
    give up;;
}`
	diags := analyzeSrc(t, src)
	if len(diags) == 0 {
		t.Fatal("expected diagnostics, got none")
	}
}

func TestAnalyze_DuplicateLocalDeclaration(t *testing.T) {
	src := `outy seq Foremost<> --> <<nada>> {
    allow exact Count 2b=2 1;;
    allow exact Count 2b=2 2;;
    give up;;
}`
	diags := analyzeSrc(t, src)
	if len(diags) == 0 {
		t.Fatal("expected diagnostics, got none")
	}
}

func analyzeSrc(t *testing.T, src string) []string {
	t.Helper()
	toks, lexDiags := lexer.Lex(src)
	if len(lexDiags) != 0 {
		t.Fatalf("unexpected lexer diagnostics: %d", len(lexDiags))
	}
	prog, parseDiags := parser.Parse(toks)
	if len(parseDiags) != 0 {
		t.Fatalf("unexpected parser diagnostics: %d", len(parseDiags))
	}
	out := sema.Analyze(prog)
	msgs := make([]string, 0, len(out))
	for _, d := range out {
		msgs = append(msgs, d.Message)
	}
	return msgs
}
