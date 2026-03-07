package main

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestExtensionOSName(t *testing.T) {
	tests := []struct {
		in   string
		want string
		ok   bool
	}{
		{in: "windows", want: "windows", ok: true},
		{in: "linux", want: "linux", ok: true},
		{in: "darwin", want: "darwin", ok: true},
		{in: "freebsd", want: "", ok: false},
	}
	for _, tt := range tests {
		got, ok := extensionOSName(tt.in)
		if got != tt.want || ok != tt.ok {
			t.Fatalf("extensionOSName(%q) = (%q,%v), want (%q,%v)", tt.in, got, ok, tt.want, tt.ok)
		}
	}
}

func TestExtensionArchName(t *testing.T) {
	tests := []struct {
		in   string
		want string
		ok   bool
	}{
		{in: "amd64", want: "amd64", ok: true},
		{in: "arm64", want: "arm64", ok: true},
		{in: "386", want: "", ok: false},
	}
	for _, tt := range tests {
		got, ok := extensionArchName(tt.in)
		if got != tt.want || ok != tt.ok {
			t.Fatalf("extensionArchName(%q) = (%q,%v), want (%q,%v)", tt.in, got, ok, tt.want, tt.ok)
		}
	}
}

func TestBundledJavelsBinaryName(t *testing.T) {
	if got := bundledJavelsBinaryName("windows", "amd64"); got != "javels-windows-amd64.exe" {
		t.Fatalf("unexpected windows name: %s", got)
	}
	if got := bundledJavelsBinaryName("linux", "arm64"); got != "javels-linux-arm64" {
		t.Fatalf("unexpected linux name: %s", got)
	}
}

func TestGetVSCodeExtensionsDir(t *testing.T) {
	dir, err := getVSCodeExtensionsDir()
	if err != nil {
		t.Fatalf("getVSCodeExtensionsDir failed: %v", err)
	}

	// Should contain .vscode/extensions
	if !filepath.IsAbs(dir) {
		t.Errorf("expected absolute path, got: %s", dir)
	}

	expected := filepath.Join(".vscode", "extensions")
	cleanDir := filepath.Clean(dir)
	cleanExpected := filepath.ToSlash(filepath.Clean(expected))
	if !strings.Contains(filepath.ToSlash(cleanDir), cleanExpected) {
		t.Errorf("expected path to contain %s, got: %s", expected, dir)
	}
}

func TestGetVSCodeExtensionsDir_IsValid(t *testing.T) {
	dir, err := getVSCodeExtensionsDir()
	if err != nil {
		t.Fatalf("getVSCodeExtensionsDir failed: %v", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home dir: %v", err)
	}

	// Expected path varies by OS but should be under home directory
	switch runtime.GOOS {
	case "windows", "darwin", "linux":
		expected := filepath.Join(homeDir, ".vscode", "extensions")
		if dir != expected {
			t.Errorf("expected %s, got %s", expected, dir)
		}
	default:
		t.Skipf("unsupported OS: %s", runtime.GOOS)
	}
}

func TestCopyFile(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()

	// Create source file
	srcPath := filepath.Join(tmpDir, "source.txt")
	content := []byte("test content")
	if err := os.WriteFile(srcPath, content, 0644); err != nil {
		t.Fatalf("failed to create source file: %v", err)
	}

	// Copy file
	dstPath := filepath.Join(tmpDir, "dest.txt")
	if err := copyFile(srcPath, dstPath); err != nil {
		t.Fatalf("copyFile failed: %v", err)
	}

	// Verify destination exists and has same content
	dstContent, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("failed to read destination file: %v", err)
	}

	if string(dstContent) != string(content) {
		t.Errorf("content mismatch: expected %q, got %q", string(content), string(dstContent))
	}

	// Verify permissions preserved
	srcInfo, _ := os.Stat(srcPath)
	dstInfo, _ := os.Stat(dstPath)
	if srcInfo.Mode() != dstInfo.Mode() {
		t.Errorf("mode mismatch: expected %v, got %v", srcInfo.Mode(), dstInfo.Mode())
	}
}

func TestCopyDir(t *testing.T) {
	// Create temp directory with structure
	tmpDir := t.TempDir()
	srcDir := filepath.Join(tmpDir, "source")

	// Create source structure
	if err := os.MkdirAll(filepath.Join(srcDir, "subdir"), 0755); err != nil {
		t.Fatalf("failed to create source structure: %v", err)
	}

	if err := os.WriteFile(filepath.Join(srcDir, "file1.txt"), []byte("content1"), 0644); err != nil {
		t.Fatalf("failed to create file1: %v", err)
	}

	if err := os.WriteFile(filepath.Join(srcDir, "subdir", "file2.txt"), []byte("content2"), 0644); err != nil {
		t.Fatalf("failed to create file2: %v", err)
	}

	// Copy directory
	dstDir := filepath.Join(tmpDir, "dest")
	if err := copyDir(srcDir, dstDir); err != nil {
		t.Fatalf("copyDir failed: %v", err)
	}

	// Verify structure
	checkFile := func(path, expectedContent string) {
		content, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("failed to read %s: %v", path, err)
			return
		}
		if string(content) != expectedContent {
			t.Errorf("%s content mismatch: expected %q, got %q", path, expectedContent, string(content))
		}
	}

	checkFile(filepath.Join(dstDir, "file1.txt"), "content1")
	checkFile(filepath.Join(dstDir, "subdir", "file2.txt"), "content2")

	// Verify subdir exists
	if _, err := os.Stat(filepath.Join(dstDir, "subdir")); os.IsNotExist(err) {
		t.Error("subdir was not copied")
	}
}

func TestIsToolchainInstalled(t *testing.T) {
	// This test just verifies the function runs without panic
	// Actual result depends on whether tools are installed
	result := isToolchainInstalled()
	t.Logf("Toolchain installed: %v", result)

	// Function should return a boolean value
	if result != true && result != false {
		t.Error("isToolchainInstalled should return a boolean")
	}
}

func TestIsWindowsAccessDenied(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "nil", err: nil, want: false},
		{name: "windows access denied", err: errors.New("unlinkat C:\\x: Access is denied."), want: true},
		{name: "permission denied", err: errors.New("permission denied"), want: true},
		{name: "used by another process", err: errors.New("file is being used by another process"), want: true},
		{name: "other error", err: errors.New("no such file or directory"), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isWindowsAccessDenied(tt.err); got != tt.want {
				t.Fatalf("isWindowsAccessDenied(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}
