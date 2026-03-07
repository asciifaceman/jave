package parser_test

import (
	"os"
	"testing"

	"github.com/asciifaceman/jave/internal/ast"
	"github.com/asciifaceman/jave/internal/lexer"
	"github.com/asciifaceman/jave/internal/parser"
)

func TestParse_HelloWorldSequence(t *testing.T) {
	src := `outy seq Foremost<> --> <<nada>> {
    pront("hello, jave");;
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
	if len(prog.Sequences) != 1 {
		t.Fatalf("expected 1 sequence, got %d", len(prog.Sequences))
	}
	seq := prog.Sequences[0]
	if seq.Name != "Foremost" {
		t.Fatalf("expected Foremost, got %q", seq.Name)
	}
	if seq.ReturnType != "nada" {
		t.Fatalf("expected return type nada, got %q", seq.ReturnType)
	}
	if len(seq.Body) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(seq.Body))
	}
	if _, ok := seq.Body[0].(ast.ExprStmt); !ok {
		t.Fatalf("expected first statement to be ExprStmt")
	}
	if _, ok := seq.Body[1].(ast.GiveStmt); !ok {
		t.Fatalf("expected second statement to be GiveStmt")
	}
}

func TestParse_VarDeclAndGiveValue(t *testing.T) {
	src := `outy seq Add<> --> <<exact>> {
    allow exact A 2b=2 1;;
    allow exact B 2b=2 2;;
    give A + B up;;
}`
	toks, lexDiags := lexer.Lex(src)
	if len(lexDiags) != 0 {
		t.Fatalf("unexpected lexer diagnostics: %d", len(lexDiags))
	}

	prog, parseDiags := parser.Parse(toks)
	if len(parseDiags) != 0 {
		t.Fatalf("unexpected parser diagnostics: %d", len(parseDiags))
	}
	if len(prog.Sequences) != 1 {
		t.Fatalf("expected 1 sequence, got %d", len(prog.Sequences))
	}
	if len(prog.Sequences[0].Body) != 3 {
		t.Fatalf("expected 3 statements, got %d", len(prog.Sequences[0].Body))
	}
}

func TestParse_MissingStatementEndDiagnostic(t *testing.T) {
	src := `outy seq Foremost<> --> <<nada>> {
    pront("oops")
    give up;;
}`
	toks, _ := lexer.Lex(src)
	_, parseDiags := parser.Parse(toks)
	if len(parseDiags) == 0 {
		t.Fatal("expected parse diagnostics, got none")
	}
}

func TestParse_ConditionsExample(t *testing.T) {
	srcBytes, err := os.ReadFile("../../examples/conditions/main.jave")
	if err != nil {
		t.Fatalf("read example: %v", err)
	}
	toks, lexDiags := lexer.Lex(string(srcBytes))
	if len(lexDiags) != 0 {
		t.Fatalf("unexpected lexer diagnostics: %d", len(lexDiags))
	}

	prog, parseDiags := parser.Parse(toks)
	if len(parseDiags) != 0 {
		t.Fatalf("unexpected parser diagnostics: %d", len(parseDiags))
	}
	if len(prog.Sequences) != 1 {
		t.Fatalf("expected 1 sequence, got %d", len(prog.Sequences))
	}
	if len(prog.Sequences[0].Body) < 2 {
		t.Fatalf("expected condition sequence body to contain statements, got %d", len(prog.Sequences[0].Body))
	}
}

func TestParse_LoopsExample(t *testing.T) {
	srcBytes, err := os.ReadFile("../../examples/loops/main.jave")
	if err != nil {
		t.Fatalf("read example: %v", err)
	}
	toks, lexDiags := lexer.Lex(string(srcBytes))
	if len(lexDiags) != 0 {
		t.Fatalf("unexpected lexer diagnostics: %d", len(lexDiags))
	}

	prog, parseDiags := parser.Parse(toks)
	if len(parseDiags) != 0 {
		t.Fatalf("unexpected parser diagnostics: %d", len(parseDiags))
	}
	if len(prog.Imports) != 1 {
		t.Fatalf("expected 1 top-level import, got %d", len(prog.Imports))
	}
	if len(prog.Sequences) != 1 {
		t.Fatalf("expected 1 sequence, got %d", len(prog.Sequences))
	}
}

func TestParse_MultiDimensionalTablesExample(t *testing.T) {
	srcBytes, err := os.ReadFile("../../examples/multi_dimensional_tables/main.jave")
	if err != nil {
		t.Fatalf("read example: %v", err)
	}
	toks, lexDiags := lexer.Lex(string(srcBytes))
	if len(lexDiags) != 0 {
		t.Fatalf("unexpected lexer diagnostics: %d", len(lexDiags))
	}

	prog, parseDiags := parser.Parse(toks)
	if len(parseDiags) != 0 {
		t.Fatalf("unexpected parser diagnostics: %d", len(parseDiags))
	}
	if len(prog.Imports) != 1 {
		t.Fatalf("expected 1 import, got %d", len(prog.Imports))
	}
	if len(prog.Sequences) != 1 {
		t.Fatalf("expected 1 sequence, got %d", len(prog.Sequences))
	}
}
