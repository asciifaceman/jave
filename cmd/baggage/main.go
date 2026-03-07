package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		printUsage()
		return
	}

	if args[0] == "--version" || args[0] == "version" {
		fmt.Println("baggage v0.1.0")
		return
	}

	cmd := args[0]
	rest := args[1:]

	switch cmd {
	case "build":
		opts, err := parseBuildArgs(rest)
		if err != nil {
			fmt.Fprintf(os.Stderr, "baggage: %v\n", err)
			os.Exit(2)
		}
		if err := runBuild(opts); err != nil {
			fmt.Fprintf(os.Stderr, "baggage build failed: %v\n", err)
			os.Exit(1)
		}
	case "run":
		opts, err := parseRunArgs(rest)
		if err != nil {
			fmt.Fprintf(os.Stderr, "baggage: %v\n", err)
			os.Exit(2)
		}
		if err := runProgram(opts); err != nil {
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
	case "docs":
		opts, err := parseDocsArgs(rest)
		if err != nil {
			fmt.Fprintf(os.Stderr, "baggage: %v\n", err)
			os.Exit(2)
		}
		if err := runDocs(opts); err != nil {
			fmt.Fprintf(os.Stderr, "baggage docs failed: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "baggage: unknown command %q\n", cmd)
		printUsage()
		os.Exit(2)
	}
}

func printUsage() {
	fmt.Println("usage: baggage <build|run|check|test|new|add|docs|version>")
	fmt.Println("  baggage new <project-name> [--force]")
	fmt.Println("  baggage add <carryon-path> [--manifest baggage.jave]")
	fmt.Println("  baggage docs [--project-root dir] [--out-dir site/reference] [--manifest-dir docs-manifests]")
	fmt.Println("  baggage build [input.jave] [-o output.jbin] [--trace-imports] [--project-root dir] [--sponsor-notice mode] [--sponsor-redacted] [--sponsor-quiet]")
	fmt.Println("  baggage run [input.jave|program.jbin] [--trace-imports] [--project-root dir] [--sponsor-notice mode] [--sponsor-redacted] [--sponsor-quiet]")
}

type buildOptions struct {
	input           string
	out             string
	traceImports    bool
	projectRoot     string
	sponsorMode     string
	sponsorQuiet    bool
	sponsorRedacted bool
}

type runOptions struct {
	input           string
	traceImports    bool
	projectRoot     string
	sponsorMode     string
	sponsorQuiet    bool
	sponsorRedacted bool
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

func parseBuildArgs(args []string) (buildOptions, error) {
	opts := buildOptions{}

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch a {
		case "-o", "--out":
			if i+1 >= len(args) {
				return buildOptions{}, fmt.Errorf("missing value for %s", a)
			}
			i++
			opts.out = args[i]
		case "--trace-imports":
			opts.traceImports = true
		case "--sponsor-notice":
			if i+1 >= len(args) {
				return buildOptions{}, fmt.Errorf("missing value for --sponsor-notice")
			}
			i++
			opts.sponsorMode = args[i]
		case "--sponsor-redacted":
			opts.sponsorRedacted = true
		case "--sponsor-quiet":
			opts.sponsorQuiet = true
		case "--project-root":
			if i+1 >= len(args) {
				return buildOptions{}, fmt.Errorf("missing value for --project-root")
			}
			i++
			opts.projectRoot = args[i]
		default:
			if strings.HasPrefix(a, "-") {
				return buildOptions{}, fmt.Errorf("unknown flag %s", a)
			}
			opts.input = a
		}
	}

	if opts.input == "" {
		input, err := discoverDefaultInput()
		if err != nil {
			return buildOptions{}, err
		}
		opts.input = input
	}

	if opts.out == "" {
		opts.out = artifactPathFor(opts.input)
	}
	if opts.projectRoot == "" {
		opts.projectRoot = defaultProjectRoot()
	}
	if opts.sponsorMode == "" {
		opts.sponsorMode = "full"
	}
	return opts, nil
}

func parseRunArgs(args []string) (runOptions, error) {
	opts := runOptions{}

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch a {
		case "--trace-imports":
			opts.traceImports = true
		case "--sponsor-notice":
			if i+1 >= len(args) {
				return runOptions{}, fmt.Errorf("missing value for --sponsor-notice")
			}
			i++
			opts.sponsorMode = args[i]
		case "--sponsor-redacted":
			opts.sponsorRedacted = true
		case "--sponsor-quiet":
			opts.sponsorQuiet = true
		case "--project-root":
			if i+1 >= len(args) {
				return runOptions{}, fmt.Errorf("missing value for --project-root")
			}
			i++
			opts.projectRoot = args[i]
		default:
			if strings.HasPrefix(a, "-") {
				return runOptions{}, fmt.Errorf("unknown flag %s", a)
			}
			if opts.input != "" {
				return runOptions{}, fmt.Errorf("run accepts at most one input argument")
			}
			opts.input = a
		}
	}

	if opts.input == "" {
		input, err := discoverDefaultInput()
		if err != nil {
			return runOptions{}, err
		}
		opts.input = input
	}
	if opts.projectRoot == "" {
		opts.projectRoot = defaultProjectRoot()
	}
	if opts.sponsorMode == "" {
		opts.sponsorMode = "full"
	}
	return opts, nil
}

func runBuild(opts buildOptions) error {
	args := []string{"run", "./cmd/javec"}
	if opts.traceImports {
		args = append(args, "--trace-imports")
	}
	if opts.sponsorQuiet {
		args = append(args, "--sponsor-quiet")
	}
	if opts.sponsorRedacted {
		args = append(args, "--sponsor-redacted")
	}
	if opts.sponsorMode != "" {
		args = append(args, "--sponsor-notice", opts.sponsorMode)
	}
	if opts.projectRoot != "" {
		args = append(args, "--project-root", opts.projectRoot)
	}
	args = append(args, "--out", opts.out, opts.input)
	return runGo(args...)
}

func runProgram(opts runOptions) error {
	if strings.HasSuffix(strings.ToLower(opts.input), ".jbin") {
		return runGo("run", "./cmd/javevm", opts.input)
	}
	out := artifactPathFor(opts.input)
	if err := runBuild(buildOptions{input: opts.input, out: out, traceImports: opts.traceImports, projectRoot: opts.projectRoot, sponsorMode: opts.sponsorMode, sponsorQuiet: opts.sponsorQuiet, sponsorRedacted: opts.sponsorRedacted}); err != nil {
		return err
	}
	return runGo("run", "./cmd/javevm", out)
}

func defaultProjectRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return wd
}

