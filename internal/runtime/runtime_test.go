package runtime_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/asciifaceman/jave/internal/ast"
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

func TestExecute_LexisIndexing(t *testing.T) {
	src := `outy seq Foremost<> --> <<nada>> {
    allow lexis<strang, exact> Ages 2b=2 {"Ada": 36, "Linus": 55};;
    pront(Ages["Ada"]);;
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

func TestExecute_LexisMissingKeyErrors(t *testing.T) {
	src := `outy seq Foremost<> --> <<nada>> {
    allow lexis<strang, exact> Ages 2b=2 {"Ada": 36};;
    pront(Ages["Grace"]);;
    give up;;
}`
	buf := &bytes.Buffer{}
	err := runSource(src, buf)
	if err == nil {
		t.Fatal("expected runtime error for missing lexis key")
	}
	if !strings.Contains(err.Error(), "lexis key not found") {
		t.Fatalf("expected lexis key not found error, got: %v", err)
	}
}

func TestExecute_ProntulateBuiltin(t *testing.T) {
	src := `outy seq Foremost<> --> <<nada>> {
	prontulate<"Count=%exact", 2>;;
    give up;;
}`
	buf := &bytes.Buffer{}
	if err := runSource(src, buf); err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	if got := strings.TrimSpace(buf.String()); got != "Count=2" {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestExecute_SlotifyBuiltin(t *testing.T) {
	src := `outy seq Foremost<> --> <<nada>> {
    pront(slotify("A=%exact B=%strang", 7));;
    pront(slotify(slotify("A=%exact B=%strang", 7), "yee"));;
    give up;;
}`
	buf := &bytes.Buffer{}
	if err := runSource(src, buf); err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	got := strings.TrimSpace(buf.String())
	want := "A=7 B=%strang\nA=7 B=yee"
	if got != want {
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

func TestExecute_CollectionsExample(t *testing.T) {
	srcBytes, err := os.ReadFile("../../examples/collections/main.jave")
	if err != nil {
		t.Fatalf("read collections example: %v", err)
	}
	buf := &bytes.Buffer{}
	if err := runSource(string(srcBytes), buf); err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	got := strings.TrimSpace(buf.String())
	if !strings.Contains(got, "Scores girth: 3") {
		t.Fatalf("expected scores output, got: %q", got)
	}
	if !strings.Contains(got, "Ada age: 36") {
		t.Fatalf("expected lexis output, got: %q", got)
	}
}

func TestExecute_MathAlgebraOpsExample(t *testing.T) {
	src := `outy seq MostExact<exact Left, exact Right> --> <<exact>> {
	maybe (<Left biglysame Right>) -> {
		give Left up;;
	}
	give Right up;;
}

outy seq LeastExact<exact Left, exact Right> --> <<exact>> {
	maybe (<Left lesslysame Right>) -> {
		give Left up;;
	}
	give Right up;;
}

outy seq MostVag<vag Left, vag Right> --> <<vag>> {
	maybe (<Left biglysame Right>) -> {
		give Left up;;
	}
	give Right up;;
}

outy seq LeastVag<vag Left, vag Right> --> <<vag>> {
	maybe (<Left lesslysame Right>) -> {
		give Left up;;
	}
	give Right up;;
}

outy seq PosidirExact<exact Value> --> <<exact>> {
	maybe (<Value bigly 0>) -> {
		give 1 up;;
	} furthermore (<Value lessly 0>) -> {
		give 0 - 1 up;;
	}
	give 0 up;;
}

outy seq PosidirVag<vag Value> --> <<exact>> {
	maybe (<Value bigly 0>) -> {
		give 1 up;;
	} furthermore (<Value lessly 0>) -> {
		give 0 - 1 up;;
	}
	give 0 up;;
}

outy seq Nearlydont<vag Value> --> <<truther>> {
	allow vag Epsilon 2b=2 0.000001;;
	maybe (<Value lessly 0>) -> {
		give (0 - Value) lesslysame Epsilon up;;
	}
	give Value lesslysame Epsilon up;;
}

outy seq Stretch<vag Start, vag End, vag Progress> --> <<vag>> {
	give Start + ((End - Start) * Progress) up;;
}

outy seq Foremost<> --> <<nada>> {
	pront(MostExact<3, 9>);;
	pront(LeastExact<3, 9>);;
	pront(MostVag<1.25, 1.2>);;
	pront(LeastVag<1.25, 1.2>);;

	pront(PosidirExact<0 - 5>);;
	pront(PosidirExact<0>);;
	pront(PosidirVag<2.5>);;

	pront(Nearlydont<0.0000004>);;
	pront(Nearlydont<0.01>);;

	pront(Stretch<10.0, 20.0, 0.25>);;
	give up;;
}`
	buf := &bytes.Buffer{}
	if err := runSource(src, buf); err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	got := strings.TrimSpace(buf.String())
	want := "9\n3\n1.25\n1.2\n-1\n0\n1\nyee\nnee\n12.5"
	if got != want {
		t.Fatalf("unexpected output:\n%s\nwant:\n%s", got, want)
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
	if err := appendImportedSequences(prog); err != nil {
		return err
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

func appendImportedSequences(prog *ast.Program) error {
	for _, imp := range prog.Imports {
		module, err := loadModuleProgram(imp.From)
		if err != nil {
			return err
		}
		for _, seq := range module.Sequences {
			seq.SourceModule = imp.Name
			prog.Sequences = append(prog.Sequences, seq)
		}
	}
	return nil
}

func loadModuleProgram(from string) (*ast.Program, error) {
	modulePath, err := resolveModuleMainFile(from)
	if err != nil {
		return nil, err
	}
	srcBytes, err := os.ReadFile(modulePath)
	if err != nil {
		return nil, fmt.Errorf("read module %s: %w", modulePath, err)
	}
	toks, lexDiags := lexer.Lex(string(srcBytes))
	if len(lexDiags) != 0 {
		return nil, fmt.Errorf("lexer diagnostics while loading module %s", from)
	}
	module, parseDiags := parser.Parse(toks)
	if len(parseDiags) != 0 {
		return nil, fmt.Errorf("parser diagnostics while loading module %s", from)
	}
	return module, nil
}

func resolveModuleMainFile(from string) (string, error) {
	relModulePath := filepath.Join(filepath.FromSlash(from), "main.jave")
	candidates := []string{
		relModulePath,
		filepath.Join("..", "..", relModulePath),
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("module source for %s not found", from)
}

func lexErr(v any) error   { return &testErr{msg: "lexer diagnostics"} }
func parseErr(v any) error { return &testErr{msg: "parser diagnostics"} }
func semaErr(v any) error  { return &testErr{msg: "semantic diagnostics"} }
func lowerErr(v any) error { return &testErr{msg: "lowering diagnostics"} }

type testErr struct{ msg string }

func (e *testErr) Error() string { return e.msg }
