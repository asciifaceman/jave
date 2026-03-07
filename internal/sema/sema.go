package sema

import (
	"strings"

	"github.com/asciifaceman/jave/internal/ast"
	"github.com/asciifaceman/jave/internal/diagnostics"
	"github.com/asciifaceman/jave/internal/token"
)

// Analyze runs semantic checks over the parsed program.
func Analyze(program *ast.Program) []diagnostics.Diagnostic {
	a := analyzer{}
	a.analyzeProgram(program)
	return a.diags
}

type analyzer struct {
	diags           []diagnostics.Diagnostic
	globalSequences map[string]struct{}
	sequenceArity   map[string]int
	imports         map[string]struct{}
}

func (a *analyzer) analyzeProgram(program *ast.Program) {
	a.globalSequences = map[string]struct{}{}
	a.sequenceArity = map[string]int{}
	a.imports = map[string]struct{}{}

	for _, imp := range program.Imports {
		if _, exists := a.imports[imp.Name]; exists {
			a.errorAt(imp.Pos, "duplicate import declaration: "+imp.Name)
			continue
		}

		if imp.Name == "Srangs" {
			a.warnAt(imp.Pos, "legacy module alias 'Srangs' remains supported for ecosystem continuity")
		}

		if (imp.Name == "Strangs" || imp.Name == "Srangs" || imp.Name == "Pronts") && !strings.HasPrefix(imp.From, "highschool/") {
			a.errorAt(imp.Pos, "standard library import '"+imp.Name+"' must use highschool/... path")
		}

		a.imports[imp.Name] = struct{}{}
	}

	seen := map[string]bool{}
	foremostCount := 0

	for _, seq := range program.Sequences {
		if seq.Name != "Foreward" && seen[seq.Name] {
			a.errorAt(seq.Pos, "duplicate sequence declaration: "+seq.Name)
		}
		if seq.Name != "Foreward" {
			seen[seq.Name] = true
			a.globalSequences[seq.Name] = struct{}{}
			a.sequenceArity[seq.Name] = len(seq.Params)
		}

		if seq.Name == "Foremost" {
			foremostCount++
			if seq.ReturnType != "nada" {
				a.errorAt(seq.Pos, "Foremost must return nada in v0.1")
			}
		}
	}

	if foremostCount == 0 {
		a.error("program is missing required Foremost sequence")
	}

	for _, seq := range program.Sequences {
		seqScope := a.newSequenceScope()
		for _, param := range seq.Params {
			if seqScope.hasHere(param.Name) {
				a.errorAt(seq.Pos, "duplicate sequence parameter: "+param.Name)
				continue
			}
			seqScope.define(param.Name)
		}
		a.checkStatements(seq.ReturnType, seq.Body, seqScope)
	}
}

func (a *analyzer) checkStatements(returnType string, statements []ast.Stmt, scope *scope) {
	for _, stmt := range statements {
		switch s := stmt.(type) {
		case ast.GiveStmt:
			hasValue := s.Value != nil
			if returnType == "nada" && hasValue {
				a.errorAt(s.Pos, "nada sequence cannot return a value")
			}
			if returnType != "nada" && !hasValue {
				a.errorAt(s.Pos, "non-nada sequence must return a value")
			}
			if s.Value != nil {
				a.checkExpr(s.Value, scope)
			}
		case ast.VarDeclStmt:
			if scope.hasHere(s.Name) {
				a.errorAt(s.Pos, "duplicate local declaration: "+s.Name)
			}
			a.checkExpr(s.Value, scope)
			scope.define(s.Name)
		case ast.AssignmentStmt:
			if !scope.has(s.Name) {
				a.errorAt(s.Pos, "assignment to undefined identifier: "+s.Name)
			}
			a.checkExpr(s.Value, scope)
		case ast.ExprStmt:
			a.checkExpr(s.Expr, scope)
		case ast.IfStmt:
			for _, b := range s.Branches {
				a.checkExpr(b.Condition, scope)
				a.checkStatements(returnType, b.Body, scope.child())
			}
			a.checkStatements(returnType, s.ElseBody, scope.child())
		case ast.GivenStmt:
			loopScope := scope.child()
			if s.Init != nil {
				a.checkStatements(returnType, []ast.Stmt{s.Init}, loopScope)
			}
			if s.Cond != nil {
				a.checkExpr(s.Cond, loopScope)
			}
			if s.In != nil {
				a.checkExpr(s.In, loopScope)
			}
			if s.Var != "" {
				if loopScope.hasHere(s.Var) {
					a.errorAt(s.Pos, "duplicate loop variable declaration: "+s.Var)
				} else {
					loopScope.define(s.Var)
				}
			}
			a.checkStatements(returnType, s.Body, loopScope.child())
			if s.Step != nil {
				a.checkStatements(returnType, []ast.Stmt{s.Step}, loopScope)
			}
		}
	}
}

