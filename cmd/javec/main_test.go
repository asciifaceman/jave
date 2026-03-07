package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/asciifaceman/jave/internal/diagnostics"
	"github.com/asciifaceman/jave/internal/lowering"
	"github.com/asciifaceman/jave/internal/runtime"
	"github.com/asciifaceman/jave/internal/sema"
	"github.com/asciifaceman/jave/internal/sponsor"
)

func TestResolveSponsorMode(t *testing.T) {
	tests := []struct {
		name     string
		modeText string
		redacted bool
		quiet    bool
		want     sponsor.Mode
		wantErr  bool
	}{
		{name: "default full", modeText: "full", want: sponsor.ModeFull},
		{name: "explicit redacted", modeText: "redacted", want: sponsor.ModeRedacted},
		{name: "explicit off", modeText: "off", want: sponsor.ModeOff},
		{name: "alias redacted", modeText: "full", redacted: true, want: sponsor.ModeRedacted},
		{name: "alias quiet", modeText: "full", quiet: true, want: sponsor.ModeOff},
		{name: "invalid mode", modeText: "weird", wantErr: true},
		{name: "conflict quiet", modeText: "redacted", quiet: true, wantErr: true},
		{name: "conflict redacted", modeText: "off", redacted: true, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveSponsorMode(tt.modeText, tt.redacted, tt.quiet)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("resolveSponsorMode returned error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("resolveSponsorMode(%q, %v, %v) = %q, want %q", tt.modeText, tt.redacted, tt.quiet, got, tt.want)
			}
		})
	}
}

func TestLoadProgramWithImports_TransitiveForewardOrder(t *testing.T) {
	base := t.TempDir()

	innerPath := filepath.Join(base, "inner.jave")
	depPath := filepath.Join(base, "dep.jave")
	rootPath := filepath.Join(base, "main.jave")

	if err := os.WriteFile(innerPath, []byte(`outy seq Foreward<> --> <<nada>> {
    pront("inner foreward");;
    give up;;
}`), 0o644); err != nil {
		t.Fatalf("write inner: %v", err)
	}

	if err := os.WriteFile(depPath, []byte(`install Inner from ./inner;;
outy seq Foreward<> --> <<nada>> {
    pront("dep foreward");;
    give up;;
}`), 0o644); err != nil {
		t.Fatalf("write dep: %v", err)
	}

	if err := os.WriteFile(rootPath, []byte(`install Warmup from ./dep;;
outy seq Foremost<> --> <<nada>> {
    pront("root foremost");;
    give up;;
}`), 0o644); err != nil {
		t.Fatalf("write root: %v", err)
	}

	program, _, trace, loadDiags, _, err := loadProgramWithImports(rootPath, loadOptions{})
	if err != nil {
		t.Fatalf("loadProgramWithImports error: %v", err)
	}
	if len(loadDiags) != 0 {
		t.Fatalf("expected no load diagnostics, got %d", len(loadDiags))
	}
	if got := len(trace.Resolutions); got != 2 {
		t.Fatalf("expected 2 import resolutions, got %d", got)
	}

	semaDiags := sema.Analyze(program)
	for _, d := range semaDiags {
		if d.Severity == "error" {
			t.Fatalf("unexpected sema error: %s", d.Message)
		}
	}

	irProgram, lowerDiags := lowering.Lower(program)
	if len(lowerDiags) != 0 {
		t.Fatalf("unexpected lowering diagnostics: %d", len(lowerDiags))
	}
	if got := len(irProgram.Forewards); got != 2 {
		t.Fatalf("expected 2 forewards, got %d", got)
	}

	buf := &bytes.Buffer{}
	if err := runtime.Execute(irProgram, buf); err != nil {
		t.Fatalf("runtime execute failed: %v", err)
	}
	got := strings.TrimSpace(buf.String())
	want := "inner foreward\ndep foreward\nroot foremost"
	if got != want {
		t.Fatalf("unexpected output:\n%s\nwant:\n%s", got, want)
	}
}

