package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		printUsage()
		return
	}

	if args[0] == "--version" || args[0] == "version" {
		fmt.Println("baggage v0.1.0-bootstrap")
		return
	}

	cmd := args[0]
	rest := args[1:]

	switch cmd {
	case "build":
		input, out, err := parseBuildArgs(rest)
		if err != nil {
			fmt.Fprintf(os.Stderr, "baggage: %v\n", err)
			os.Exit(2)
		}
		if err := runBuild(input, out); err != nil {
			fmt.Fprintf(os.Stderr, "baggage build failed: %v\n", err)
			os.Exit(1)
		}
	case "run":
		input, err := parseRunArgs(rest)
		if err != nil {
			fmt.Fprintf(os.Stderr, "baggage: %v\n", err)
			os.Exit(2)
		}
		if err := runProgram(input); err != nil {
			fmt.Fprintf(os.Stderr, "baggage run failed: %v\n", err)
			os.Exit(1)
		}
	case "check":
		if err := runGo("test", "./..."); err != nil {
			fmt.Fprintf(os.Stderr, "baggage check failed: %v\n", err)
			os.Exit(1)
		}
	case "test":
		if err := runGo("test", "./..."); err != nil {
			fmt.Fprintf(os.Stderr, "baggage test failed: %v\n", err)
			os.Exit(1)
		}
	case "new":
		fmt.Println("baggage new: scaffold not implemented yet")
	case "add":
		fmt.Println("baggage add: dependency workflow not implemented yet")
	default:
		fmt.Fprintf(os.Stderr, "baggage: unknown command %q\n", cmd)
		printUsage()
		os.Exit(2)
	}
}

func printUsage() {
	fmt.Println("usage: baggage <build|run|check|test|new|add|version>")
	fmt.Println("  baggage build [input.jave] [-o output.jbin]")
	fmt.Println("  baggage run [input.jave|program.jbin]")
}

func parseBuildArgs(args []string) (input string, out string, err error) {
	input = defaultInput()
	out = ""

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch a {
		case "-o", "--out":
			if i+1 >= len(args) {
				return "", "", fmt.Errorf("missing value for %s", a)
			}
			i++
			out = args[i]
		default:
			if strings.HasPrefix(a, "-") {
				return "", "", fmt.Errorf("unknown flag %s", a)
			}
			input = a
		}
	}

	if out == "" {
		out = artifactPathFor(input)
	}
	return input, out, nil
}

func parseRunArgs(args []string) (string, error) {
	if len(args) > 1 {
		return "", fmt.Errorf("run accepts at most one input argument")
	}
	if len(args) == 0 {
		if env := os.Getenv("JAVE_FILE"); env != "" {
			return env, nil
		}
		return defaultInput(), nil
	}
	return args[0], nil
}

func runBuild(input, out string) error {
	return runGo("run", "./cmd/javec", "--out", out, input)
}

func runProgram(input string) error {
	if strings.HasSuffix(strings.ToLower(input), ".jbin") {
		return runGo("run", "./cmd/javevm", input)
	}
	out := artifactPathFor(input)
	if err := runBuild(input, out); err != nil {
		return err
	}
	return runGo("run", "./cmd/javevm", out)
}

func runGo(args ...string) error {
	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func defaultInput() string {
	if env := os.Getenv("JAVE_FILE"); env != "" {
		return env
	}
	return "examples/hello_world/main.jave"
}

func artifactPathFor(input string) string {
	ext := filepath.Ext(input)
	if ext == "" {
		return input + ".jbin"
	}
	return strings.TrimSuffix(input, ext) + ".jbin"
}
