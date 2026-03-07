package runtime_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/asciifaceman/jave/internal/diagnostics"
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

func TestExecute_ForewardRunsBeforeForemost(t *testing.T) {
	srcBytes, err := os.ReadFile("../../examples/foreward_foremost/main.jave")
	if err != nil {
		t.Fatalf("read foreward example: %v", err)
	}
	buf := &bytes.Buffer{}
	if err := runSource(string(srcBytes), buf); err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	if got := strings.TrimSpace(buf.String()); got != "warming carryon\nrunning foremost" {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestExecute_ImportsExample(t *testing.T) {
	srcBytes, err := os.ReadFile("../../examples/imports/main.jave")
	if err != nil {
		t.Fatalf("read imports example: %v", err)
	}
	buf := &bytes.Buffer{}
	if err := runSource(string(srcBytes), buf); err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	got := strings.TrimSpace(buf.String())
	if !strings.Contains(got, "Imported systems online. Count=2") {
		t.Fatalf("expected prontulate output, got: %q", got)
	}
	if !strings.Contains(got, "Direct pront via Strangs too: 2") {
		t.Fatalf("expected combobulate output, got: %q", got)
	}
}

func TestExecute_SrangsLegacyAlias(t *testing.T) {
	src := `install Srangs from highschool/English;;
outy seq Foremost<> --> <<nada>> {
    pront(Srangs.Combobulate<"Legacy alias says: %exact", 7>);;
    give up;;
}`
	buf := &bytes.Buffer{}
	if err := runSource(src, buf); err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	if got := strings.TrimSpace(buf.String()); got != "Legacy alias says: 7" {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestExecute_SequenceParamsAndCalls(t *testing.T) {
	src := `outy seq Add<exact A, exact B> --> <<exact>> {
    give A + B up;;
}

outy seq Scale<exact Base, exact Multiplier> --> <<exact>> {
    give Base * Multiplier up;;
}

outy seq Foremost<> --> <<nada>> {
    allow exact Sum 2b=2 Add<7, 5>;;
    allow exact Product 2b=2 Scale<Sum, 3>;;
    pront(Product);;
    give up;;
}`

	buf := &bytes.Buffer{}
	if err := runSource(src, buf); err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	if got := strings.TrimSpace(buf.String()); got != "36" {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestExecute_PortfolioReviewExample(t *testing.T) {
	srcBytes, err := os.ReadFile("../../examples/portfolio_review/main.jave")
	if err != nil {
		t.Fatalf("read portfolio review example: %v", err)
	}
	buf := &bytes.Buffer{}
	if err := runSource(string(srcBytes), buf); err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	got := strings.TrimSpace(buf.String())
	if !strings.Contains(got, "Initiative=Payments") {
		t.Fatalf("expected initiative output, got: %q", got)
	}
	if !strings.Contains(got, "Best initiative: Mobile (63)") {
		t.Fatalf("expected best initiative summary, got: %q", got)
	}
	if !strings.Contains(got, "Average signal: 32.8") {
		t.Fatalf("expected average summary, got: %q", got)
	}
}

func TestExecute_IncidentTriageExample(t *testing.T) {
	srcBytes, err := os.ReadFile("../../examples/incident_triage/main.jave")
	if err != nil {
		t.Fatalf("read incident triage example: %v", err)
	}
	buf := &bytes.Buffer{}
	if err := runSource(string(srcBytes), buf); err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	got := strings.TrimSpace(buf.String())
	if !strings.Contains(got, "service=Auth score=92 lane=EXEC-WARROOM") {
		t.Fatalf("expected auth escalation output, got: %q", got)
	}
	if !strings.Contains(got, "Most critical service: Auth (92)") {
		t.Fatalf("expected critical service summary, got: %q", got)
	}
}

func TestExecute_BudgetPlanningExample(t *testing.T) {
	srcBytes, err := os.ReadFile("../../examples/budget_planning/main.jave")
	if err != nil {
		t.Fatalf("read budget planning example: %v", err)
	}
	buf := &bytes.Buffer{}
	if err := runSource(string(srcBytes), buf); err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	got := strings.TrimSpace(buf.String())
	if !strings.Contains(got, "Department=Platform annual=570") {
		t.Fatalf("expected department annual output, got: %q", got)
	}
	if !strings.Contains(got, "Annual portfolio budget=1797") {
		t.Fatalf("expected annual budget summary, got: %q", got)
	}
}

func TestExecute_ServiceCapacityPlanningExample(t *testing.T) {
	srcBytes, err := os.ReadFile("../../examples/service_capacity_planning/main.jave")
	if err != nil {
		t.Fatalf("read service capacity example: %v", err)
	}
	buf := &bytes.Buffer{}
	if err := runSource(string(srcBytes), buf); err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	got := strings.TrimSpace(buf.String())
	if !strings.Contains(got, "service=Checkout signal=114 lane=EXEC-SPONSOR") {
		t.Fatalf("expected checkout line, got: %q", got)
	}
	if !strings.Contains(got, "Top priority service: Checkout (114)") {
		t.Fatalf("expected top priority summary, got: %q", got)
	}
	if !strings.Contains(got, "Average portfolio signal: 72.8") {
		t.Fatalf("expected average summary, got: %q", got)
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
	for _, d := range semaDiags {
		if d.Severity == diagnostics.SeverityError {
			return semaErr(semaDiags)
		}
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
