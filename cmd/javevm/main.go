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
		fmt.Println("usage: javevm <program.jbin> [program args...]")
		return
	}

	path := flag.Arg(0)
	runtimeArgs := []string{}
	if flag.NArg() > 1 {
		runtimeArgs = append(runtimeArgs, flag.Args()[1:]...)
	}
	program, err := jbin.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "javevm: unable to read %q: %v\n", path, err)
		os.Exit(1)
	}

	runOpts := runtime.ExecuteOptions{Stdout: os.Stdout, Stderr: os.Stderr, Args: runtimeArgs}
	if err := runtime.ExecuteWithOptions(program, runOpts); err != nil {
		fmt.Fprintf(os.Stderr, "javevm: runtime error: %v\n", err)
		os.Exit(runtime.ExitCodeForError(err))
	}
}
