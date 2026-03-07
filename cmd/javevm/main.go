package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/asciifaceman/jave/internal/jbin"
	"github.com/asciifaceman/jave/internal/runtime"
)

func main() {
	showVersion := flag.Bool("version", false, "print javevm version")
	flag.Parse()

	if *showVersion {
		fmt.Println("javevm v0.1.0")
		return
	}

	if flag.NArg() == 0 {
		fmt.Println("usage: javevm <program.jbin>")
		return
	}

	path := flag.Arg(0)
	program, err := jbin.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "javevm: unable to read %q: %v\n", path, err)
		os.Exit(1)
	}

	if err := runtime.Execute(program, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "javevm: runtime error: %v\n", err)
		os.Exit(1)
	}
}
