package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/asciifaceman/jave/internal/jbin"
	"github.com/asciifaceman/jave/internal/lexer"
	"github.com/asciifaceman/jave/internal/lowering"
	"github.com/asciifaceman/jave/internal/parser"
	"github.com/asciifaceman/jave/internal/runtime"
	"github.com/asciifaceman/jave/internal/sema"
)

func main() {
	showVersion := flag.Bool("version", false, "print javec version")
	showTokens := flag.Bool("tokens", false, "print lexer tokens")
	runProgram := flag.Bool("run", false, "run Foremost after successful analysis")
	outPath := flag.String("out", "", "output .jbin file path (defaults to <input>.jbin)")
	flag.Parse()

	if *showVersion {
		fmt.Println("javec v0.1.0-bootstrap")
		return
	}

	if flag.NArg() == 0 {
		fmt.Println("usage: javec [--version] [--tokens] [--run] [--out file.jbin] <input.jave>")
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

	program, parseDiags := parser.Parse(tokens)
	if len(parseDiags) > 0 {
		for _, d := range parseDiags {
			fmt.Fprintf(os.Stderr, "%s:%d:%d: %s: %s\n", path, d.Pos.Line, d.Pos.Column, d.Severity, d.Message)
		}
		os.Exit(1)
	}

	semaDiags := sema.Analyze(program)
	if len(semaDiags) > 0 {
		for _, d := range semaDiags {
			fmt.Fprintf(os.Stderr, "%s:%d:%d: %s: %s\n", path, d.Pos.Line, d.Pos.Column, d.Severity, d.Message)
		}
		os.Exit(1)
	}

	irProgram, lowerDiags := lowering.Lower(program)
	if len(lowerDiags) > 0 {
		for _, d := range lowerDiags {
			fmt.Fprintf(os.Stderr, "%s:%d:%d: %s: %s\n", path, d.Pos.Line, d.Pos.Column, d.Severity, d.Message)
		}
		os.Exit(1)
	}

	if *showTokens {
		for _, tok := range tokens {
			fmt.Printf("%d:%d %-14s %q\n", tok.Pos.Line, tok.Pos.Column, tok.Kind, tok.Lexeme)
		}
	}

	fmt.Printf("javec: lowering succeeded (%d tokens, %d sequences)\n", len(tokens), len(program.Sequences))

	emitPath := *outPath
	if emitPath == "" {
		ext := filepath.Ext(path)
		emitPath = strings.TrimSuffix(path, ext) + ".jbin"
	}
	if err := jbin.WriteFile(emitPath, irProgram); err != nil {
		fmt.Fprintf(os.Stderr, "javec: unable to write %q: %v\n", emitPath, err)
		os.Exit(1)
	}
	fmt.Printf("javec: emitted %s\n", emitPath)

	if *runProgram {
		if err := runtime.Execute(irProgram, os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "javec: runtime error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	fmt.Println("next: run with javevm <file.jbin>")
}
