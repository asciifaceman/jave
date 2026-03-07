package lsp

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

// SymbolDoc is the LSP-facing normalized symbol documentation model.
type SymbolDoc struct {
	Name       string
	Kind       string
	Signature  string
	Title      string
	About      string
	Params     []ParamDoc
	ReturnType string
	ReturnDesc string
	Source     string
}

// ParamDoc contains one parameter in signature/help output.
type ParamDoc struct {
	Name        string
	TypeName    string
	Description string
	Variadic    bool
}

// Index stores symbol docs lookup by exact symbol name.
type Index struct {
	Symbols map[string]SymbolDoc
}

// NewIndex creates an empty symbol index.
func NewIndex() Index {
	return Index{Symbols: map[string]SymbolDoc{}}
}

// Merge merges another index into the receiver without overwriting existing keys.
func (i *Index) Merge(other Index) {
	for k, v := range other.Symbols {
		if _, exists := i.Symbols[k]; !exists {
			i.Symbols[k] = v
		}
	}
}

// BuildIndexFromSource parses a source unit and builds sequence docs.
func BuildIndexFromSource(source, sourceName string) (Index, error) {
	idx := NewIndex()
	toks, lexDiags := lexer.Lex(source)
	if len(lexDiags) != 0 {
		return idx, fmt.Errorf("lexer diagnostics: %s", lexDiags[0].Message)
	}
	prog, parseDiags := parser.Parse(toks)
	for _, d := range parseDiags {
		if d.Severity == diagnostics.SeverityError {
			return idx, fmt.Errorf("parser diagnostics: %s", d.Message)
		}
	}
	for _, seq := range prog.Sequences {
		if seq.Visibility != "outy" {
			continue
		}
		doc := symbolDocFromSequence(seq, sourceName)
		idx.Symbols[doc.Name] = doc
	}
	return idx, nil
}

// BuildIndexFromHighschool scans highschool carryons for exported symbol docs.
func BuildIndexFromHighschool(projectRoot string) (Index, error) {
	idx := NewIndex()
	highschool := filepath.Join(projectRoot, "highschool")
	err := filepath.WalkDir(highschool, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() || filepath.Ext(path) != ".jave" {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		fileIdx, err := BuildIndexFromSource(string(b), path)
		if err != nil {
			return err
		}
		idx.Merge(fileIdx)
		return nil
	})
	if err != nil {
		if os.IsNotExist(err) {
			return idx, nil
		}
		return idx, err
	}
	return idx, nil
}

type manifestDoc struct {
	Kind      string   `yaml:"kind"`
	Name      string   `yaml:"name"`
	Signature string   `yaml:"signature"`
	Title     string   `yaml:"title"`
	About     string   `yaml:"about"`
	Examples  []string `yaml:"examples"`
}

// BuildIndexFromManifests loads builtin/language docs from docs-manifests.
func BuildIndexFromManifests(manifestDir string) (Index, error) {
	idx := NewIndex()
	if _, err := os.Stat(manifestDir); err != nil {
		if os.IsNotExist(err) {
			return idx, nil
		}
		return idx, err
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
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var m manifestDoc
		if err := yaml.Unmarshal(b, &m); err != nil {
			return err
		}
		if strings.TrimSpace(m.Name) == "" {
			return nil
		}
		doc := SymbolDoc{
			Name:      m.Name,
			Kind:      firstNonEmpty(m.Kind, "manifest"),
			Signature: m.Signature,
			Title:     firstNonEmpty(m.Title, m.Name),
			About:     m.About,
			Source:    path,
		}
		idx.Symbols[m.Name] = doc
		return nil
	})
	if err != nil {
		return idx, err
	}
	return idx, nil
}

func symbolDocFromSequence(seq ast.SequenceDecl, source string) SymbolDoc {
	doc := SymbolDoc{
		Name:       seq.Name,
		Kind:       "sequence",
		Signature:  formatSequenceSignature(seq),
		Title:      seq.Name,
		ReturnType: seq.ReturnType,
		Params:     make([]ParamDoc, 0, len(seq.Params)),
		Source:     source,
	}
	for _, p := range seq.Params {
		doc.Params = append(doc.Params, ParamDoc{Name: p.Name, TypeName: p.TypeName, Variadic: p.Variadic})
	}
	if seq.Doc != nil {
		applyDocSections(&doc, seq.Doc.Sections)
	}
	return doc
}

func applyDocSections(doc *SymbolDoc, sections []ast.DocSection) {
	paramDesc := map[string]string{}
	for _, section := range sections {
		label := strings.ToLower(strings.TrimSpace(section.Label))
		body := strings.TrimSpace(section.Body)
		switch label {
		case "title":
			doc.Title = body
		case "about":
			doc.About = body
		case "param":
			name, desc, ok := splitParamBody(body)
			if ok {
				paramDesc[name] = desc
			}
		case "return":
			doc.ReturnDesc = body
		}
	}
	for i := range doc.Params {
		if desc, ok := paramDesc[doc.Params[i].Name]; ok {
			doc.Params[i].Description = desc
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

func formatSequenceSignature(seq ast.SequenceDecl) string {
	parts := make([]string, 0, len(seq.Params))
	for _, p := range seq.Params {
		prefix := ""
		if p.Variadic {
			prefix = "..."
		}
		parts = append(parts, fmt.Sprintf("%s%s %s", prefix, p.TypeName, p.Name))
	}
	return fmt.Sprintf("%s<%s> -> %s", seq.Name, strings.Join(parts, ", "), seq.ReturnType)
}

func firstNonEmpty(a, b string) string {
	if strings.TrimSpace(a) != "" {
		return a
	}
	return b
}

// SortedNames returns sorted symbol names.
func (i Index) SortedNames() []string {
	names := make([]string, 0, len(i.Symbols))
	for k := range i.Symbols {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}