func runGo(args ...string) error {
	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func defaultInput() string {
	input, err := discoverDefaultInput()
	if err != nil {
		return "examples/hello_world/main.jave"
	}
	return input
}

func discoverDefaultInput() (string, error) {
	if env := os.Getenv("JAVE_FILE"); env != "" {
		return env, nil
	}

	manifestPath := "baggage.jave"
	if _, err := os.Stat(manifestPath); err == nil {
		return inputFromManifest(manifestPath)
	} else if !os.IsNotExist(err) {
		return "", err
	}

	return "examples/hello_world/main.jave", nil
}

func inputFromManifest(path string) (string, error) {
	manifest, err := readManifest(path)
	if err != nil {
		return "", err
	}
	if manifest.Entry == "" {
		return "", fmt.Errorf("manifest %s is missing entry", path)
	}

	entry := manifest.Entry
	if filepath.IsAbs(entry) {
		return filepath.Clean(entry), nil
	}

	base := filepath.Dir(path)
	return filepath.Clean(filepath.Join(base, entry)), nil
}

type baggageManifest struct {
	Carryon string
	Lang    string
	Entry   string
	Deps    []string
}

func readManifest(path string) (baggageManifest, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return baggageManifest{}, err
	}

	manifest := baggageManifest{}
	lines := strings.Split(string(b), "\n")
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			return baggageManifest{}, fmt.Errorf("invalid manifest line: %q", line)
		}

		key := strings.TrimSpace(parts[0])
		valuePart := strings.TrimSpace(parts[1])
		value, err := strconv.Unquote(valuePart)
		if err != nil {
			return baggageManifest{}, fmt.Errorf("invalid quoted value in line %q: %w", line, err)
		}

		switch key {
		case "carryon":
			manifest.Carryon = value
		case "lang":
			manifest.Lang = value
		case "entry":
			manifest.Entry = value
		case "dep":
			manifest.Deps = append(manifest.Deps, value)
		default:
			return baggageManifest{}, fmt.Errorf("unknown manifest key %q", key)
		}
	}

	return manifest, nil
}

func artifactPathFor(input string) string {
	ext := filepath.Ext(input)
	if ext == "" {
		return input + ".jbin"
	}
	return strings.TrimSuffix(input, ext) + ".jbin"
}
