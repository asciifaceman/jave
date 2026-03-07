package parser_test

import (
	"os"
	"testing"

	"github.com/asciifaceman/jave/internal/ast"
	"github.com/asciifaceman/jave/internal/diagnostics"
	"github.com/asciifaceman/jave/internal/lexer"
	"github.com/asciifaceman/jave/internal/parser"
)

func TestParse_HelloWorldSequence(t *testing.T) {
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

func TestParse_SequenceParams(t *testing.T) {
	src := `outy seq Add<exact A, exact B> --> <<exact>> {
    give A + B up;;
}

outy seq Foremost<> --> <<nada>> {
    allow exact Sum 2b=2 Add<2, 3>;;
    Pront(Sum);;
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
	if len(prog.Sequences) != 2 {
		t.Fatalf("expected 2 sequences, got %d", len(prog.Sequences))
	}
	if got := len(prog.Sequences[0].Params); got != 2 {
		t.Fatalf("expected 2 params on Add, got %d", got)
	}
	if prog.Sequences[0].Params[0].Name != "A" || prog.Sequences[0].Params[1].Name != "B" {
		t.Fatalf("unexpected param names: %+v", prog.Sequences[0].Params)
	}
}

func TestParse_SequenceParamsVariadicLast(t *testing.T) {
	src := `outy seq Combobulate<strang Template, ...strang Args> --> <<strang>> {
    give Template up;;
}

outy seq Foremost<> --> <<nada>> {
    Pront(Combobulate<"x=%exact", 1>);;
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
	if len(prog.Sequences) != 2 {
		t.Fatalf("expected 2 sequences, got %d", len(prog.Sequences))
	}
	params := prog.Sequences[0].Params
	if len(params) != 2 {
		t.Fatalf("expected 2 params, got %d", len(params))
	}
	if !params[1].Variadic {
		t.Fatal("expected last parameter to be variadic")
	}
}

func TestParse_SequenceParamsVariadicNotLastDiagnostic(t *testing.T) {
	src := `outy seq Bad<...strang Args, exact Tail> --> <<nada>> {
    give up;;
}`
	toks, lexDiags := lexer.Lex(src)
	if len(lexDiags) != 0 {
		t.Fatalf("unexpected lexer diagnostics: %d", len(lexDiags))
	}

	_, parseDiags := parser.Parse(toks)
	if len(parseDiags) == 0 {
		t.Fatal("expected parser diagnostics for non-final variadic parameter")
	}
}

func TestParse_DocstringAttachedToSequence(t *testing.T) {
	src := `doc<
Title:
    Positive exact value
About:
    Returns non-negative exact value.
Param:
    Value: Input exact to normalize.
>
outy seq PosiExact<exact Value> --> <<exact>> {
    give Value up;;
}`
	toks, lexDiags := lexer.Lex(src)
	if len(lexDiags) != 0 {
		t.Fatalf("unexpected lexer diagnostics: %d", len(lexDiags))
	}

	prog, parseDiags := parser.Parse(toks)
	for _, d := range parseDiags {
		if d.Severity == diagnostics.SeverityError {
			t.Fatalf("unexpected parser error: %s", d.Message)
		}
	}
	if len(prog.Sequences) != 1 {
		t.Fatalf("expected 1 sequence, got %d", len(prog.Sequences))
	}
	if prog.Sequences[0].Doc == nil {
		t.Fatal("expected docstring attached to sequence")
	}
	if len(prog.Sequences[0].Doc.Sections) < 2 {
		t.Fatalf("expected parsed sections, got %d", len(prog.Sequences[0].Doc.Sections))
	}
}

func TestParse_DocstringOrphanWarning(t *testing.T) {
	src := `doc<
Title:
    Orphaned docs
>
install Strangs from highschool/English;;

outy seq Foremost<> --> <<nada>> {
    give up;;
}`
	toks, lexDiags := lexer.Lex(src)
	if len(lexDiags) != 0 {
		t.Fatalf("unexpected lexer diagnostics: %d", len(lexDiags))
	}

	_, parseDiags := parser.Parse(toks)
	foundWarning := false
	for _, d := range parseDiags {
		if d.Severity == diagnostics.SeverityWarning && d.Message == "DOC-ATTACH: docstring was not attached to a sequence" {
			foundWarning = true
		}
	}
	if !foundWarning {
		t.Fatal("expected orphan docstring warning")
	}
}

func TestParse_DocstringAttachmentAcrossComments(t *testing.T) {
	src := `doc<
Title:
    Says hello
>
>>| comment between docstring and sequence
=[
block comment
]=
outy seq Hello<> --> <<nada>> {
    give up;;
}`
	toks, lexDiags := lexer.Lex(src)
	if len(lexDiags) != 0 {
		t.Fatalf("unexpected lexer diagnostics: %d", len(lexDiags))
	}
	prog, parseDiags := parser.Parse(toks)
	for _, d := range parseDiags {
		if d.Severity == diagnostics.SeverityError {
			t.Fatalf("unexpected parser error: %s", d.Message)
		}
	}
	if len(prog.Sequences) != 1 || prog.Sequences[0].Doc == nil {
		t.Fatal("expected docstring attachment to survive comments")
	}
}

func TestParse_MissingStatementEndDiagnostic(t *testing.T) {
	src := `outy seq Foremost<> --> <<nada>> {
    Pront("oops")
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

func TestParse_AdvLogAnomalyTriageExample(t *testing.T) {
	srcBytes, err := os.ReadFile("../../examples/adv-log-anomaly-triage/main.jave")
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

func TestParse_IncidentTriageExample(t *testing.T) {
	srcBytes, err := os.ReadFile("../../examples/incident_triage/main.jave")
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

func TestParse_AdvGameLobbyBalancerExample(t *testing.T) {
	srcBytes, err := os.ReadFile("../../examples/adv-game-lobby-balancer/main.jave")
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
	if len(prog.Imports) != 2 {
		t.Fatalf("expected 2 imports, got %d", len(prog.Imports))
	}
	if len(prog.Sequences) != 1 {
		t.Fatalf("expected 1 sequence, got %d", len(prog.Sequences))
	}
}

func TestParse_AdvMapSpawnSelectorExample(t *testing.T) {
	srcBytes, err := os.ReadFile("../../examples/adv-map-spawn-selector/main.jave")
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
	if len(prog.Imports) != 2 {
		t.Fatalf("expected 2 imports, got %d", len(prog.Imports))
	}
	if len(prog.Sequences) != 1 {
		t.Fatalf("expected 1 sequence, got %d", len(prog.Sequences))
	}
}
