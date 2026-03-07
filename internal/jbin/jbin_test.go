package jbin_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/asciifaceman/jave/internal/jbin"
	"github.com/asciifaceman/jave/internal/lexer"
	"github.com/asciifaceman/jave/internal/lowering"
	"github.com/asciifaceman/jave/internal/parser"
	"github.com/asciifaceman/jave/internal/sema"
)

func TestRoundTrip(t *testing.T) {
	src := `outy seq Foremost<> --> <<nada>> {
    Pront("hello, jave");;
    give up;;
}`
	toks, lexDiags := lexer.Lex(src)
	if len(lexDiags) != 0 {
		t.Fatalf("unexpected lexer diagnostics: %d", len(lexDiags))
	}
	prog, parseDiags := parser.Parse(toks)
	if len(parseDiags) != 0 {
		t.Fatalf("unexpected parser diagnostics: %d", len(parseDiags))
	}
	semaDiags := sema.Analyze(prog)
	if len(semaDiags) != 0 {
		t.Fatalf("unexpected sema diagnostics: %d", len(semaDiags))
	}
	irProg, lowerDiags := lowering.Lower(prog)
	if len(lowerDiags) != 0 {
		t.Fatalf("unexpected lower diagnostics: %d", len(lowerDiags))
	}

	path := filepath.Join(t.TempDir(), "hello.jbin")
	if err := jbin.WriteFile(path, irProg); err != nil {
		t.Fatalf("write jbin failed: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected jbin output to exist: %v", err)
	}
	decoded, err := jbin.ReadFile(path)
	if err != nil {
		t.Fatalf("read jbin failed: %v", err)
	}
	if decoded == nil {
		t.Fatal("expected decoded program, got nil")
	}
	if len(decoded.Foremost.Instructions) == 0 {
		t.Fatal("expected decoded instructions, got none")
	}
}
