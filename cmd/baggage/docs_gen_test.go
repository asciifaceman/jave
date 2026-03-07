package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseDocsArgsDefaults(t *testing.T) {
	base := t.TempDir()
	opts, err := parseDocsArgs([]string{"--project-root", base})
	if err != nil {
		t.Fatalf("parseDocsArgs returned error: %v", err)
	}
	if opts.projectRoot != base {
		t.Fatalf("projectRoot = %q", opts.projectRoot)
	}
	if opts.outDir != filepath.Join(base, "site", "reference") {
		t.Fatalf("outDir = %q", opts.outDir)
	}
	if opts.manifestDir != filepath.Join(base, "docs-manifests") {
		t.Fatalf("manifestDir = %q", opts.manifestDir)
	}
}

func TestRunDocsGeneratesJekyllMarkdown(t *testing.T) {
	base := t.TempDir()
	highschool := filepath.Join(base, "highschool", "Algebra")
	if err := os.MkdirAll(highschool, 0o755); err != nil {
		t.Fatalf("mkdir highschool: %v", err)
	}
	stdlib := `doc<
Title:
    Absolute value for exact numbers.
About:
    Returns non-negative exact values.
Param:
    Value: Input exact.
Return:
    Non-negative exact.
Example:
    allow exact A 2b=2 Algebra.PosiExact<0 - 3>;;
>
outy seq PosiExact<exact Value> --> <<exact>> {
    maybe (<Value lessly 0>) -> {
        give 0 - Value up;;
    }
    give Value up;;
}`
	if err := os.WriteFile(filepath.Join(highschool, "main.jave"), []byte(stdlib), 0o644); err != nil {
		t.Fatalf("write carryon: %v", err)
	}

	manifestDir := filepath.Join(base, "docs-manifests", "builtins")
	if err := os.MkdirAll(manifestDir, 0o755); err != nil {
		t.Fatalf("mkdir manifests: %v", err)
	}
	manifest := `kind: builtin
name: pront
signature: pront(Value)
about: Builtin print.`
	if err := os.WriteFile(filepath.Join(manifestDir, "pront.yaml"), []byte(manifest), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	outDir := filepath.Join(base, "site", "reference")
	err := runDocs(docsOptions{projectRoot: base, outDir: outDir, manifestDir: filepath.Join(base, "docs-manifests")})
	if err != nil {
		t.Fatalf("runDocs returned error: %v", err)
	}

	carryonPage := filepath.Join(outDir, "carryons", "highschool-algebra.md")
	b, err := os.ReadFile(carryonPage)
	if err != nil {
		t.Fatalf("read carryon page: %v", err)
	}
	text := string(b)
	if !strings.Contains(text, "title: highschool/Algebra Reference") {
		t.Fatalf("expected front matter title in %s", carryonPage)
	}
	if !strings.Contains(text, "`PosiExact`") {
		t.Fatalf("expected PosiExact table-of-contents entry in %s", carryonPage)
	}
	if !strings.Contains(text, "install Algebra from highschool/Algebra;;") {
		t.Fatalf("expected import hint in %s", carryonPage)
	}

	indexPage := filepath.Join(outDir, "index.md")
	indexText, err := os.ReadFile(indexPage)
	if err != nil {
		t.Fatalf("read index page: %v", err)
	}
	if !strings.Contains(string(indexText), "Jave Reference") {
		t.Fatalf("expected index header in %s", indexPage)
	}

	builtinsPage := filepath.Join(outDir, "builtins.md")
	builtinsText, err := os.ReadFile(builtinsPage)
	if err != nil {
		t.Fatalf("read builtins page: %v", err)
	}
	if !strings.Contains(string(builtinsText), "pront") {
		t.Fatalf("expected manifest builtin in %s", builtinsPage)
	}
}