func TestLoadProgramWithImports_HighschoolCarryon(t *testing.T) {
	base := t.TempDir()

	if err := os.WriteFile(filepath.Join(base, "go.mod"), []byte("module example.com/test\n\ngo 1.26\n"), 0o644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	englishDir := filepath.Join(base, "highschool", "English")
	if err := os.MkdirAll(englishDir, 0o755); err != nil {
		t.Fatalf("mkdir english dir: %v", err)
	}
	englishPath := filepath.Join(englishDir, "main.jave")
	if err := os.WriteFile(englishPath, []byte(`outy seq Foreward<> --> <<nada>> {
    pront("highschool english foreward");;
    give up;;
}`), 0o644); err != nil {
		t.Fatalf("write english carryon: %v", err)
	}

	appDir := filepath.Join(base, "apps", "demo")
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		t.Fatalf("mkdir app dir: %v", err)
	}
	rootPath := filepath.Join(appDir, "main.jave")
	if err := os.WriteFile(rootPath, []byte(`install Strangs from highschool/English;;
outy seq Foremost<> --> <<nada>> {
    pront("root foremost");;
    give up;;
}`), 0o644); err != nil {
		t.Fatalf("write root: %v", err)
	}

	program, _, trace, loadDiags, _, err := loadProgramWithImports(rootPath, loadOptions{})
	if err != nil {
		t.Fatalf("loadProgramWithImports error: %v", err)
	}
	if len(loadDiags) != 0 {
		t.Fatalf("expected no load diagnostics, got %d", len(loadDiags))
	}
	if got := len(trace.Resolutions); got != 1 {
		t.Fatalf("expected 1 import resolution, got %d", got)
	}

	semaDiags := sema.Analyze(program)
	for _, d := range semaDiags {
		if d.Severity == "error" {
			t.Fatalf("unexpected sema error: %s", d.Message)
		}
	}

	irProgram, lowerDiags := lowering.Lower(program)
	if len(lowerDiags) != 0 {
		t.Fatalf("unexpected lowering diagnostics: %d", len(lowerDiags))
	}
	if got := len(irProgram.Forewards); got != 1 {
		t.Fatalf("expected 1 foreward, got %d", got)
	}

	buf := &bytes.Buffer{}
	if err := runtime.Execute(irProgram, buf); err != nil {
		t.Fatalf("runtime execute failed: %v", err)
	}
	got := strings.TrimSpace(buf.String())
	want := "highschool english foreward\nroot foremost"
	if got != want {
		t.Fatalf("unexpected output:\n%s\nwant:\n%s", got, want)
	}
}

