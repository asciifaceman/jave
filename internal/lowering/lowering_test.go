package lowering_test

import (
	"os"
	"testing"

	"github.com/asciifaceman/jave/internal/lexer"
	"github.com/asciifaceman/jave/internal/lowering"
	"github.com/asciifaceman/jave/internal/parser"
	"github.com/asciifaceman/jave/internal/sema"
)

func TestLower_HelloWorldForemost(t *testing.T) {
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
		t.Fatalf("unexpected semantic diagnostics: %d", len(semaDiags))
	}

	irProg, lowerDiags := lowering.Lower(prog)
	if len(lowerDiags) != 0 {
		t.Fatalf("unexpected lowering diagnostics: %d", len(lowerDiags))
	}
	if irProg == nil {
		t.Fatal("expected lowered program, got nil")
	}
	if got := len(irProg.Foremost.Instructions); got != 2 {
		t.Fatalf("expected 2 instructions, got %d", got)
	}
}

func TestLower_ConditionsExample(t *testing.T) {
	src := `outy seq Foremost<> --> <<nada>> {
    allow vag Foo 2b=2 0.6;;

    maybe (<Foo bigly 0.5>) -> {
        Pront("Over half");;
    } furthermore (<Foo lessly 0.5>) -> {
        Pront("Under half");;
    } otherwise -> {
        Pront("Exactly half");;
    }

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
		t.Fatalf("unexpected semantic diagnostics: %d", len(semaDiags))
	}

	irProg, lowerDiags := lowering.Lower(prog)
	if len(lowerDiags) != 0 {
		t.Fatalf("unexpected lowering diagnostics: %d", len(lowerDiags))
	}
	if irProg == nil {
		t.Fatal("expected lowered program, got nil")
	}
	if got := len(irProg.Foremost.Instructions); got < 3 {
		t.Fatalf("expected lowered instructions to include condition flow, got %d", got)
	}
}

func TestLower_WhileGivenLoop(t *testing.T) {
	src := `outy seq Foremost<> --> <<nada>> {
    allow exact X 2b=2 0;;
    given (<X lesslysame 2>) again -> {
        X 2b=2 X + 1;;
    }
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
		t.Fatalf("unexpected semantic diagnostics: %d", len(semaDiags))
	}

	irProg, lowerDiags := lowering.Lower(prog)
	if len(lowerDiags) != 0 {
		t.Fatalf("unexpected lowering diagnostics: %d", len(lowerDiags))
	}
	if got := len(irProg.Foremost.Instructions); got < 3 {
		t.Fatalf("expected loop instructions to be lowered, got %d", got)
	}
}

func TestLower_LoopsExampleAllModes(t *testing.T) {
	srcBytes, err := os.ReadFile("../../examples/loops/main.jave")
	if err != nil {
		t.Fatalf("read loops example: %v", err)
	}
	toks, lexDiags := lexer.Lex(string(srcBytes))
	if len(lexDiags) != 0 {
		t.Fatalf("unexpected lexer diagnostics: %d", len(lexDiags))
	}
	prog, parseDiags := parser.Parse(toks)
	if len(parseDiags) != 0 {
		t.Fatalf("unexpected parser diagnostics: %d", len(parseDiags))
	}
	semaDiags := sema.Analyze(prog)
	if len(semaDiags) != 0 {
		t.Fatalf("unexpected semantic diagnostics: %d", len(semaDiags))
	}
	irProg, lowerDiags := lowering.Lower(prog)
	if len(lowerDiags) != 0 {
		t.Fatalf("unexpected lowering diagnostics: %d", len(lowerDiags))
	}
	if irProg == nil {
		t.Fatal("expected lowered program, got nil")
	}
}

func TestLower_ForewardAndForemost(t *testing.T) {
	src := `outy seq Foreward<> --> <<nada>> {
    Pront("warming carryon");;
    give up;;
}

outy seq Foremost<> --> <<nada>> {
    Pront("running foremost");;
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
		t.Fatalf("unexpected semantic diagnostics: %d", len(semaDiags))
	}

	irProg, lowerDiags := lowering.Lower(prog)
	if len(lowerDiags) != 0 {
		t.Fatalf("unexpected lowering diagnostics: %d", len(lowerDiags))
	}
	if len(irProg.Forewards) == 0 {
		t.Fatal("expected Foreward sequence to be lowered")
	}
	if irProg.Forewards[0].Name != "Foreward" {
		t.Fatalf("unexpected Foreward name: %q", irProg.Forewards[0].Name)
	}
}
