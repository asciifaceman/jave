package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestArtifactPathFor(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "with jave extension", input: "examples/hello_world/main.jave", want: "examples/hello_world/main.jbin"},
		{name: "already jbin", input: "out/program.jbin", want: "out/program.jbin"},
		{name: "no extension", input: "program", want: "program.jbin"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := artifactPathFor(tt.input)
			if got != tt.want {
				t.Fatalf("artifactPathFor(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseBuildArgs(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		input, out, err := parseBuildArgs(nil)
		if err != nil {
			t.Fatalf("parseBuildArgs returned error: %v", err)
		}
		if input != defaultInput() {
			t.Fatalf("input = %q, want %q", input, defaultInput())
		}
		if out != artifactPathFor(defaultInput()) {
			t.Fatalf("out = %q, want %q", out, artifactPathFor(defaultInput()))
		}
	})

	t.Run("input and output", func(t *testing.T) {
		input, out, err := parseBuildArgs([]string{"examples/conditions/main.jave", "--out", "bin/conditions.jbin"})
		if err != nil {
			t.Fatalf("parseBuildArgs returned error: %v", err)
		}
		if input != "examples/conditions/main.jave" {
			t.Fatalf("input = %q", input)
		}
		if out != "bin/conditions.jbin" {
			t.Fatalf("out = %q", out)
		}
	})

	t.Run("missing out value", func(t *testing.T) {
		_, _, err := parseBuildArgs([]string{"-o"})
		if err == nil {
			t.Fatal("expected error for missing -o value")
		}
	})

	t.Run("unknown flag", func(t *testing.T) {
		_, _, err := parseBuildArgs([]string{"--wat"})
		if err == nil {
			t.Fatal("expected error for unknown flag")
		}
	})
}

func TestParseRunArgs(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		input, err := parseRunArgs(nil)
		if err != nil {
			t.Fatalf("parseRunArgs returned error: %v", err)
		}
		if input != defaultInput() {
			t.Fatalf("input = %q, want %q", input, defaultInput())
		}
	})

	t.Run("explicit input", func(t *testing.T) {
		input, err := parseRunArgs([]string{"program.jbin"})
		if err != nil {
			t.Fatalf("parseRunArgs returned error: %v", err)
		}
		if input != "program.jbin" {
			t.Fatalf("input = %q", input)
		}
	})

	t.Run("too many args", func(t *testing.T) {
		_, err := parseRunArgs([]string{"a", "b"})
		if err == nil {
			t.Fatal("expected error for too many args")
		}
	})
}

func TestParseNewArgs(t *testing.T) {
	t.Run("project only", func(t *testing.T) {
		project, force, err := parseNewArgs([]string{"hello-jave"})
		if err != nil {
			t.Fatalf("parseNewArgs returned error: %v", err)
		}
		if project != "hello-jave" {
			t.Fatalf("project = %q", project)
		}
		if force {
			t.Fatal("force = true, want false")
		}
	})

	t.Run("project with force", func(t *testing.T) {
		project, force, err := parseNewArgs([]string{"hello-jave", "--force"})
		if err != nil {
			t.Fatalf("parseNewArgs returned error: %v", err)
		}
		if project != "hello-jave" {
			t.Fatalf("project = %q", project)
		}
		if !force {
			t.Fatal("force = false, want true")
		}
	})

	t.Run("missing project", func(t *testing.T) {
		_, _, err := parseNewArgs(nil)
		if err == nil {
			t.Fatal("expected error for missing project")
		}
	})

	t.Run("unknown flag", func(t *testing.T) {
		_, _, err := parseNewArgs([]string{"--wat"})
		if err == nil {
			t.Fatal("expected error for unknown flag")
		}
	})

	t.Run("multiple names", func(t *testing.T) {
		_, _, err := parseNewArgs([]string{"one", "two"})
		if err == nil {
			t.Fatal("expected error for multiple names")
		}
	})
}

func TestScaffoldNewProject(t *testing.T) {
	base := t.TempDir()
	target := filepath.Join(base, "hello-jave")

	if err := scaffoldNewProject(target, false); err != nil {
		t.Fatalf("scaffoldNewProject returned error: %v", err)
	}

	mainPath := filepath.Join(target, "main.jave")
	manifestPath := filepath.Join(target, "baggage.jave")
	if _, err := os.Stat(mainPath); err != nil {
		t.Fatalf("expected %s: %v", mainPath, err)
	}
	if _, err := os.Stat(manifestPath); err != nil {
		t.Fatalf("expected %s: %v", manifestPath, err)
	}

	if err := scaffoldNewProject(target, false); err == nil {
		t.Fatal("expected error when target exists without --force")
	}

	if err := scaffoldNewProject(target, true); err != nil {
		t.Fatalf("expected force overwrite to succeed: %v", err)
	}
}

func TestParseAddArgs(t *testing.T) {
	t.Run("default manifest", func(t *testing.T) {
		dep, manifest, err := parseAddArgs([]string{"some/carryon"})
		if err != nil {
			t.Fatalf("parseAddArgs returned error: %v", err)
		}
		if dep != "some/carryon" {
			t.Fatalf("dep = %q", dep)
		}
		if manifest != "baggage.jave" {
			t.Fatalf("manifest = %q", manifest)
		}
	})

	t.Run("custom manifest", func(t *testing.T) {
		dep, manifest, err := parseAddArgs([]string{"--manifest", "project/baggage.jave", "dep/pkg"})
		if err != nil {
			t.Fatalf("parseAddArgs returned error: %v", err)
		}
		if dep != "dep/pkg" {
			t.Fatalf("dep = %q", dep)
		}
		if manifest != "project/baggage.jave" {
			t.Fatalf("manifest = %q", manifest)
		}
	})

	t.Run("missing dep", func(t *testing.T) {
		_, _, err := parseAddArgs(nil)
		if err == nil {
			t.Fatal("expected missing dependency error")
		}
	})

	t.Run("unknown flag", func(t *testing.T) {
		_, _, err := parseAddArgs([]string{"--wat"})
		if err == nil {
			t.Fatal("expected unknown flag error")
		}
	})

	t.Run("multiple deps", func(t *testing.T) {
		_, _, err := parseAddArgs([]string{"a", "b"})
		if err == nil {
			t.Fatal("expected multiple dependency error")
		}
	})
}

func TestAddDependencyToManifest(t *testing.T) {
	manifest := filepath.Join(t.TempDir(), "baggage.jave")
	initial := "carryon \"hello-jave\"\nlang \"v0.1\"\nentry \"main.jave\"\n"
	if err := os.WriteFile(manifest, []byte(initial), 0o644); err != nil {
		t.Fatalf("write initial manifest: %v", err)
	}

	added, err := addDependencyToManifest(manifest, "some/carryon")
	if err != nil {
		t.Fatalf("addDependencyToManifest returned error: %v", err)
	}
	if !added {
		t.Fatal("expected dependency to be added")
	}

	b, err := os.ReadFile(manifest)
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	content := string(b)
	if !strings.Contains(content, "dep \"some/carryon\"\n") {
		t.Fatalf("expected dep line in manifest, got: %q", content)
	}

	added, err = addDependencyToManifest(manifest, "some/carryon")
	if err != nil {
		t.Fatalf("second addDependencyToManifest returned error: %v", err)
	}
	if added {
		t.Fatal("expected duplicate dependency to be ignored")
	}
}
