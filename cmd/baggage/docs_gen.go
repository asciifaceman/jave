package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/asciifaceman/jave/internal/ast"
	"github.com/asciifaceman/jave/internal/diagnostics"
	"github.com/asciifaceman/jave/internal/lexer"
	"github.com/asciifaceman/jave/internal/parser"
	"gopkg.in/yaml.v3"
)

type docsOptions struct {
	projectRoot string
	outDir      string
	manifestDir string
}

type docsEntry struct {
	Origin     string
	Kind       string
	Name       string
	Aliases    []string
	Carryon    string
	Signature  string
	Visibility string
	ImportHint string
	Title      string
	About      string
	Notes      []string
	Warnings   []string
	Examples   []string
	Links      []string
	SeeAlso    []string
	Params     []seqParamDoc
	ReturnType string
	ReturnDesc string
	Extras     []labeledText
}

type seqParamDoc struct {
	Name        string
	TypeName    string
	Description string
	Variadic    bool
}

type labeledText struct {
	Label string
	Body  string
}

type manifestDoc struct {
	Kind      string   `yaml:"kind"`
	Name      string   `yaml:"name"`
	Aliases   []string `yaml:"aliases"`
	Signature string   `yaml:"signature"`
	Carryon   string   `yaml:"carryon"`
	Title     string   `yaml:"title"`
	About     string   `yaml:"about"`
	Notes     []string `yaml:"notes"`
	Warnings  []string `yaml:"warnings"`
	Examples  []string `yaml:"examples"`
	Links     []string `yaml:"links"`
	SeeAlso   []string `yaml:"seealso"`
	Status    string   `yaml:"status"`
}

func parseDocsArgs(args []string) (docsOptions, error) {
	opts := docsOptions{}
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch a {
		case "--project-root":
			if i+1 >= len(args) {
				return docsOptions{}, fmt.Errorf("missing value for --project-root")
			}
			i++
			opts.projectRoot = args[i]
		case "--out-dir":
			if i+1 >= len(args) {
				return docsOptions{}, fmt.Errorf("missing value for --out-dir")
			}
			i++
			opts.outDir = args[i]
		case "--manifest-dir":
			if i+1 >= len(args) {
				return docsOptions{}, fmt.Errorf("missing value for --manifest-dir")
			}
			i++
			opts.manifestDir = args[i]
		default:
			return docsOptions{}, fmt.Errorf("unknown flag %s", a)
		}
	}

	if opts.projectRoot == "" {
		opts.projectRoot = defaultProjectRoot()
	}
	if opts.outDir == "" {
		opts.outDir = filepath.Join(opts.projectRoot, "site", "reference")
	} else if !filepath.IsAbs(opts.outDir) {
		opts.outDir = filepath.Join(opts.projectRoot, opts.outDir)
	}
	if opts.manifestDir == "" {
		opts.manifestDir = filepath.Join(opts.projectRoot, "docs-manifests")
	} else if !filepath.IsAbs(opts.manifestDir) {
		opts.manifestDir = filepath.Join(opts.projectRoot, opts.manifestDir)
	}
	return opts, nil
}

func runDocs(opts docsOptions) error {
	sourceEntries, err := collectSourceDocs(opts.projectRoot)
	if err != nil {
		return err
	}
	manifestEntries, err := collectManifestDocs(opts.manifestDir)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(opts.outDir, 0o755); err != nil {
		return err
	}
	carryonDir := filepath.Join(opts.outDir, "carryons")
	if err := os.MkdirAll(carryonDir, 0o755); err != nil {
		return err
	}

	if err := writeCarryonPages(carryonDir, sourceEntries); err != nil {
		return err
	}
	if err := writeManifestPages(opts.outDir, manifestEntries); err != nil {
		return err
	}
	if err := writeReferenceIndex(opts.outDir, sourceEntries, manifestEntries); err != nil {
		return err
	}
	fmt.Printf("baggage: generated docs in %s\n", opts.outDir)
	return nil
}

func collectSourceDocs(projectRoot string) ([]docsEntry, error) {
	highschoolDir := filepath.Join(projectRoot, "highschool")
	entries := make([]docsEntry, 0)
	err := filepath.WalkDir(highschoolDir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".jave" {
			return nil
		}
		seqEntries, err := sourceDocsFromFile(projectRoot, path)
		if err != nil {
			return err
		}
		entries = append(entries, seqEntries...)
		return nil
	})
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Carryon == entries[j].Carryon {
			return entries[i].Name < entries[j].Name
		}
		return entries[i].Carryon < entries[j].Carryon
	})
	return entries, nil
}

