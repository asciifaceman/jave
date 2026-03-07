package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/asciifaceman/jave/internal/ast"
	"github.com/asciifaceman/jave/internal/diagnostics"
	"github.com/asciifaceman/jave/internal/jbin"
	"github.com/asciifaceman/jave/internal/lexer"
	"github.com/asciifaceman/jave/internal/lowering"
	"github.com/asciifaceman/jave/internal/parser"
	"github.com/asciifaceman/jave/internal/runtime"
	"github.com/asciifaceman/jave/internal/sema"
	"github.com/asciifaceman/jave/internal/sponsor"
	"github.com/asciifaceman/jave/internal/token"
)

func main() {
	showVersion := flag.Bool("version", false, "print javec version")
	showTokens := flag.Bool("tokens", false, "print lexer tokens")
	traceImports := flag.Bool("trace-imports", false, "print import resolution trace")
	runProgram := flag.Bool("run", false, "run Foremost after successful analysis")
	projectRoot := flag.String("project-root", "", "project root for highschool import resolution")
	sponsorNoticeMode := flag.String("sponsor-notice", "full", "sponsor notice mode: full|redacted|off")
	sponsorRedacted := flag.Bool("sponsor-redacted", false, "alias for --sponsor-notice redacted")
	sponsorQuiet := flag.Bool("sponsor-quiet", false, "alias for --sponsor-notice off")
	outPath := flag.String("out", "", "output .jbin file path (defaults to <input>.jbin)")
	flag.Parse()

	if *showVersion {
		fmt.Println("javec v0.1.0")
		return
	}

	if flag.NArg() == 0 {
		fmt.Println("usage: javec [--version] [--tokens] [--trace-imports] [--project-root dir] [--sponsor-notice mode] [--sponsor-redacted] [--sponsor-quiet] [--run] [--out file.jbin] <input.jave>")
		return
	}

	mode, err := resolveSponsorMode(*sponsorNoticeMode, *sponsorRedacted, *sponsorQuiet)
	if err != nil {
		fmt.Fprintf(os.Stderr, "javec: %v\n", err)
		os.Exit(2)
	}

	path := flag.Arg(0)
	for _, line := range sponsor.RenderLines(mode, path) {
		fmt.Fprintln(os.Stderr, line)
	}

	loadOpts := loadOptions{ProjectRoot: strings.TrimSpace(*projectRoot)}
	program, tokensByFile, importTrace, loadDiags, sourceMap, err := loadProgramWithImports(path, loadOpts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "javec: import/load error: %v\n", err)
		os.Exit(1)
	}
	if len(loadDiags) > 0 {
		hasErrors := false
		for _, fd := range loadDiags {
			fmt.Fprintf(os.Stderr, "%s:%d:%d: %s: %s\n", fd.path, fd.diag.Pos.Line, fd.diag.Pos.Column, fd.diag.Severity, fd.diag.Message)
			if fd.diag.Severity == diagnostics.SeverityError {
				hasErrors = true
			}
		}
		if hasErrors {
			os.Exit(1)
		}
	}

	semaDiags := sema.Analyze(program)
	hasSemaErrors := false
	if len(semaDiags) > 0 {
		for _, d := range semaDiags {
			diagPath, diagLine, diagColumn := resolveDiagnosticLocation(d.Pos, path, sourceMap)
			fmt.Fprintf(os.Stderr, "%s:%d:%d: %s: %s\n", diagPath, diagLine, diagColumn, d.Severity, d.Message)
			if d.Severity == diagnostics.SeverityError {
				hasSemaErrors = true
			}
		}
	}
	if hasSemaErrors {
		os.Exit(1)
	}

	irProgram, lowerDiags := lowering.Lower(program)
	if len(lowerDiags) > 0 {
		for _, d := range lowerDiags {
			diagPath, diagLine, diagColumn := resolveDiagnosticLocation(d.Pos, path, sourceMap)
			fmt.Fprintf(os.Stderr, "%s:%d:%d: %s: %s\n", diagPath, diagLine, diagColumn, d.Severity, d.Message)
		}
		os.Exit(1)
	}

	if *showTokens {
		for filePath, fileTokens := range tokensByFile {
			fmt.Printf("tokens for %s\n", filePath)
			for _, tok := range fileTokens {
				fmt.Printf("%d:%d %-14s %q\n", tok.Pos.Line, tok.Pos.Column, tok.Kind, tok.Lexeme)
			}
		}
	}

	if *traceImports {
		printImportTrace(os.Stdout, importTrace)
	}

	totalTokens := 0
	for _, fileTokens := range tokensByFile {
		totalTokens += len(fileTokens)
	}
	fmt.Printf("javec: lowering succeeded (%d tokens, %d sequences)\n", totalTokens, len(program.Sequences))

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

