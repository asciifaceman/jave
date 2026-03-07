package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/asciifaceman/jave/internal/lsp"
)

type rpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type rpcResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *rpcError   `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type server struct {
	docs        map[string]string
	stdlib      lsp.Index
	manifests   lsp.Index
	initialized bool
	shutdown    bool
}

type initializeResult struct {
	Capabilities serverCapabilities `json:"capabilities"`
	ServerInfo   serverInfo         `json:"serverInfo"`
}

type serverCapabilities struct {
	TextDocumentSync int                   `json:"textDocumentSync"`
	HoverProvider    bool                  `json:"hoverProvider"`
	SignatureHelp    signatureHelpProvider `json:"signatureHelpProvider"`
}

type signatureHelpProvider struct {
	TriggerCharacters []string `json:"triggerCharacters"`
}

type serverInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type didOpenParams struct {
	TextDocument textDocumentItem `json:"textDocument"`
}

type textDocumentItem struct {
	URI  string `json:"uri"`
	Text string `json:"text"`
}

type didChangeParams struct {
	TextDocument   versionedTextDocumentIdentifier `json:"textDocument"`
	ContentChanges []textContentChangeEvent        `json:"contentChanges"`
}

type versionedTextDocumentIdentifier struct {
	URI string `json:"uri"`
}

type textContentChangeEvent struct {
	Text string `json:"text"`
}

type textDocumentPositionParams struct {
	TextDocument textDocumentIdentifier `json:"textDocument"`
	Position     position               `json:"position"`
}

type textDocumentIdentifier struct {
	URI string `json:"uri"`
}

type position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

type hoverResult struct {
	Contents markupContent `json:"contents"`
}

type markupContent struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

type signatureHelpResult struct {
	Signatures      []signatureInformation `json:"signatures"`
	ActiveSignature int                    `json:"activeSignature"`
	ActiveParameter int                    `json:"activeParameter"`
}

type signatureInformation struct {
	Label         string                 `json:"label"`
	Documentation *markupContent         `json:"documentation,omitempty"`
	Parameters    []parameterInformation `json:"parameters"`
}

type parameterInformation struct {
	Label         string         `json:"label"`
	Documentation *markupContent `json:"documentation,omitempty"`
}

func main() {
	if len(os.Args) > 1 {
		a := os.Args[1]
		if a == "--version" || a == "version" {
			fmt.Println("javels v0.1.0")
			return
		}
	}

	projectRoot := defaultProjectRoot()
	stdlib, _ := lsp.BuildIndexFromHighschool(projectRoot)
	manifests, _ := lsp.BuildIndexFromManifests(filepath.Join(projectRoot, "docs-manifests"))

	s := &server{
		docs:      map[string]string{},
		stdlib:    stdlib,
		manifests: manifests,
	}
	if err := s.serve(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "javels: %v\n", err)
		os.Exit(1)
	}
}

func defaultProjectRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return wd
}

func (s *server) serve(in io.Reader, out io.Writer) error {
	reader := bufio.NewReader(in)
	writer := bufio.NewWriter(out)
	defer writer.Flush()

	for {
		payload, err := readMessage(reader)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		var req rpcRequest
		if err := json.Unmarshal(payload, &req); err != nil {
			continue
		}
		if req.Method == "exit" {
			if s.shutdown {
				return nil
			}
			return fmt.Errorf("exit before shutdown")
		}

		resp := s.handle(req)
		if resp == nil {
			continue
		}
		if err := writeMessage(writer, resp); err != nil {
			return err
		}
	}
}

func (s *server) handle(req rpcRequest) *rpcResponse {
	id := rawID(req.ID)
	respond := len(req.ID) > 0

	switch req.Method {
	case "initialize":
		s.initialized = true
		if !respond {
			return nil
		}
		return &rpcResponse{
			JSONRPC: "2.0",
			ID:      id,
			Result: initializeResult{
				Capabilities: serverCapabilities{
					TextDocumentSync: 1,
					HoverProvider:    true,
					SignatureHelp: signatureHelpProvider{
						TriggerCharacters: []string{"<", "(", ","},
					},
				},
				ServerInfo: serverInfo{Name: "javels", Version: "0.1.0"},
			},
		}
	case "shutdown":
		s.shutdown = true
		if !respond {
			return nil
		}
		return &rpcResponse{JSONRPC: "2.0", ID: id, Result: nil}
	case "textDocument/didOpen":
		var p didOpenParams
		_ = json.Unmarshal(req.Params, &p)
		s.docs[p.TextDocument.URI] = p.TextDocument.Text
		return nil
	case "textDocument/didChange":
		var p didChangeParams
		_ = json.Unmarshal(req.Params, &p)
		if len(p.ContentChanges) > 0 {
			s.docs[p.TextDocument.URI] = p.ContentChanges[len(p.ContentChanges)-1].Text
		}
		return nil
	case "textDocument/hover":
		if !respond {
			return nil
		}
		var p textDocumentPositionParams
		_ = json.Unmarshal(req.Params, &p)
		res, ok := s.hover(p)
		if !ok {
			return &rpcResponse{JSONRPC: "2.0", ID: id, Result: nil}
		}
		return &rpcResponse{JSONRPC: "2.0", ID: id, Result: res}
	case "textDocument/signatureHelp":
		if !respond {
			return nil
		}
		var p textDocumentPositionParams
		_ = json.Unmarshal(req.Params, &p)
		res, ok := s.signatureHelp(p)
		if !ok {
			return &rpcResponse{JSONRPC: "2.0", ID: id, Result: nil}
		}
		return &rpcResponse{JSONRPC: "2.0", ID: id, Result: res}
	default:
		if !respond {
			return nil
		}
		return &rpcResponse{JSONRPC: "2.0", ID: id, Error: &rpcError{Code: -32601, Message: "method not found"}}
	}
}

func (s *server) hover(p textDocumentPositionParams) (hoverResult, bool) {
	source, ok := s.sourceForURI(p.TextDocument.URI)
	if !ok {
		return hoverResult{}, false
	}
	symbol := symbolAtPosition(source, p.Position)
	if symbol == "" {
		return hoverResult{}, false
	}
	doc, ok := s.lookupSymbol(symbol, source, p.TextDocument.URI)
	if !ok {
		return hoverResult{}, false
	}
	return hoverResult{Contents: markupContent{Kind: "markdown", Value: renderHover(doc)}}, true
}

func (s *server) signatureHelp(p textDocumentPositionParams) (signatureHelpResult, bool) {
	source, ok := s.sourceForURI(p.TextDocument.URI)
	if !ok {
		return signatureHelpResult{}, false
	}
	name, argIndex := callContextAtPosition(source, p.Position)
	if name == "" {
		return signatureHelpResult{}, false
	}
	doc, ok := s.lookupSymbol(name, source, p.TextDocument.URI)
	if !ok {
		return signatureHelpResult{}, false
	}
	sig := signatureInformation{
		Label:         doc.Signature,
		Documentation: &markupContent{Kind: "markdown", Value: doc.About},
		Parameters:    make([]parameterInformation, 0, len(doc.Params)),
	}
	for _, p := range doc.Params {
		prefix := ""
		if p.Variadic {
			prefix = "..."
		}
		param := parameterInformation{Label: fmt.Sprintf("%s%s %s", prefix, p.TypeName, p.Name)}
		if p.Description != "" {
			param.Documentation = &markupContent{Kind: "markdown", Value: p.Description}
		}
		sig.Parameters = append(sig.Parameters, param)
	}
	if argIndex < 0 {
		argIndex = 0
	}
	if len(sig.Parameters) > 0 && argIndex >= len(sig.Parameters) {
		argIndex = len(sig.Parameters) - 1
	}
	return signatureHelpResult{Signatures: []signatureInformation{sig}, ActiveSignature: 0, ActiveParameter: argIndex}, true
}

func (s *server) lookupSymbol(symbol, source, uri string) (lsp.SymbolDoc, bool) {
	local, err := lsp.BuildIndexFromSource(source, uri)
	if err == nil {
		if d, ok := local.Symbols[symbol]; ok {
			return d, true
		}
	}
	if d, ok := s.stdlib.Symbols[symbol]; ok {
		return d, true
	}
	if d, ok := s.manifests.Symbols[symbol]; ok {
		return d, true
	}
	return lsp.SymbolDoc{}, false
}

func (s *server) sourceForURI(uri string) (string, bool) {
	if src, ok := s.docs[uri]; ok {
		return src, true
	}
	path, err := uriToPath(uri)
	if err != nil {
		return "", false
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}
	return string(b), true
}

func renderHover(doc lsp.SymbolDoc) string {
	var b strings.Builder
	b.WriteString("### " + doc.Name + "\n\n")
	if doc.Signature != "" {
		b.WriteString("```jave\n")
		b.WriteString(doc.Signature)
		b.WriteString("\n```\n\n")
	}
	if doc.Title != "" && doc.Title != doc.Name {
		b.WriteString(doc.Title + "\n\n")
	}
	if doc.About != "" {
		b.WriteString(doc.About + "\n\n")
	}
	if len(doc.Params) > 0 {
		b.WriteString("**Parameters**\n")
		for _, p := range doc.Params {
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
	if doc.ReturnType != "" {
		if doc.ReturnDesc != "" {
			b.WriteString(fmt.Sprintf("**Returns:** `%s` - %s\n", doc.ReturnType, doc.ReturnDesc))
		} else {
			b.WriteString(fmt.Sprintf("**Returns:** `%s`\n", doc.ReturnType))
		}
	}
	return b.String()
}

func symbolAtPosition(source string, pos position) string {
	lines := strings.Split(source, "\n")
	if pos.Line < 0 || pos.Line >= len(lines) {
		return ""
	}
	line := lines[pos.Line]
	if len(line) == 0 {
		return ""
	}
	idx := pos.Character
	if idx >= len(line) {
		idx = len(line) - 1
	}
	if idx < 0 {
		return ""
	}

	if !isIdentChar(line[idx]) {
		if idx > 0 && isIdentChar(line[idx-1]) {
			idx--
		} else {
			return ""
		}
	}

	start := idx
	for start > 0 && isIdentChar(line[start-1]) {
		start--
	}
	end := idx + 1
	for end < len(line) && isIdentChar(line[end]) {
		end++
	}
	if start >= end {
		return ""
	}
	return line[start:end]
}

func callContextAtPosition(source string, pos position) (string, int) {
	lines := strings.Split(source, "\n")
	if pos.Line < 0 || pos.Line >= len(lines) {
		return "", 0
	}
	line := lines[pos.Line]
	cursor := pos.Character
	if cursor > len(line) {
		cursor = len(line)
	}
	left := line[:cursor]
	open := strings.LastIndexAny(left, "<(")
	if open < 0 {
		return "", 0
	}

	nameEnd := open
	nameStart := nameEnd
	for nameStart > 0 && isIdentChar(line[nameStart-1]) {
		nameStart--
	}
	if nameStart == nameEnd {
		return "", 0
	}
	name := line[nameStart:nameEnd]
	argIdx := strings.Count(left[open:cursor], ",")
	return name, argIdx
}

func isIdentChar(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}

func uriToPath(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	if u.Scheme != "file" {
		return "", fmt.Errorf("unsupported uri scheme %q", u.Scheme)
	}
	path, err := url.PathUnescape(u.Path)
	if err != nil {
		return "", err
	}
	if len(path) >= 3 && path[0] == '/' && path[2] == ':' {
		path = path[1:]
	}
	if u.Host != "" && u.Host != "localhost" {
		path = `\\` + u.Host + filepath.FromSlash(path)
	}
	return filepath.FromSlash(path), nil
}

func rawID(v json.RawMessage) interface{} {
	if len(v) == 0 {
		return nil
	}
	var any interface{}
	if err := json.Unmarshal(v, &any); err != nil {
		return nil
	}
	return any
}

func readMessage(r *bufio.Reader) ([]byte, error) {
	contentLength := -1
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			break
		}
		if strings.HasPrefix(strings.ToLower(line), "content-length:") {
			v := strings.TrimSpace(line[len("content-length:"):])
			n, err := strconv.Atoi(v)
			if err != nil {
				return nil, err
			}
			contentLength = n
		}
	}
	if contentLength < 0 {
		return nil, fmt.Errorf("missing Content-Length")
	}
	payload := make([]byte, contentLength)
	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func writeMessage(w *bufio.Writer, v interface{}) error {
	payload, err := json.Marshal(v)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	_, _ = fmt.Fprintf(&buf, "Content-Length: %d\r\n\r\n", len(payload))
	_, _ = buf.Write(payload)
	if _, err := w.Write(buf.Bytes()); err != nil {
		return err
	}
	return w.Flush()
}
