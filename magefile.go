//go:build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/sh"
)

// Default target when running plain `mage`.
var Default = Help

// Help prints the most useful project workflows.
func Help() {
	fmt.Println("Jave Mage Commands")
	fmt.Println("")
	fmt.Println("Core:")
	fmt.Println("  mage build        # go build ./...")
	fmt.Println("  mage test         # go test ./...")
	fmt.Println("  mage check        # fmt + vet + test")
	fmt.Println("  mage bootstrap    # go mod tidy + build")
	fmt.Println("")
	fmt.Println("Tooling:")
	fmt.Println("  mage installExtension  # install VS Code extension for Jave syntax")
	fmt.Println("")
	fmt.Println("Compiler/Runtime:")
	fmt.Println("  mage runJavec                 # analyze JAVE_FILE or examples/hello_world/main.jave")
	fmt.Println("                                # PowerShell: $env:JAVE_FILE='examples/conditions/main.jave' ; mage runJavec")
	fmt.Println("  mage runJavecTokens           # token dump + analysis for hello world")
	fmt.Println("  mage runExampleConditions     # javec --run examples/conditions/main.jave")
	fmt.Println("  mage runExampleMultiTable     # javec --run examples/multi_dimensional_tables/main.jave")
	fmt.Println("")
	fmt.Println("Tool Stubs:")
	fmt.Println("  mage runBaggage    # run baggage CLI")
	fmt.Println("  mage runJavevm     # run javevm for JBIN_FILE or examples/hello_world/main.jbin")
	fmt.Println("")
	fmt.Println("Tip: use `mage -l` to see all generated targets.")
}

// Build compiles all Go packages in the repository.
func Build() error {
	return sh.RunV("go", "build", "./...")
}

// Test runs all Go tests in the repository.
func Test() error {
	return sh.RunV("go", "test", "./...")
}

// Check runs formatting, vet, and tests for a basic quality gate.
func Check() error {
	if err := sh.RunV("go", "fmt", "./..."); err != nil {
		return err
	}
	if err := sh.RunV("go", "vet", "./..."); err != nil {
		return err
	}
	return Test()
}

// Bootstrap prepares local dependencies and validates the Go workspace.
func Bootstrap() error {
	if err := sh.RunV("go", "mod", "tidy"); err != nil {
		return err
	}
	return Build()
}

// Run lists runnable workflow commands.
func Run() { Help() }

// RunJavec runs javec for JAVE_FILE or a default example file.
func RunJavec() error {
	input := "examples/hello_world/main.jave"
	if env := os.Getenv("JAVE_FILE"); env != "" {
		input = env
	}
	return sh.RunV("go", "run", "./cmd/javec", input)
}

// RunJavecTokens runs javec with token output against a default example file.
func RunJavecTokens() error {
	return sh.RunV("go", "run", "./cmd/javec", "--tokens", "examples/hello_world/main.jave")
}

// RunBaggage runs the baggage CLI.
func RunBaggage() error {
	return sh.RunV("go", "run", "./cmd/baggage")
}

// RunJavevm runs the javevm CLI.
func RunJavevm() error {
	input := "examples/hello_world/main.jbin"
	if env := os.Getenv("JBIN_FILE"); env != "" {
		input = env
	}
	if _, err := os.Stat(input); err != nil {
		if err := sh.RunV("go", "run", "./cmd/javec", "examples/hello_world/main.jave"); err != nil {
			return err
		}
	}
	return sh.RunV("go", "run", "./cmd/javevm", input)
}

// RunExampleConditions executes the conditions example through javec runtime mode.
func RunExampleConditions() error {
	return sh.RunV("go", "run", "./cmd/javec", "--run", "examples/conditions/main.jave")
}

// RunExampleMultiTable executes the nested table example through javec runtime mode.
func RunExampleMultiTable() error {
	return sh.RunV("go", "run", "./cmd/javec", "--run", "examples/multi_dimensional_tables/main.jave")
}

// Version prints a short local workflow version marker.
func Version() {
	fmt.Println("mage workflow: jave v0.1 bootstrap")
}

// InstallExtension installs the VS Code extension for Jave syntax highlighting.
func InstallExtension() error {
	return sh.RunV("go", "run", "./tools/install-extension")
}