func resolveSponsorMode(modeText string, redacted bool, quiet bool) (sponsor.Mode, error) {
	selected, err := sponsor.ParseMode(modeText)
	if err != nil {
		return "", err
	}
	if quiet {
		if selected != sponsor.ModeFull {
			return "", fmt.Errorf("--sponsor-quiet cannot be combined with --sponsor-notice %q", selected)
		}
		return sponsor.ModeOff, nil
	}
	if redacted {
		if selected != sponsor.ModeFull {
			return "", fmt.Errorf("--sponsor-redacted cannot be combined with --sponsor-notice %q", selected)
		}
		return sponsor.ModeRedacted, nil
	}
	return selected, nil
}

type fileDiagnostic struct {
	path string
	diag diagnostics.Diagnostic
}

type loadOptions struct {
	ProjectRoot string
}

type importResolution struct {
	Importer string
	From     string
	Resolved string
}

type importTrace struct {
	LoadOrder   []string
	Resolutions []importResolution
}

const sourceLineStride = 1_000_000

type sourceIndex struct {
	idByPath map[string]int
	pathByID map[int]string
}

func newSourceIndex() sourceIndex {
	return sourceIndex{
		idByPath: map[string]int{},
		pathByID: map[int]string{},
	}
}

func (s *sourceIndex) register(path string) int {
	if id, ok := s.idByPath[path]; ok {
		return id
	}
	id := len(s.idByPath) + 1
	s.idByPath[path] = id
	s.pathByID[id] = path
	return id
}

func (s sourceIndex) encodePos(path string, pos token.Position) token.Position {
	id, ok := s.idByPath[path]
	if !ok {
		return pos
	}
	if pos.Line <= 0 {
		return pos
	}
	return token.Position{Line: id*sourceLineStride + pos.Line, Column: pos.Column}
}

func (s sourceIndex) decode(pos diagnostics.Position) (string, int, int, bool) {
	if pos.Line <= sourceLineStride {
		return "", 0, 0, false
	}
	id := pos.Line / sourceLineStride
	line := pos.Line % sourceLineStride
	if line == 0 {
		return "", 0, 0, false
	}
	path, ok := s.pathByID[id]
	if !ok {
		return "", 0, 0, false
	}
	return path, line, pos.Column, true
}

func resolveDiagnosticLocation(pos diagnostics.Position, defaultPath string, sources sourceIndex) (string, int, int) {
	if path, line, col, ok := sources.decode(pos); ok {
		return path, line, col
	}
	return defaultPath, pos.Line, pos.Column
}

func loadProgramWithImports(entryPath string, opts loadOptions) (*ast.Program, map[string][]token.Token, importTrace, []fileDiagnostic, sourceIndex, error) {
	visited := map[string]bool{}
	active := map[string]int{}
	stack := make([]string, 0)
	sources := newSourceIndex()
	tokensByFile := map[string][]token.Token{}
	sequences := make([]ast.SequenceDecl, 0)
	rootImports := make([]ast.ImportDecl, 0)
	diags := make([]fileDiagnostic, 0)
	trace := importTrace{LoadOrder: make([]string, 0), Resolutions: make([]importResolution, 0)}

	var visit func(path string, isRoot bool) error
	visit = func(path string, isRoot bool) error {
		abs, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		abs = filepath.Clean(abs)
		sources.register(abs)
		if cycleStart, ok := active[abs]; ok {
			return fmt.Errorf("import cycle detected: %s", formatImportCycle(stack, cycleStart, abs))
		}
		if visited[abs] {
			return nil
		}
		visited[abs] = true
		active[abs] = len(stack)
		stack = append(stack, abs)
		defer func() {
			delete(active, abs)
			stack = stack[:len(stack)-1]
		}()

		trace.LoadOrder = append(trace.LoadOrder, abs)

		src, err := os.ReadFile(abs)
		if err != nil {
			return err
		}

		tokens, lexDiags := lexer.Lex(string(src))
		tokensByFile[abs] = tokens
		for _, d := range lexDiags {
			diags = append(diags, fileDiagnostic{path: abs, diag: d})
		}
		if len(lexDiags) > 0 {
			return nil
		}

		program, parseDiags := parser.Parse(tokens)
		for _, d := range parseDiags {
			diags = append(diags, fileDiagnostic{path: abs, diag: d})
		}
		if len(parseDiags) > 0 {
			return nil
		}

		rewriteProgramPositions(program, abs, sources)

		for _, imp := range program.Imports {
			resolved, err := resolveImportPath(abs, imp.From, opts.ProjectRoot)
			if err != nil {
				return err
			}
			trace.Resolutions = append(trace.Resolutions, importResolution{Importer: abs, From: imp.From, Resolved: resolved})
			if err := visit(resolved, false); err != nil {
				return err
			}
		}

		if isRoot {
			rootImports = append(rootImports, program.Imports...)
		}
		sequences = append(sequences, program.Sequences...)
		return nil
	}

	if err := visit(entryPath, true); err != nil {
		return nil, nil, importTrace{}, nil, sourceIndex{}, err
	}

	merged := &ast.Program{Imports: rootImports, Sequences: sequences}
	return merged, tokensByFile, trace, diags, sources, nil
}