func TestLoadProgramWithImports_HighschoolAlgebraPosiExactAndPosiVag(t *testing.T) {
	base := t.TempDir()

	if err := os.WriteFile(filepath.Join(base, "go.mod"), []byte("module example.com/test\n\ngo 1.26\n"), 0o644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	algebraDir := filepath.Join(base, "highschool", "Algebra")
	if err := os.MkdirAll(algebraDir, 0o755); err != nil {
		t.Fatalf("mkdir algebra dir: %v", err)
	}
	algebraPath := filepath.Join(algebraDir, "main.jave")
	if err := os.WriteFile(algebraPath, []byte(`outy seq PosiExact<exact Value> --> <<exact>> {
	maybe (<Value lessly 0>) -> {
		give 0 - Value up;;
	}
	give Value up;;
}

outy seq PosiVag<vag Value> --> <<vag>> {
    maybe (<Value lessly 0>) -> {
        give 0 - Value up;;
    }
    give Value up;;
}`), 0o644); err != nil {
		t.Fatalf("write algebra carryon: %v", err)
	}

	appDir := filepath.Join(base, "apps", "demo")
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		t.Fatalf("mkdir app dir: %v", err)
	}
	rootPath := filepath.Join(appDir, "main.jave")
	if err := os.WriteFile(rootPath, []byte(`install Algebra from highschool/Algebra;;
outy seq Foremost<> --> <<nada>> {
	pront(Algebra.PosiExact<0 - 12>);;
    pront(Algebra.PosiExact<7>);;
	pront(Algebra.PosiVag<0.5 - 1.25>);;
    give up;;
}`), 0o644); err != nil {
		t.Fatalf("write root: %v", err)
	}

	program, _, trace, loadDiags, _, err := loadProgramWithImports(rootPath, loadOptions{})
	if err != nil {
		t.Fatalf("loadProgramWithImports error: %v", err)
	}
	if len(loadDiags) != 0 {
		t.Fatalf("expected no load diagnostics, got %d", len(loadDiags))
	}
	if got := len(trace.Resolutions); got != 1 {
		t.Fatalf("expected 1 import resolution, got %d", got)
	}

	semaDiags := sema.Analyze(program)
	for _, d := range semaDiags {
		if d.Severity == "error" {
			t.Fatalf("unexpected sema error: %s", d.Message)
		}
	}

	irProgram, lowerDiags := lowering.Lower(program)
	if len(lowerDiags) != 0 {
		t.Fatalf("unexpected lowering diagnostics: %d", len(lowerDiags))
	}

	buf := &bytes.Buffer{}
	if err := runtime.Execute(irProgram, buf); err != nil {
		t.Fatalf("runtime execute failed: %v", err)
	}
	got := strings.TrimSpace(buf.String())
	want := "12\n7\n0.75"
	if got != want {
		t.Fatalf("unexpected output:\n%s\nwant:\n%s", got, want)
	}
}

