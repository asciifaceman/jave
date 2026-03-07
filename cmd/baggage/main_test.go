package main

import "testing"

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