func rewriteProgramPositions(program *ast.Program, path string, sources sourceIndex) {
	for i := range program.Imports {
		program.Imports[i].Pos = sources.encodePos(path, program.Imports[i].Pos)
	}
	for i := range program.Sequences {
		program.Sequences[i].Pos = sources.encodePos(path, program.Sequences[i].Pos)
		program.Sequences[i].Body = rewriteStmtPositions(program.Sequences[i].Body, path, sources)
	}
}

func rewriteStmtPositions(stmts []ast.Stmt, path string, sources sourceIndex) []ast.Stmt {
	for i, stmt := range stmts {
		switch s := stmt.(type) {
		case ast.GiveStmt:
			s.Pos = sources.encodePos(path, s.Pos)
			s.Value = rewriteExprPosition(s.Value, path, sources)
			stmts[i] = s
		case ast.VarDeclStmt:
			s.Pos = sources.encodePos(path, s.Pos)
			s.Value = rewriteExprPosition(s.Value, path, sources)
			stmts[i] = s
		case ast.ExprStmt:
			s.Pos = sources.encodePos(path, s.Pos)
			s.Expr = rewriteExprPosition(s.Expr, path, sources)
			stmts[i] = s
		case ast.AssignmentStmt:
			s.Pos = sources.encodePos(path, s.Pos)
			s.Value = rewriteExprPosition(s.Value, path, sources)
			stmts[i] = s
		case ast.IfStmt:
			s.Pos = sources.encodePos(path, s.Pos)
			for bi := range s.Branches {
				s.Branches[bi].Condition = rewriteExprPosition(s.Branches[bi].Condition, path, sources)
				s.Branches[bi].Body = rewriteStmtPositions(s.Branches[bi].Body, path, sources)
			}
			s.ElseBody = rewriteStmtPositions(s.ElseBody, path, sources)
			stmts[i] = s
		case ast.GivenStmt:
			s.Pos = sources.encodePos(path, s.Pos)
			s.Cond = rewriteExprPosition(s.Cond, path, sources)
			s.In = rewriteExprPosition(s.In, path, sources)
			if s.Init != nil {
				rewritten := rewriteStmtPositions([]ast.Stmt{s.Init}, path, sources)
				s.Init = rewritten[0]
			}
			if s.Step != nil {
				rewritten := rewriteStmtPositions([]ast.Stmt{s.Step}, path, sources)
				s.Step = rewritten[0]
			}
			s.Body = rewriteStmtPositions(s.Body, path, sources)
			stmts[i] = s
		}
	}
	return stmts
}

