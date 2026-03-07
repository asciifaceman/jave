package main

import "testing"

func TestSymbolAtPosition(t *testing.T) {
	src := "Prontulate<\"Count=%exact\", 2>;;\n"
	sym := symbolAtPosition(src, position{Line: 0, Character: 2})
	if sym != "Prontulate" {
		t.Fatalf("symbol = %q", sym)
	}
}

func TestCallContextAtPosition(t *testing.T) {
	src := "Prontulate<\"A=%exact B=%exact\", 1, 2>;;\n"
	name, arg := callContextAtPosition(src, position{Line: 0, Character: 34})
	if name != "Prontulate" {
		t.Fatalf("name = %q", name)
	}
	if arg < 1 {
		t.Fatalf("arg index too low: %d", arg)
	}
}
