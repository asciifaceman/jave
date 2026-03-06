//go:build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Default target when running plain `mage`.
var Default = Build

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

// Cmd namespace for running tool stubs during bootstrap.
type Cmd mg.Namespace

// Javec runs the compiler CLI stub.
func (Cmd) Javec() error {
	return sh.RunV("go", "run", "./cmd/javec", "--help")
}

// Baggage runs the package/build manager CLI stub.
func (Cmd) Baggage() error {
	return sh.RunV("go", "run", "./cmd/baggage", "--help")
}

// Javevm runs the VM CLI stub.
func (Cmd) Javevm() error {
	return sh.RunV("go", "run", "./cmd/javevm", "--help")
}

// Version prints a short local workflow version marker.
func Version() {
	fmt.Println("mage workflow: jave v0.1 bootstrap")
}