func (a *analyzer) checkExpr(expr ast.Expr, scope *scope) {
	switch e := expr.(type) {
	case ast.IdentifierExpr:
		if !scope.has(e.Name) {
			a.errorAt(e.Pos, "undefined identifier: "+e.Name)
		}
	case ast.BinaryExpr:
		a.checkExpr(e.Left, scope)
		a.checkExpr(e.Right, scope)
	case ast.MemberExpr:
		a.checkExpr(e.Target, scope)
	case ast.IndexExpr:
		a.checkExpr(e.Target, scope)
		a.checkExpr(e.Index, scope)
	case ast.CallExpr:
		if callee, ok := e.Callee.(ast.IdentifierExpr); ok {
			if expected, exists := a.sequenceArity[callee.Name]; exists && expected != len(e.Args) {
				a.errorAt(exprPos(e), "sequence call arity mismatch for "+callee.Name)
			}
		}
		a.checkExpr(e.Callee, scope)
		for _, arg := range e.Args {
			a.checkExpr(arg, scope)
		}
	case ast.CollectionLiteralExpr:
		for _, item := range e.Items {
			a.checkExpr(item, scope)
		}
		for _, pair := range e.Pairs {
			a.checkExpr(pair.Key, scope)
			a.checkExpr(pair.Value, scope)
		}
	}
}

func (a *analyzer) newSequenceScope() *scope {
	s := newScope(nil)
	for _, builtin := range []string{"pront", "girth", "Strangs", "Pronts"} {
		s.define(builtin)
	}
	for name := range a.imports {
		s.define(name)
	}
	for name := range a.globalSequences {
		s.define(name)
	}
	return s
}

func (a *analyzer) error(message string) {
	a.errorAt(token.Position{Line: 1, Column: 1}, message)
}

func (a *analyzer) errorAt(pos token.Position, message string) {
	a.diags = append(a.diags, diagnostics.Diagnostic{
		Severity: diagnostics.SeverityError,
		Message:  message,
		Pos: diagnostics.Position{
			Line:   pos.Line,
			Column: pos.Column,
		},
	})
}

func (a *analyzer) warnAt(pos token.Position, message string) {
	a.diags = append(a.diags, diagnostics.Diagnostic{
		Severity: diagnostics.SeverityWarning,
		Message:  message,
		Pos: diagnostics.Position{
			Line:   pos.Line,
			Column: pos.Column,
		},
	})
}

func exprPos(expr ast.Expr) token.Position {
	pos := token.Position{Line: 1, Column: 1}
	switch e := expr.(type) {
	case ast.IdentifierExpr:
		pos = e.Pos
	case ast.BinaryExpr:
		pos = exprPos(e.Left)
	case ast.MemberExpr:
		pos = exprPos(e.Target)
	case ast.IndexExpr:
		pos = exprPos(e.Target)
	case ast.CallExpr:
		pos = exprPos(e.Callee)
	}
	return pos
}

type scope struct {
	parent *scope
	names  map[string]struct{}
}

func newScope(parent *scope) *scope {
	return &scope{parent: parent, names: map[string]struct{}{}}
}

func (s *scope) child() *scope {
	return newScope(s)
}

func (s *scope) define(name string) {
	s.names[name] = struct{}{}
}

func (s *scope) hasHere(name string) bool {
	_, ok := s.names[name]
	return ok
}

func (s *scope) has(name string) bool {
	if _, ok := s.names[name]; ok {
		return true
	}
	if s.parent != nil {
		return s.parent.has(name)
	}
	return false
}
