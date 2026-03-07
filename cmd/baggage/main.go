package main

import (
	"fmt"
	"io/fs"
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
		if err := runNew(rest); err != nil {
			fmt.Fprintf(os.Stderr, "baggage new failed: %v\n", err)
			os.Exit(1)
		}
	case "add":
		if err := runAdd(rest); err != nil {
			fmt.Fprintf(os.Stderr, "baggage add failed: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "baggage: unknown command %q\n", cmd)
		printUsage()
		os.Exit(2)
	}
}

func printUsage() {
	fmt.Println("usage: baggage <build|run|check|test|new|add|version>")
	fmt.Println("  baggage new <project-name> [--force]")
	fmt.Println("  baggage add <carryon-path> [--manifest baggage.jave]")
	fmt.Println("  baggage build [input.jave] [-o output.jbin]")
	fmt.Println("  baggage run [input.jave|program.jbin]")
}

func runAdd(args []string) error {
	dep, manifest, err := parseAddArgs(args)
	if err != nil {
		return err
	}
	added, err := addDependencyToManifest(manifest, dep)
	if err != nil {
		return err
	}
	if added {
		fmt.Printf("baggage: added dependency %q to %s\n", dep, manifest)
		return nil
	}
	fmt.Printf("baggage: dependency %q already present in %s\n", dep, manifest)
	return nil
}

func parseAddArgs(args []string) (dep string, manifest string, err error) {
	manifest = "baggage.jave"

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch a {
		case "--manifest":
			if i+1 >= len(args) {
				return "", "", fmt.Errorf("missing value for --manifest")
			}
			i++
			manifest = args[i]
		default:
			if strings.HasPrefix(a, "-") {
				return "", "", fmt.Errorf("unknown flag %s", a)
			}
			if dep != "" {
				return "", "", fmt.Errorf("add accepts exactly one dependency path")
			}
			dep = a
		}
	}

	if dep == "" {
		return "", "", fmt.Errorf("missing dependency path")
	}

	if strings.ContainsAny(dep, "\r\n") {
		return "", "", fmt.Errorf("dependency path may not contain newlines")
	}

	return dep, manifest, nil
}

func addDependencyToManifest(path, dep string) (bool, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}

	content := string(b)
	line := fmt.Sprintf("dep %q", dep)
	if strings.Contains(content, line) {
		return false, nil
	}

	if content != "" && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	content += line + "\n"

	if err := os.WriteFile(path, []byte(content), fs.FileMode(0o644)); err != nil {
		return false, err
	}
	return true, nil
}

func runNew(args []string) error {
	project, force, err := parseNewArgs(args)
	if err != nil {
		return err
	}
	if err := scaffoldNewProject(project, force); err != nil {
		return err
	}
	fmt.Printf("baggage: created project %q\n", project)
	fmt.Printf("next: go run ./cmd/javec %s\n", filepath.ToSlash(filepath.Join(project, "main.jave")))
	return nil
}

func parseNewArgs(args []string) (project string, force bool, err error) {
	for _, a := range args {
		switch a {
		case "--force":
			force = true
		default:
			if strings.HasPrefix(a, "-") {
				return "", false, fmt.Errorf("unknown flag %s", a)
			}
			if project != "" {
				return "", false, fmt.Errorf("new accepts exactly one project name")
			}
			project = a
		}
	}

	if project == "" {
		return "", false, fmt.Errorf("missing project name")
	}

	if project == "." || project == ".." {
		return "", false, fmt.Errorf("invalid project name %q", project)
	}

	clean := filepath.Clean(project)
	if clean != project && clean != filepath.FromSlash(project) {
		return "", false, fmt.Errorf("project path must be a clean relative path")
	}

	return project, force, nil
}

func scaffoldNewProject(project string, force bool) error {
	if st, err := os.Stat(project); err == nil {
		if !st.IsDir() {
			return fmt.Errorf("%q exists and is not a directory", project)
		}
		if !force {
			return fmt.Errorf("%q already exists (use --force to overwrite scaffold files)", project)
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	if err := os.MkdirAll(project, fs.FileMode(0o755)); err != nil {
		return err
	}

	mainPath := filepath.Join(project, "main.jave")
	manifestPath := filepath.Join(project, "baggage.jave")

	manifestName := filepath.Base(project)
	manifest := fmt.Sprintf("carryon %q\nlang %q\nentry %q\n", manifestName, "v0.1", "main.jave")
	mainProgram := "outy seq Foremost<> --> <<nada>> {\n    pront(\"hello, jave\");;\n    give up;;\n}\n"

	if err := os.WriteFile(mainPath, []byte(mainProgram), fs.FileMode(0o644)); err != nil {
		return err
	}
	if err := os.WriteFile(manifestPath, []byte(manifest), fs.FileMode(0o644)); err != nil {
		return err
	}
	return nil
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