func TestLoadProgramWithImports_SemaDiagnosticResolvesImportedFile(t *testing.T) {
	base := t.TempDir()

	depPath := filepath.Join(base, "dep.jave")
	rootPath := filepath.Join(base, "main.jave")

	if err := os.WriteFile(depPath, []byte(`outy seq Foreward<> --> <<nada>> {
    allow exact Count 2b=2 1;;
    allow exact Count 2b=2 2;;
    give up;;
}`), 0o644); err != nil {
		t.Fatalf("write dep: %v", err)
	}

	if err := os.WriteFile(rootPath, []byte(`install Dep from ./dep;;
outy seq Foremost<> --> <<nada>> {
    give up;;
}`), 0o644); err != nil {
		t.Fatalf("write root: %v", err)
	}

	program, _, _, _, sourceMap, err := loadProgramWithImports(rootPath, loadOptions{})
	if err != nil {
		t.Fatalf("loadProgramWithImports error: %v", err)
	}

	diags := sema.Analyze(program)
	if len(diags) == 0 {
		t.Fatal("expected sema diagnostics")
	}

	var target diagnostics.Diagnostic
	found := false
	for _, d := range diags {
		if d.Severity == diagnostics.SeverityError && strings.Contains(d.Message, "duplicate local declaration") {
			target = d
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected duplicate local declaration diagnostic, got: %+v", diags)
	}

	diagPath, diagLine, _ := resolveDiagnosticLocation(target.Pos, rootPath, sourceMap)
	if filepath.Clean(diagPath) != filepath.Clean(depPath) {
		t.Fatalf("diagnostic path = %q, want %q", diagPath, depPath)
	}
	if diagLine != 3 {
		t.Fatalf("diagnostic line = %d, want 3", diagLine)
	}
}

func TestLoadProgramWithImports_ImportCycleDetected(t *testing.T) {
	base := t.TempDir()

	rootPath := filepath.Join(base, "main.jave")
	aPath := filepath.Join(base, "a.jave")
	bPath := filepath.Join(base, "b.jave")

	if err := os.WriteFile(rootPath, []byte(`install A from ./a;;
outy seq Foremost<> --> <<nada>> {
    give up;;
}`), 0o644); err != nil {
		t.Fatalf("write root: %v", err)
	}

	if err := os.WriteFile(aPath, []byte(`install B from ./b;;
outy seq Foreward<> --> <<nada>> {
    give up;;
}`), 0o644); err != nil {
		t.Fatalf("write a: %v", err)
	}

	if err := os.WriteFile(bPath, []byte(`install A from ./a;;
outy seq Foreward<> --> <<nada>> {
    give up;;
}`), 0o644); err != nil {
		t.Fatalf("write b: %v", err)
	}

	_, _, _, _, _, err := loadProgramWithImports(rootPath, loadOptions{})
	if err == nil {
		t.Fatal("expected import cycle error")
	}
	msg := err.Error()
	if !strings.Contains(msg, "import cycle detected") {
		t.Fatalf("unexpected cycle error: %q", msg)
	}
	if !strings.Contains(msg, filepath.Clean(aPath)) || !strings.Contains(msg, filepath.Clean(bPath)) {
		t.Fatalf("expected cycle path to include a and b, got: %q", msg)
	}
}

func TestLoadProgramWithImports_SelfImportCycleDetected(t *testing.T) {
	base := t.TempDir()

	rootPath := filepath.Join(base, "main.jave")
	if err := os.WriteFile(rootPath, []byte(`install Self from ./main;;
outy seq Foremost<> --> <<nada>> {
    give up;;
}`), 0o644); err != nil {
		t.Fatalf("write root: %v", err)
	}

	_, _, _, _, _, err := loadProgramWithImports(rootPath, loadOptions{})
	if err == nil {
		t.Fatal("expected self import cycle error")
	}
	msg := err.Error()
	if !strings.Contains(msg, "import cycle detected") {
		t.Fatalf("unexpected cycle error: %q", msg)
	}
	if !strings.Contains(msg, filepath.Clean(rootPath)) {
		t.Fatalf("expected cycle to include root path, got: %q", msg)
	}
}

func TestFindWorkspaceRoot(t *testing.T) {
	base := t.TempDir()
	if err := os.WriteFile(filepath.Join(base, "go.mod"), []byte("module example.com/test\n\ngo 1.26\n"), 0o644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	nested := filepath.Join(base, "a", "b", "main.jave")
	if err := os.MkdirAll(filepath.Dir(nested), 0o755); err != nil {
		t.Fatalf("mkdir nested dir: %v", err)
	}
	if err := os.WriteFile(nested, []byte(""), 0o644); err != nil {
		t.Fatalf("write nested file: %v", err)
	}

	got := findWorkspaceRoot(nested)
	if filepath.Clean(got) != filepath.Clean(base) {
		t.Fatalf("workspace root = %q, want %q", got, base)
	}
}

func TestResolveHighschoolImportPath_MissingIncludesTriedPaths(t *testing.T) {
	base := t.TempDir()
	if err := os.WriteFile(filepath.Join(base, "go.mod"), []byte("module example.com/test\n\ngo 1.26\n"), 0o644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	entry := filepath.Join(base, "apps", "demo", "main.jave")
	if err := os.MkdirAll(filepath.Dir(entry), 0o755); err != nil {
		t.Fatalf("mkdir app dir: %v", err)
	}
	if err := os.WriteFile(entry, []byte(""), 0o644); err != nil {
		t.Fatalf("write entry file: %v", err)
	}

	_, err := resolveHighschoolImportPath(entry, "highschool/History", "")
	if err == nil {
		t.Fatal("expected missing highschool import error")
	}
	msg := err.Error()
	if !strings.Contains(msg, "tried:") {
		t.Fatalf("expected tried paths in error, got: %q", msg)
	}
	if !strings.Contains(msg, filepath.Join("highschool", "History", "main.jave")) {
		t.Fatalf("expected highschool candidate path in error, got: %q", msg)
	}
}

func TestResolveHighschoolImportPath_Ambiguous(t *testing.T) {
	base := t.TempDir()
	if err := os.WriteFile(filepath.Join(base, "go.mod"), []byte("module example.com/test\n\ngo 1.26\n"), 0o644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	entry := filepath.Join(base, "apps", "demo", "main.jave")
	if err := os.MkdirAll(filepath.Dir(entry), 0o755); err != nil {
		t.Fatalf("mkdir app dir: %v", err)
	}
	if err := os.WriteFile(entry, []byte(""), 0o644); err != nil {
		t.Fatalf("write entry file: %v", err)
	}

	pathA := filepath.Join(base, "highschool", "English", "main.jave")
	pathB := filepath.Join(base, "stdlib", "highschool", "English", "main.jave")
	if err := os.MkdirAll(filepath.Dir(pathA), 0o755); err != nil {
		t.Fatalf("mkdir pathA: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(pathB), 0o755); err != nil {
		t.Fatalf("mkdir pathB: %v", err)
	}
	if err := os.WriteFile(pathA, []byte(""), 0o644); err != nil {
		t.Fatalf("write pathA: %v", err)
	}
	if err := os.WriteFile(pathB, []byte(""), 0o644); err != nil {
		t.Fatalf("write pathB: %v", err)
	}

	_, err := resolveHighschoolImportPath(entry, "highschool/English", "")
	if err == nil {
		t.Fatal("expected ambiguous highschool import error")
	}
	msg := err.Error()
	if !strings.Contains(msg, "ambiguous highschool import") {
		t.Fatalf("unexpected error message: %q", msg)
	}
	if !strings.Contains(msg, filepath.Clean(pathA)) || !strings.Contains(msg, filepath.Clean(pathB)) {
		t.Fatalf("expected both matching paths in error, got: %q", msg)
	}
}

func TestPrintImportTrace(t *testing.T) {
	trace := importTrace{
		LoadOrder: []string{"a.jave", "b.jave"},
		Resolutions: []importResolution{
			{Importer: "a.jave", From: "./b", Resolved: "b.jave"},
		},
	}

	buf := &bytes.Buffer{}
	printImportTrace(buf, trace)

	out := buf.String()
	if !strings.Contains(out, "javec: import trace") {
		t.Fatalf("missing header in output: %q", out)
	}
	if !strings.Contains(out, "load[0]: a.jave") {
		t.Fatalf("missing load row in output: %q", out)
	}
	if !strings.Contains(out, "resolve: a.jave -> b.jave (./b)") {
		t.Fatalf("missing resolve row in output: %q", out)
	}
}

func TestResolveHighschoolImportPath_ProjectRootOverride(t *testing.T) {
	root := t.TempDir()
	projectRoot := filepath.Join(root, "my-project")
	outsideRoot := filepath.Join(root, "outside")

	if err := os.MkdirAll(filepath.Join(projectRoot, "highschool", "English"), 0o755); err != nil {
		t.Fatalf("mkdir project highschool dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(projectRoot, "highschool", "English", "main.jave"), []byte(""), 0o644); err != nil {
		t.Fatalf("write project highschool file: %v", err)
	}

	entry := filepath.Join(outsideRoot, "app", "main.jave")
	if err := os.MkdirAll(filepath.Dir(entry), 0o755); err != nil {
		t.Fatalf("mkdir outside app dir: %v", err)
	}
	if err := os.WriteFile(entry, []byte(""), 0o644); err != nil {
		t.Fatalf("write outside entry file: %v", err)
	}

	resolved, err := resolveHighschoolImportPath(entry, "highschool/English", projectRoot)
	if err != nil {
		t.Fatalf("resolveHighschoolImportPath returned error: %v", err)
	}
	want := filepath.Join(projectRoot, "highschool", "English", "main.jave")
	if filepath.Clean(resolved) != filepath.Clean(want) {
		t.Fatalf("resolved = %q, want %q", resolved, want)
	}
}
