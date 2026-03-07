package sema_test

import (
	"strings"
	"testing"

	"github.com/asciifaceman/jave/internal/ast"
	"github.com/asciifaceman/jave/internal/diagnostics"
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

func TestAnalyze_SrangsLegacyAliasWarning(t *testing.T) {
	src := `install Srangs from highschool/English;;
outy seq Foremost<> --> <<nada>> {
    pront(Srangs.Combobulate<"Legacy %exact", 1>);;
    give up;;
}`
	diags := analyzeSrcDetailed(t, src)
	found := false
	for _, d := range diags {
		if d.Severity == diagnostics.SeverityWarning && strings.Contains(d.Message, "legacy module alias 'Srangs'") {
			found = true
		}
		if d.Severity == diagnostics.SeverityError {
			t.Fatalf("unexpected error diagnostic: %s", d.Message)
		}
	}
	if !found {
		t.Fatal("expected legacy Srangs warning")
	}
}

func TestAnalyze_DuplicateImports(t *testing.T) {
	src := `install Strangs from highschool/English;;
install Strangs from highschool/English;;
outy seq Foremost<> --> <<nada>> {
    give up;;
}`
	diags := analyzeSrcDetailed(t, src)
	found := false
	for _, d := range diags {
		if d.Severity == diagnostics.SeverityError && strings.Contains(d.Message, "duplicate import declaration") {
			found = true
		}
	}
	if !found {
		t.Fatal("expected duplicate import diagnostic")
	}
}

func TestAnalyze_StdlibImportPathValidation(t *testing.T) {
	src := `install Strangs from corp/English;;
outy seq Foremost<> --> <<nada>> {
    give up;;
}`
	diags := analyzeSrcDetailed(t, src)
	found := false
	for _, d := range diags {
		if d.Severity == diagnostics.SeverityError && strings.Contains(d.Message, "must use highschool/... path") {
			found = true
		}
	}
	if !found {
		t.Fatal("expected stdlib import path diagnostic")
	}
}

func TestAnalyze_SequenceParamsAreInScope(t *testing.T) {
	src := `outy seq Add<exact A, exact B> --> <<exact>> {
    give A + B up;;
}

outy seq Foremost<> --> <<nada>> {
    allow exact Sum 2b=2 Add<2, 3>;;
    pront(Sum);;
    give up;;
}`

	diags := analyzeSrcDetailed(t, src)
	for _, d := range diags {
		if d.Severity == diagnostics.SeverityError {
			t.Fatalf("unexpected semantic error: %s", d.Message)
		}
	}
}

func TestAnalyze_SequenceCallArityMismatch(t *testing.T) {
	src := `outy seq Add<exact A, exact B> --> <<exact>> {
    give A + B up;;
}
outy seq Foremost<> --> <<nada>> {
    allow exact Sum 2b=2 Add<2>;;
    pront(Sum);;
    give up;;
}`

	diags := analyzeSrcDetailed(t, src)
	found := false
	for _, d := range diags {
		if d.Severity == diagnostics.SeverityError && strings.Contains(d.Message, "arity mismatch") {
			found = true
		}
	}
	if !found {
		t.Fatal("expected sequence call arity mismatch diagnostic")
	}
}

func TestAnalyze_ProntulateBuiltinIdentifier(t *testing.T) {
	src := `outy seq Foremost<> --> <<nada>> {
    Prontulate<"Count=%exact", 2>;;
    give up;;
}`

	diags := analyzeSrcDetailed(t, src)
	for _, d := range diags {
		if d.Severity == diagnostics.SeverityError {
			t.Fatalf("unexpected semantic error: %s", d.Message)
		}
	}
}

func TestAnalyze_ModuleMemberCallArityMismatch_PosiExact(t *testing.T) {
	program := &ast.Program{
		Imports: []ast.ImportDecl{{Name: "Algebra", From: "highschool/Algebra"}},
		Sequences: []ast.SequenceDecl{
			{
				SourceModule: "Algebra",
				Name:         "PosiExact",
				Params:       []ast.SequenceParam{{TypeName: "exact", Name: "Value"}},
				ReturnType:   "exact",
				Body:         []ast.Stmt{ast.GiveStmt{Value: ast.NumberExpr{Value: "1"}}},
			},
			{
				Name:       "Foremost",
				ReturnType: "nada",
				Body: []ast.Stmt{
					ast.ExprStmt{Expr: ast.CallExpr{Callee: ast.MemberExpr{Target: ast.IdentifierExpr{Name: "Algebra"}, Name: "PosiExact"}, Args: []ast.Expr{}}},
					ast.GiveStmt{},
				},
			},
		},
	}

	diags := sema.Analyze(program)
	found := false
	for _, d := range diags {
		if d.Severity == diagnostics.SeverityError && strings.Contains(d.Message, "sequence call arity mismatch for Algebra.PosiExact") {
			found = true
		}
	}
	if !found {
		t.Fatal("expected module member arity mismatch diagnostic")
	}
}

func TestAnalyze_ModuleMemberMissingDiagnostic(t *testing.T) {
	program := &ast.Program{
		Imports: []ast.ImportDecl{{Name: "Algebra", From: "highschool/Algebra"}},
		Sequences: []ast.SequenceDecl{
			{
				SourceModule: "Algebra",
				Name:         "PosiExact",
				Params:       []ast.SequenceParam{{TypeName: "exact", Name: "Value"}},
				ReturnType:   "exact",
				Body:         []ast.Stmt{ast.GiveStmt{Value: ast.NumberExpr{Value: "1"}}},
			},
			{
				Name:       "Foremost",
				ReturnType: "nada",
				Body: []ast.Stmt{
					ast.ExprStmt{Expr: ast.CallExpr{Callee: ast.MemberExpr{Target: ast.IdentifierExpr{Name: "Algebra"}, Name: "Nope"}, Args: []ast.Expr{ast.NumberExpr{Value: "1"}}}},
					ast.GiveStmt{},
				},
			},
		},
	}

	diags := sema.Analyze(program)
	found := false
	for _, d := range diags {
		if d.Severity == diagnostics.SeverityError && strings.Contains(d.Message, "undefined module sequence: Algebra.Nope") {
			found = true
		}
	}
	if !found {
		t.Fatal("expected undefined module sequence diagnostic")
	}
}

func analyzeSrc(t *testing.T, src string) []string {
	t.Helper()
	msgs := make([]string, 0)
	for _, d := range analyzeSrcDetailed(t, src) {
		msgs = append(msgs, d.Message)
	}
	return msgs
}

func analyzeSrcDetailed(t *testing.T, src string) []diagnostics.Diagnostic {
	t.Helper()
	toks, lexDiags := lexer.Lex(src)
	if len(lexDiags) != 0 {
		t.Fatalf("unexpected lexer diagnostics: %d", len(lexDiags))
	}
	prog, parseDiags := parser.Parse(toks)
	if len(parseDiags) != 0 {
		t.Fatalf("unexpected parser diagnostics: %d", len(parseDiags))
	}
	return sema.Analyze(prog)
}
