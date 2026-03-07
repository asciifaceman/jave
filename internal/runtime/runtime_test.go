package runtime_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/asciifaceman/jave/internal/lexer"
	"github.com/asciifaceman/jave/internal/lowering"
	"github.com/asciifaceman/jave/internal/parser"
	"github.com/asciifaceman/jave/internal/runtime"
	"github.com/asciifaceman/jave/internal/sema"
)

func TestExecute_PrintsHelloWorld(t *testing.T) {
	src := `outy seq Foremost<> --> <<nada>> {
    pront("hello, jave");;
    give up;;
}`
	buf := &bytes.Buffer{}
	if err := runSource(src, buf); err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	if got := strings.TrimSpace(buf.String()); got != "hello, jave" {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestExecute_CombobulateAndGirth(t *testing.T) {
	src := `install Strangs from highschool/English;;
outy seq Foremost<> --> <<nada>> {
    allow table<exact> Scores 2b=2 [1, 2, 3];;
    pront(Strangs.Combobulate<"Scores girth: %exact", girth(Scores)>);;
    give up;;
}`
	buf := &bytes.Buffer{}
	if err := runSource(src, buf); err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	if got := strings.TrimSpace(buf.String()); got != "Scores girth: 3" {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestExecute_ConditionsBranching(t *testing.T) {
	src := `outy seq Foremost<> --> <<nada>> {
    allow vag Foo 2b=2 0.6;;

    maybe (<Foo bigly 0.5>) -> {
        pront("Over half");;
    } furthermore (<Foo lessly 0.5>) -> {
        pront("Under half");;
    } otherwise -> {
        pront("Exactly half");;
    }

    give up;;
}`
	buf := &bytes.Buffer{}
	if err := runSource(src, buf); err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	if got := strings.TrimSpace(buf.String()); got != "Over half" {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestExecute_WhileGivenLoop(t *testing.T) {
	src := `outy seq Foremost<> --> <<nada>> {
    allow exact X 2b=2 0;;
    given (<X lesslysame 3>) again -> {
        pront(X);;
        X 2b=2 X + 1;;
    }
    give up;;
}`
	buf := &bytes.Buffer{}
	if err := runSource(src, buf); err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	if got := strings.TrimSpace(buf.String()); got != "0\n1\n2\n3" {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestExecute_LoopsExample(t *testing.T) {
	srcBytes, err := os.ReadFile("../../examples/loops/main.jave")
	if err != nil {
		t.Fatalf("read loops example: %v", err)
	}
	buf := &bytes.Buffer{}
	if err := runSource(string(srcBytes), buf); err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	got := strings.TrimSpace(buf.String())
	if !strings.Contains(got, "while-ish X: 0") {
		t.Fatalf("expected while-ish output, got: %q", got)
	}
	if !strings.Contains(got, "for-ish I: 2") {
		t.Fatalf("expected for-ish output, got: %q", got)
	}
	if !strings.Contains(got, "Grace") {
		t.Fatalf("expected within output, got: %q", got)
	}
}

func runSource(src string, out *bytes.Buffer) error {
	toks, lexDiags := lexer.Lex(src)
	if len(lexDiags) != 0 {
		return lexErr(lexDiags)
	}
	prog, parseDiags := parser.Parse(toks)
	if len(parseDiags) != 0 {
		return parseErr(parseDiags)
	}
	semaDiags := sema.Analyze(prog)
	if len(semaDiags) != 0 {
		return semaErr(semaDiags)
	}
	irProg, lowerDiags := lowering.Lower(prog)
	if len(lowerDiags) != 0 {
		return lowerErr(lowerDiags)
	}
	return runtime.Execute(irProg, out)
}

func lexErr(v any) error   { return &testErr{msg: "lexer diagnostics"} }
func parseErr(v any) error { return &testErr{msg: "parser diagnostics"} }
func semaErr(v any) error  { return &testErr{msg: "semantic diagnostics"} }
func lowerErr(v any) error { return &testErr{msg: "lowering diagnostics"} }

type testErr struct{ msg string }

func (e *testErr) Error() string { return e.msg }
