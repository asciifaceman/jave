package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Get workspace root (where this tool is being run from)
	workspaceRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Locate vscode-jave directory
	extensionSource := filepath.Join(workspaceRoot, "vscode-jave")
	if _, err := os.Stat(extensionSource); os.IsNotExist(err) {
		return fmt.Errorf("vscode-jave directory not found at %s\nPlease run this tool from the repository root", extensionSource)
	}

	// Ensure javels binary is present for this host platform when possible.
	if msg, err := ensureJavelsForHost(workspaceRoot, extensionSource); err != nil {
		fmt.Printf("Warning: %v\n", err)
	} else if msg != "" {
		fmt.Println(msg)
	}

	// Get VS Code extensions directory
	extensionsDir, err := getVSCodeExtensionsDir()
	if err != nil {
		return err
	}

	extensionTarget := filepath.Join(extensionsDir, "jave-language-0.1.0")

	fmt.Printf("Installing Jave VS Code extension...\n")
	fmt.Printf("  Source: %s\n", extensionSource)
	fmt.Printf("  Target: %s\n", extensionTarget)
	fmt.Println()

	// Create extensions directory if it doesn't exist
	if err := os.MkdirAll(extensionsDir, 0755); err != nil {
		return fmt.Errorf("failed to create extensions directory: %w", err)
	}

	// Remove existing installation
	if _, err := os.Stat(extensionTarget); err == nil {
		fmt.Println("Removing existing installation...")
		if err := os.RemoveAll(extensionTarget); err != nil {
			return fmt.Errorf("failed to remove existing installation: %w", err)
		}
	}

	// Try to create symlink first (preferred)
	if err := createSymlink(extensionSource, extensionTarget); err != nil {
		fmt.Printf("Symlink creation failed: %v\n", err)
		fmt.Println("Falling back to directory copy...")

		// Fall back to copying
		if err := copyDir(extensionSource, extensionTarget); err != nil {
			return fmt.Errorf("failed to copy extension: %w", err)
		}
		fmt.Println("✓ Extension installed successfully (copied)")
	} else {
		fmt.Println("✓ Extension installed successfully (symlinked)")
	}

	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("1. Reload VS Code: Ctrl+Shift+P → 'Developer: Reload Window'")
	fmt.Println("2. Open a .jave file to verify syntax highlighting")
	fmt.Println("3. Verify hover/signature help (javels) in a .jave file")

	// Check if Jave toolchain is installed
	if !isToolchainInstalled() {
		fmt.Println()
		fmt.Println("Note: To compile and run Jave programs, install the toolchain:")
		fmt.Println("  go install github.com/asciifaceman/jave/cmd/javec@latest")
		fmt.Println("  go install github.com/asciifaceman/jave/cmd/baggage@latest")
		fmt.Println("  go install github.com/asciifaceman/jave/cmd/javevm@latest")
	}

	return nil
}

func ensureJavelsForHost(workspaceRoot, extensionSource string) (string, error) {
	osName, ok := extensionOSName(runtime.GOOS)
	if !ok {
		return "", nil
	}
	archName, ok := extensionArchName(runtime.GOARCH)
	if !ok {
		return "", nil
	}

	binDir := filepath.Join(extensionSource, "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create extension bin dir: %w", err)
	}

	binName := bundledJavelsBinaryName(osName, archName)
	outPath := filepath.Join(binDir, binName)
	if _, err := os.Stat(outPath); err == nil {
		return fmt.Sprintf("✓ javels already bundled for host: %s", outPath), nil
	}

	if !toolExists("go") {
		return "", fmt.Errorf("javels binary missing for host (%s/%s) and Go is not installed", osName, archName)
	}

	cmd := exec.Command("go", "build", "-trimpath", "-ldflags=-s -w", "-o", outPath, "./cmd/javels")
	cmd.Dir = workspaceRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to build javels for host: %w", err)
	}

	if runtime.GOOS != "windows" {
		_ = os.Chmod(outPath, 0755)
	}
	return fmt.Sprintf("✓ built javels for host: %s", outPath), nil
}

func extensionOSName(goos string) (string, bool) {
	switch goos {
	case "windows":
		return "windows", true
	case "darwin":
		return "darwin", true
	case "linux":
		return "linux", true
	default:
		return "", false
	}
}

func extensionArchName(goarch string) (string, bool) {
	switch goarch {
	case "amd64", "arm64":
		return goarch, true
	default:
		return "", false
	}
}

func bundledJavelsBinaryName(osName, archName string) string {
	if osName == "windows" {
		return fmt.Sprintf("javels-%s-%s.exe", osName, archName)
	}
	return fmt.Sprintf("javels-%s-%s", osName, archName)
}

func toolExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func getVSCodeExtensionsDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	switch runtime.GOOS {
	case "windows":
		return filepath.Join(homeDir, ".vscode", "extensions"), nil
	case "darwin", "linux":
		return filepath.Join(homeDir, ".vscode", "extensions"), nil
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func createSymlink(source, target string) error {
	// Convert to absolute paths
	absSource, err := filepath.Abs(source)
	if err != nil {
		return fmt.Errorf("failed to resolve source path: %w", err)
	}

	return os.Symlink(absSource, target)
}

func copyDir(src, dst string) error {
	// Get source directory info
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create destination directory
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	// Read source directory
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// Copy each entry
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectory
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy file
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	// Read source file
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	// Get source file info for permissions
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Write to destination
	return os.WriteFile(dst, data, srcInfo.Mode())
}

func isToolchainInstalled() bool {
	for _, cmd := range []string{"javec", "baggage", "javevm"} {
		if _, err := exec.LookPath(cmd); err != nil {
			return false
		}
	}
	return true
}
