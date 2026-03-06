package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/asciifaceman/jave/internal/lexer"
)

func main() {
	showVersion := flag.Bool("version", false, "print javec version")
	showTokens := flag.Bool("tokens", false, "print lexer tokens")
	flag.Parse()

	if *showVersion {
		fmt.Println("javec v0.1.0-bootstrap")
		return
	}

	if flag.NArg() == 0 {
		fmt.Println("javec: compiler bootstrap stub")
		fmt.Println("usage: javec [--version] [--tokens] <input.jave>")
		return
	}

	path := flag.Arg(0)
	src, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "javec: unable to read %q: %v\n", path, err)
		os.Exit(1)
	}

	tokens, diags := lexer.Lex(string(src))
	if len(diags) > 0 {
		for _, d := range diags {
			fmt.Fprintf(os.Stderr, "%s:%d:%d: %s: %s\n", path, d.Pos.Line, d.Pos.Column, d.Severity, d.Message)
		}
		os.Exit(1)
	}

	if *showTokens {
		for _, tok := range tokens {
			fmt.Printf("%d:%d %-14s %q\n", tok.Pos.Line, tok.Pos.Column, tok.Kind, tok.Lexeme)
		}
	}

	fmt.Printf("javec: lexical analysis succeeded (%d tokens)\n", len(tokens))
	fmt.Println("next: parser + ast + diagnostics pipeline")
}
