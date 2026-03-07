package main

import "testing"

func TestSymbolAtPosition(t *testing.T) {
	src := "prontulate<\"Count=%exact\", 2>;;\n"
	sym := symbolAtPosition(src, position{Line: 0, Character: 2})
	if sym != "prontulate" {
		t.Fatalf("symbol = %q", sym)
	}
}

func TestCallContextAtPosition(t *testing.T) {
	src := "prontulate<\"A=%exact B=%exact\", 1, 2>;;\n"
	name, arg := callContextAtPosition(src, position{Line: 0, Character: 34})
	if name != "prontulate" {
		t.Fatalf("name = %q", name)
	}
	if arg < 1 {
		t.Fatalf("arg index too low: %d", arg)
	}
}