func rewriteExprPosition(expr ast.Expr, path string, sources sourceIndex) ast.Expr {
	if expr == nil {
		return nil
	}
	switch e := expr.(type) {
	case ast.IdentifierExpr:
		e.Pos = sources.encodePos(path, e.Pos)
		return e
	case ast.MemberExpr:
		e.Target = rewriteExprPosition(e.Target, path, sources)
		return e
	case ast.IndexExpr:
		e.Target = rewriteExprPosition(e.Target, path, sources)
		e.Index = rewriteExprPosition(e.Index, path, sources)
		return e
	case ast.CallExpr:
		e.Callee = rewriteExprPosition(e.Callee, path, sources)
		for i, a := range e.Args {
			e.Args[i] = rewriteExprPosition(a, path, sources)
		}
		return e
	case ast.BinaryExpr:
		e.Left = rewriteExprPosition(e.Left, path, sources)
		e.Right = rewriteExprPosition(e.Right, path, sources)
		return e
	case ast.CollectionLiteralExpr:
		for i, a := range e.Items {
			e.Items[i] = rewriteExprPosition(a, path, sources)
		}
		for i, p := range e.Pairs {
			p.Key = rewriteExprPosition(p.Key, path, sources)
			p.Value = rewriteExprPosition(p.Value, path, sources)
			e.Pairs[i] = p
		}
		return e
	}
	return expr
}

func printImportTrace(out io.Writer, trace importTrace) {
	fmt.Fprintln(out, "javec: import trace")
	for i, p := range trace.LoadOrder {
		fmt.Fprintf(out, "  load[%d]: %s\n", i, p)
	}
	for _, r := range trace.Resolutions {
		fmt.Fprintf(out, "  resolve: %s -> %s (%s)\n", r.Importer, r.Resolved, r.From)
	}
}
func resolveImportPath(currentFile string, from string, projectRoot string) (string, error) {
	if strings.HasPrefix(from, "highschool/") {
		return resolveHighschoolImportPath(currentFile, from, projectRoot)
	}

	base := filepath.Dir(currentFile)
	candidate := filepath.Clean(filepath.Join(base, filepath.FromSlash(from)))

	if filepath.Ext(candidate) == ".jave" {
		found := existingCandidates([]string{candidate})
		if len(found) == 1 {
			return found[0], nil
		}
		return "", fmt.Errorf("import not found: %s (tried: %s)", from, candidate)
	}

	searchCandidates := []string{
		candidate + ".jave",
		filepath.Join(candidate, "main.jave"),
	}
	found := existingCandidates(searchCandidates)
	if len(found) == 1 {
		return found[0], nil
	}
	if len(found) > 1 {
		return "", fmt.Errorf("ambiguous import path %s (matched: %s)", from, strings.Join(found, ", "))
	}

	return "", fmt.Errorf("import not found: %s (tried: %s)", from, strings.Join(searchCandidates, ", "))
}

func resolveHighschoolImportPath(currentFile string, from string, projectRoot string) (string, error) {
	rel := strings.TrimPrefix(from, "highschool/")
	if rel == "" {
		return "", fmt.Errorf("invalid highschool import path: %s", from)
	}

	resolvedRoot := ""
	if trimmed := strings.TrimSpace(projectRoot); trimmed != "" {
		candidate := filepath.Clean(trimmed)
		if abs, err := filepath.Abs(candidate); err == nil {
			resolvedRoot = abs
		} else {
			resolvedRoot = candidate
		}
	}
	if resolvedRoot == "" {
		resolvedRoot = findWorkspaceRoot(currentFile)
		if resolvedRoot == "" {
			resolvedRoot = filepath.Dir(currentFile)
		}
	}

	relPath := filepath.FromSlash(rel)
	candidates := []string{
		filepath.Join(resolvedRoot, "highschool", relPath, "main.jave"),
		filepath.Join(resolvedRoot, "highschool", relPath+".jave"),
		filepath.Join(resolvedRoot, "stdlib", "highschool", relPath, "main.jave"),
		filepath.Join(resolvedRoot, "stdlib", "highschool", relPath+".jave"),
	}

	found := existingCandidates(candidates)
	if len(found) == 1 {
		return found[0], nil
	}
	if len(found) > 1 {
		return "", fmt.Errorf("ambiguous highschool import %s (matched: %s)", from, strings.Join(found, ", "))
	}

	return "", fmt.Errorf("highschool import not found: %s (tried: %s)", from, strings.Join(candidates, ", "))
}

func existingCandidates(candidates []string) []string {
	found := make([]string, 0, len(candidates))
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			found = append(found, filepath.Clean(c))
		}
	}
	return found
}

func findWorkspaceRoot(currentFile string) string {
	dir := filepath.Dir(currentFile)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

func formatImportCycle(stack []string, cycleStart int, repeated string) string {
	parts := make([]string, 0, len(stack)-cycleStart+1)
	for i := cycleStart; i < len(stack); i++ {
		parts = append(parts, stack[i])
	}
	parts = append(parts, repeated)
	return strings.Join(parts, " -> ")
}
