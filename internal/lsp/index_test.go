package lsp

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuildIndexFromSource_DocstringFields(t *testing.T) {
	src := `doc<
Title:
    Positive exact value
About:
    Returns a non-negative exact.
Param:
    Value: Exact input.
Return:
    Non-negative exact.
>
outy seq PosiExact<exact Value> --> <<exact>> {
    give Value up;;
}`
	idx, err := BuildIndexFromSource(src, "unit-test")
	if err != nil {
		t.Fatalf("BuildIndexFromSource error: %v", err)
	}
	doc, ok := idx.Symbols["PosiExact"]
	if !ok {
		t.Fatal("expected PosiExact in index")
	}
	if doc.Title != "Positive exact value" {
		t.Fatalf("title = %q", doc.Title)
	}
	if len(doc.Params) != 1 || doc.Params[0].Description != "Exact input." {
		t.Fatalf("unexpected params: %+v", doc.Params)
	}
	if doc.ReturnDesc != "Non-negative exact." {
		t.Fatalf("return desc = %q", doc.ReturnDesc)
	}
}

func TestBuildIndexFromManifests_Builtin(t *testing.T) {
	root := t.TempDir()
	manifestDir := filepath.Join(root, "docs-manifests", "builtins")
	mustMkdirAll(t, manifestDir)
	mustWriteFile(t, filepath.Join(manifestDir, "pront.yaml"), `kind: builtin
name: pront
signature: pront(Value)
title: Print value
about: Writes one value.`)

	idx, err := BuildIndexFromManifests(filepath.Join(root, "docs-manifests"))
	if err != nil {
		t.Fatalf("BuildIndexFromManifests error: %v", err)
	}
	doc, ok := idx.Symbols["pront"]
	if !ok {
		t.Fatal("expected pront in index")
	}
	if doc.Kind != "builtin" {
		t.Fatalf("kind = %q", doc.Kind)
	}
}

func mustMkdirAll(t *testing.T, p string) {
	t.Helper()
	if err := os.MkdirAll(p, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", p, err)
	}
}

func mustWriteFile(t *testing.T, p, c string) {
	t.Helper()
	if err := os.WriteFile(p, []byte(c), 0o644); err != nil {
		t.Fatalf("write %s: %v", p, err)
	}
}