func sourceDocsFromFile(projectRoot, path string) ([]docsEntry, error) {
	srcBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	toks, lexDiags := lexer.Lex(string(srcBytes))
	if len(lexDiags) != 0 {
		return nil, fmt.Errorf("lexer diagnostics in %s: %s", path, lexDiags[0].Message)
	}
	prog, parseDiags := parser.Parse(toks)
	for _, d := range parseDiags {
		if d.Severity == diagnostics.SeverityError {
			return nil, fmt.Errorf("parser diagnostics in %s: %s", path, d.Message)
		}
	}

	relPath, err := filepath.Rel(projectRoot, path)
	if err != nil {
		return nil, err
	}
	carryonPath := filepath.ToSlash(filepath.Dir(relPath))
	alias := defaultImportAlias(carryonPath)

	entries := make([]docsEntry, 0, len(prog.Sequences))
	for _, seq := range prog.Sequences {
		if seq.Visibility != "outy" {
			continue
		}
		entry := docsEntry{
			Origin:     "source",
			Kind:       "sequence",
			Name:       seq.Name,
			Carryon:    carryonPath,
			Signature:  formatSequenceSignature(seq),
			Visibility: seq.Visibility,
			ImportHint: fmt.Sprintf("install %s from %s;;", alias, carryonPath),
			ReturnType: seq.ReturnType,
			Params:     paramsFromSequence(seq),
		}
		if seq.Doc != nil {
			applyDocSections(&entry, seq.Doc.Sections)
		}
		if entry.Title == "" {
			entry.Title = seq.Name
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func defaultImportAlias(carryonPath string) string {
	canonical := map[string]string{
		"highschool/English":        "Strangs",
		"highschool/Communications": "Pronts",
	}
	if alias, ok := canonical[carryonPath]; ok {
		return alias
	}
	return filepath.Base(carryonPath)
}

func paramsFromSequence(seq ast.SequenceDecl) []seqParamDoc {
	params := make([]seqParamDoc, 0, len(seq.Params))
	for _, p := range seq.Params {
		params = append(params, seqParamDoc{Name: p.Name, TypeName: p.TypeName, Variadic: p.Variadic})
	}
	return params
}

func formatSequenceSignature(seq ast.SequenceDecl) string {
	parts := make([]string, 0, len(seq.Params))
	for _, p := range seq.Params {
		prefix := ""
		if p.Variadic {
			prefix = "..."
		}
		parts = append(parts, fmt.Sprintf("%s%s %s", prefix, p.TypeName, p.Name))
	}
	return fmt.Sprintf("%s seq %s<%s> --> <<%s>>", seq.Visibility, seq.Name, strings.Join(parts, ", "), seq.ReturnType)
}

func applyDocSections(entry *docsEntry, sections []ast.DocSection) {
	paramDescs := map[string]string{}
	for _, section := range sections {
		label := strings.ToLower(strings.TrimSpace(section.Label))
		body := strings.TrimSpace(section.Body)
		switch label {
		case "title":
			entry.Title = body
		case "about":
			entry.About = body
		case "note":
			entry.Notes = append(entry.Notes, body)
		case "warning":
			entry.Warnings = append(entry.Warnings, body)
		case "example", "examples":
			entry.Examples = append(entry.Examples, body)
		case "link", "links":
			entry.Links = append(entry.Links, body)
		case "seealso":
			entry.SeeAlso = append(entry.SeeAlso, body)
		case "return":
			if entry.ReturnDesc == "" {
				entry.ReturnDesc = body
			} else {
				entry.Extras = append(entry.Extras, labeledText{Label: section.Label, Body: body})
			}
		case "param":
			name, desc, ok := splitParamBody(body)
			if ok {
				paramDescs[name] = desc
			} else {
				entry.Extras = append(entry.Extras, labeledText{Label: section.Label, Body: body})
			}
		default:
			entry.Extras = append(entry.Extras, labeledText{Label: section.Label, Body: body})
		}
	}
	for i := range entry.Params {
		if desc, ok := paramDescs[entry.Params[i].Name]; ok {
			entry.Params[i].Description = desc
		}
	}
}

func splitParamBody(body string) (string, string, bool) {
	idx := strings.Index(body, ":")
	if idx <= 0 {
		return "", "", false
	}
	name := strings.TrimSpace(body[:idx])
	desc := strings.TrimSpace(body[idx+1:])
	if name == "" || desc == "" {
		return "", "", false
	}
	return name, desc, true
}

func collectManifestDocs(manifestDir string) ([]docsEntry, error) {
	entries := make([]docsEntry, 0)
	if _, err := os.Stat(manifestDir); err != nil {
		if os.IsNotExist(err) {
			return entries, nil
		}
		return nil, err
	}
	err := filepath.WalkDir(manifestDir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".yaml" && ext != ".yml" {
			return nil
		}
		entry, err := readManifestDoc(path)
		if err != nil {
			return err
		}
		entries = append(entries, entry)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Kind == entries[j].Kind {
			return entries[i].Name < entries[j].Name
		}
		return entries[i].Kind < entries[j].Kind
	})
	return entries, nil
}

func readManifestDoc(path string) (docsEntry, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return docsEntry{}, err
	}
	var m manifestDoc
	if err := yaml.Unmarshal(b, &m); err != nil {
		return docsEntry{}, fmt.Errorf("parse manifest %s: %w", path, err)
	}
	if m.Kind == "" || m.Name == "" {
		return docsEntry{}, fmt.Errorf("manifest %s missing required kind/name", path)
	}
	entry := docsEntry{
		Origin:    "manifest",
		Kind:      m.Kind,
		Name:      m.Name,
		Aliases:   m.Aliases,
		Carryon:   m.Carryon,
		Signature: m.Signature,
		Title:     firstNonEmpty(m.Title, m.Name),
		About:     m.About,
		Notes:     m.Notes,
		Warnings:  m.Warnings,
		Examples:  m.Examples,
		Links:     m.Links,
		SeeAlso:   m.SeeAlso,
	}
	if m.Status != "" {
		entry.Extras = append(entry.Extras, labeledText{Label: "Status", Body: m.Status})
	}
	return entry, nil
}

func firstNonEmpty(a, b string) string {
	if strings.TrimSpace(a) != "" {
		return a
	}
	return b
}

func writeCarryonPages(outDir string, entries []docsEntry) error {
	byCarryon := map[string][]docsEntry{}
	for _, e := range entries {
		byCarryon[e.Carryon] = append(byCarryon[e.Carryon], e)
	}
	carryons := make([]string, 0, len(byCarryon))
	for carryon := range byCarryon {
		carryons = append(carryons, carryon)
	}
	sort.Strings(carryons)

	for _, carryon := range carryons {
		symbols := byCarryon[carryon]
		sort.Slice(symbols, func(i, j int) bool { return symbols[i].Name < symbols[j].Name })
		slug := slugify(carryon)
		outPath := filepath.Join(outDir, slug+".md")
		body := renderCarryonPageBody(carryon, symbols)
		front := renderFrontMatter(carryon+" Reference", "/reference/carryons/"+slug+"/", "carryon-reference", "source")
		if err := os.WriteFile(outPath, []byte(front+body), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func writeManifestPages(outDir string, entries []docsEntry) error {
	builtins := make([]docsEntry, 0)
	language := make([]docsEntry, 0)
	for _, e := range entries {
		switch e.Kind {
		case "builtin":
			builtins = append(builtins, e)
		case "language-feature":
			language = append(language, e)
		}
	}
	if len(builtins) > 0 {
		body := renderManifestPageBody("Builtin Reference", builtins)
		front := renderFrontMatter("Builtin Reference", "/reference/builtins/", "builtin-reference", "manifest")
		if err := os.WriteFile(filepath.Join(outDir, "builtins.md"), []byte(front+body), 0o644); err != nil {
			return err
		}
	}
	if len(language) > 0 {
		body := renderManifestPageBody("Language Feature Reference", language)
		front := renderFrontMatter("Language Feature Reference", "/reference/language/", "language-reference", "manifest")
		if err := os.WriteFile(filepath.Join(outDir, "language.md"), []byte(front+body), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func writeReferenceIndex(outDir string, sourceEntries, manifestEntries []docsEntry) error {
	carryons := map[string]struct{}{}
	for _, e := range sourceEntries {
		carryons[e.Carryon] = struct{}{}
	}
	carryonList := make([]string, 0, len(carryons))
	for c := range carryons {
		carryonList = append(carryonList, c)
	}
	sort.Strings(carryonList)

	var b strings.Builder
	b.WriteString("# Jave Reference\n\n")
	b.WriteString("Godoc-style reference pages generated from source docstrings and YAML manifests.\n\n")
	b.WriteString("## Carryons\n\n")
	for _, c := range carryonList {
		b.WriteString(fmt.Sprintf("- [%s](./carryons/%s)\n", c, slugify(c)))
	}
	b.WriteString("\n## Core References\n\n")
	if hasKind(manifestEntries, "builtin") {
		b.WriteString("- [Builtins](./builtins)\n")
	}
	if hasKind(manifestEntries, "language-feature") {
		b.WriteString("- [Language Features](./language)\n")
	}
	front := renderFrontMatter("Jave Reference", "/reference/", "reference-index", "mixed")
	return os.WriteFile(filepath.Join(outDir, "index.md"), []byte(front+b.String()), 0o644)
}

func hasKind(entries []docsEntry, kind string) bool {
	for _, e := range entries {
		if e.Kind == kind {
			return true
		}
	}
	return false
}

func renderFrontMatter(title, permalink, category, sourceKind string) string {
	return fmt.Sprintf("---\ntitle: %s\npermalink: %s\ncategory: %s\nsource_kind: %s\n---\n\n", title, permalink, category, sourceKind)
}

func renderCarryonPageBody(carryon string, symbols []docsEntry) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("# %s\n\n", carryon))
	if len(symbols) > 0 && symbols[0].ImportHint != "" {
		b.WriteString("```jave\n")
		b.WriteString(symbols[0].ImportHint)
		b.WriteString("\n```\n\n")
	}
	b.WriteString("## Table of Contents\n\n")
	for _, s := range symbols {
		anchor := slugify(s.Name)
		b.WriteString(fmt.Sprintf("- [`%s`](#%s)\n", s.Name, anchor))
	}
	b.WriteString("\n")
	for _, s := range symbols {
		renderEntry(&b, s)
	}
	return b.String()
}

func renderManifestPageBody(title string, entries []docsEntry) string {
	var b strings.Builder
	b.WriteString("# " + title + "\n\n")
	b.WriteString("## Table of Contents\n\n")
	for _, e := range entries {
		b.WriteString(fmt.Sprintf("- [`%s`](#%s)\n", e.Name, slugify(e.Name)))
	}
	b.WriteString("\n")
	for _, e := range entries {
		renderEntry(&b, e)
	}
	return b.String()
}

func renderEntry(b *strings.Builder, e docsEntry) {
	b.WriteString(fmt.Sprintf("## %s\n\n", e.Name))
	if e.Title != "" && e.Title != e.Name {
		b.WriteString(e.Title + "\n\n")
	}
	if e.Signature != "" {
		b.WriteString("```jave\n")
		b.WriteString(e.Signature)
		b.WriteString("\n```\n\n")
	}
	if e.About != "" {
		b.WriteString(e.About + "\n\n")
	}
	if len(e.Params) > 0 {
		b.WriteString("### Parameters\n\n")
		for _, p := range e.Params {
			prefix := ""
			if p.Variadic {
				prefix = "..."
			}
			desc := p.Description
			if desc == "" {
				desc = "No description provided."
			}
			b.WriteString(fmt.Sprintf("- `%s%s %s`: %s\n", prefix, p.TypeName, p.Name, desc))
		}
		b.WriteString("\n")
	}
	if e.ReturnType != "" {
		b.WriteString("### Returns\n\n")
		if e.ReturnDesc != "" {
			b.WriteString(fmt.Sprintf("- `%s`: %s\n\n", e.ReturnType, e.ReturnDesc))
		} else {
			b.WriteString(fmt.Sprintf("- `%s`\n\n", e.ReturnType))
		}
	}
	renderSimpleSection(b, "Notes", e.Notes)
	renderSimpleSection(b, "Warnings", e.Warnings)
	if len(e.Examples) > 0 {
		b.WriteString("### Examples\n\n")
		for _, ex := range e.Examples {
			b.WriteString("```jave\n")
			b.WriteString(ex)
			b.WriteString("\n```\n\n")
		}
	}
	renderSimpleSection(b, "Links", e.Links)
	renderSimpleSection(b, "See Also", e.SeeAlso)
	for _, extra := range e.Extras {
		b.WriteString(fmt.Sprintf("### %s\n\n%s\n\n", extra.Label, extra.Body))
	}
}

func renderSimpleSection(b *strings.Builder, heading string, items []string) {
	if len(items) == 0 {
		return
	}
	b.WriteString("### " + heading + "\n\n")
	for _, item := range items {
		b.WriteString("- " + item + "\n")
	}
	b.WriteString("\n")
}

func slugify(v string) string {
	v = strings.ToLower(v)
	var b strings.Builder
	lastDash := false
	for _, r := range v {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			b.WriteRune('-')
			lastDash = true
		}
	}
	out := strings.Trim(b.String(), "-")
	if out == "" {
		return "entry"
	}
	return out
}
