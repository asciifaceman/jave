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
	diags             []diagnostics.Diagnostic
	globalSequences   map[string]struct{}
	sequenceArities   map[string]map[int]struct{}
	sequenceVariadics map[string]int
	imports           map[string]struct{}
	moduleArities     map[string]map[string]map[int]struct{}
	moduleVariadics   map[string]map[string]int
}

func (a *analyzer) analyzeProgram(program *ast.Program) {
	a.globalSequences = map[string]struct{}{}
	a.sequenceArities = map[string]map[int]struct{}{}
	a.sequenceVariadics = map[string]int{}
	a.imports = map[string]struct{}{}
	a.moduleArities = map[string]map[string]map[int]struct{}{}
	a.moduleVariadics = map[string]map[string]int{}

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

	seenLocal := map[string]bool{}
	seenLocalByArity := map[string]map[int]bool{}
	seenModule := map[string]map[string]map[int]bool{}
	foremostCount := 0

	for _, seq := range program.Sequences {
		if seq.SourceModule != "" {
			moduleSet := a.moduleArities[seq.SourceModule]
			if moduleSet == nil {
				moduleSet = map[string]map[int]struct{}{}
				a.moduleArities[seq.SourceModule] = moduleSet
			}
			nameSet := moduleSet[seq.Name]
			if nameSet == nil {
				nameSet = map[int]struct{}{}
				moduleSet[seq.Name] = nameSet
			}
			if len(seq.Params) > 0 && seq.Params[len(seq.Params)-1].Variadic {
				moduleVar := a.moduleVariadics[seq.SourceModule]
				if moduleVar == nil {
					moduleVar = map[string]int{}
					a.moduleVariadics[seq.SourceModule] = moduleVar
				}
				moduleVar[seq.Name] = len(seq.Params) - 1
			} else {
				nameSet[len(seq.Params)] = struct{}{}
			}
		}

		if seq.Name != "Foreward" {
			if seq.SourceModule == "" {
				aritySet := seenLocalByArity[seq.Name]
				if aritySet == nil {
					aritySet = map[int]bool{}
					seenLocalByArity[seq.Name] = aritySet
				}
				declArity := len(seq.Params)
				if len(seq.Params) > 0 && seq.Params[len(seq.Params)-1].Variadic {
					declArity = len(seq.Params) - 1
				}
				if aritySet[declArity] {
					a.errorAt(seq.Pos, "duplicate sequence declaration: "+seq.Name)
				}
				aritySet[declArity] = true
				seenLocal[seq.Name] = true
				a.globalSequences[seq.Name] = struct{}{}
				if len(seq.Params) > 0 && seq.Params[len(seq.Params)-1].Variadic {
					a.sequenceVariadics[seq.Name] = len(seq.Params) - 1
				} else {
					globalArities := a.sequenceArities[seq.Name]
					if globalArities == nil {
						globalArities = map[int]struct{}{}
						a.sequenceArities[seq.Name] = globalArities
					}
					globalArities[len(seq.Params)] = struct{}{}
				}
			} else {
				moduleSeen := seenModule[seq.SourceModule]
				if moduleSeen == nil {
					moduleSeen = map[string]map[int]bool{}
					seenModule[seq.SourceModule] = moduleSeen
				}
				aritySeen := moduleSeen[seq.Name]
				if aritySeen == nil {
					aritySeen = map[int]bool{}
					moduleSeen[seq.Name] = aritySeen
				}
				declArity := len(seq.Params)
				if len(seq.Params) > 0 && seq.Params[len(seq.Params)-1].Variadic {
					declArity = len(seq.Params) - 1
				}
				if aritySeen[declArity] {
					a.errorAt(seq.Pos, "duplicate sequence declaration: "+seq.SourceModule+"."+seq.Name)
				}
				aritySeen[declArity] = true
			}
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
		seqScope := a.newSequenceScope(seq.SourceModule)
		for _, param := range seq.Params {
			if seqScope.hasHere(param.Name) {
				a.errorAt(seq.Pos, "duplicate sequence parameter: "+param.Name)
				continue
			}
			seqScope.define(param.Name)
		}
		a.checkStatements(seq.ReturnType, seq.Body, seqScope, seq.SourceModule)
	}
}

func (a *analyzer) checkStatements(returnType string, statements []ast.Stmt, scope *scope, currentModule string) {
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
				a.checkExpr(s.Value, scope, currentModule)
			}
		case ast.VarDeclStmt:
			if scope.hasHere(s.Name) {
				a.errorAt(s.Pos, "duplicate local declaration: "+s.Name)
			}
			a.checkExpr(s.Value, scope, currentModule)
			scope.define(s.Name)
		case ast.AssignmentStmt:
			if !scope.has(s.Name) {
				a.errorAt(s.Pos, "assignment to undefined identifier: "+s.Name)
			}
			a.checkExpr(s.Value, scope, currentModule)
		case ast.ExprStmt:
			a.checkExpr(s.Expr, scope, currentModule)
		case ast.IfStmt:
			for _, b := range s.Branches {
				a.checkExpr(b.Condition, scope, currentModule)
				a.checkStatements(returnType, b.Body, scope.child(), currentModule)
			}
			a.checkStatements(returnType, s.ElseBody, scope.child(), currentModule)
		case ast.GivenStmt:
			loopScope := scope.child()
			if s.Init != nil {
				a.checkStatements(returnType, []ast.Stmt{s.Init}, loopScope, currentModule)
			}
			if s.Cond != nil {
				a.checkExpr(s.Cond, loopScope, currentModule)
			}
			if s.In != nil {
				a.checkExpr(s.In, loopScope, currentModule)
			}
			if s.Var != "" {
				if loopScope.hasHere(s.Var) {
					a.errorAt(s.Pos, "duplicate loop variable declaration: "+s.Var)
				} else {
					loopScope.define(s.Var)
				}
			}
			a.checkStatements(returnType, s.Body, loopScope.child(), currentModule)
			if s.Step != nil {
				a.checkStatements(returnType, []ast.Stmt{s.Step}, loopScope, currentModule)
			}
		}
	}
}

func (a *analyzer) checkExpr(expr ast.Expr, scope *scope, currentModule string) {
	switch e := expr.(type) {
	case ast.IdentifierExpr:
		if !scope.has(e.Name) {
			a.errorAt(e.Pos, "undefined identifier: "+e.Name)
		}
	case ast.BinaryExpr:
		a.checkExpr(e.Left, scope, currentModule)
		a.checkExpr(e.Right, scope, currentModule)
	case ast.MemberExpr:
		a.checkExpr(e.Target, scope, currentModule)
	case ast.IndexExpr:
		a.checkExpr(e.Target, scope, currentModule)
		a.checkExpr(e.Index, scope, currentModule)
	case ast.CallExpr:
		if callee, ok := e.Callee.(ast.IdentifierExpr); ok {
			argCount := len(e.Args)
			if currentModule != "" {
				if moduleSet, exists := a.moduleArities[currentModule]; exists {
					if arities, moduleExists := moduleSet[callee.Name]; moduleExists {
						if _, arityExists := arities[argCount]; !arityExists {
							moduleVar, hasVar := a.moduleVariadics[currentModule]
							minArgs, varExists := 0, false
							if hasVar {
								minArgs, varExists = moduleVar[callee.Name]
							}
							if !varExists || minArgs > argCount {
								a.errorAt(exprPos(e), "sequence call arity mismatch for "+currentModule+"."+callee.Name)
							}
						}
					}
				}
				if moduleVar, hasVar := a.moduleVariadics[currentModule]; hasVar {
					if minArgs, varExists := moduleVar[callee.Name]; varExists && argCount < minArgs {
						a.errorAt(exprPos(e), "sequence call arity mismatch for "+currentModule+"."+callee.Name)
					}
				}
			}
			if arities, exists := a.sequenceArities[callee.Name]; exists {
				if _, arityExists := arities[argCount]; !arityExists {
					if minArgs, hasVar := a.sequenceVariadics[callee.Name]; !hasVar || argCount < minArgs {
						a.errorAt(exprPos(e), "sequence call arity mismatch for "+callee.Name)
					}
				}
			} else if minArgs, hasVar := a.sequenceVariadics[callee.Name]; hasVar && argCount < minArgs {
				a.errorAt(exprPos(e), "sequence call arity mismatch for "+callee.Name)
			}
		}
		if callee, ok := e.Callee.(ast.MemberExpr); ok {
			if module, ok := callee.Target.(ast.IdentifierExpr); ok {
				argCount := len(e.Args)
				if moduleSet, exists := a.moduleArities[module.Name]; exists {
					arities, memberExists := moduleSet[callee.Name]
					if !memberExists {
						moduleVar, hasVar := a.moduleVariadics[module.Name]
						if !hasVar {
							a.errorAt(exprPos(e), "undefined module sequence: "+module.Name+"."+callee.Name)
						} else if _, varExists := moduleVar[callee.Name]; !varExists {
							a.errorAt(exprPos(e), "undefined module sequence: "+module.Name+"."+callee.Name)
						}
					} else if _, arityExists := arities[argCount]; !arityExists {
						moduleVar, hasVar := a.moduleVariadics[module.Name]
						minArgs, varExists := 0, false
						if hasVar {
							minArgs, varExists = moduleVar[callee.Name]
						}
						if !varExists || minArgs > argCount {
							a.errorAt(exprPos(e), "sequence call arity mismatch for "+module.Name+"."+callee.Name)
						}
					}
				} else if moduleVar, hasVar := a.moduleVariadics[module.Name]; hasVar {
					if minArgs, varExists := moduleVar[callee.Name]; varExists {
						if argCount < minArgs {
							a.errorAt(exprPos(e), "sequence call arity mismatch for "+module.Name+"."+callee.Name)
						}
					} else {
						a.errorAt(exprPos(e), "undefined module sequence: "+module.Name+"."+callee.Name)
					}
				}
			}
		}
		a.checkExpr(e.Callee, scope, currentModule)
		for _, arg := range e.Args {
			a.checkExpr(arg, scope, currentModule)
		}
	case ast.CollectionLiteralExpr:
		for _, item := range e.Items {
			a.checkExpr(item, scope, currentModule)
		}
		for _, pair := range e.Pairs {
			a.checkExpr(pair.Key, scope, currentModule)
			a.checkExpr(pair.Value, scope, currentModule)
		}
	}
}

func (a *analyzer) newSequenceScope(currentModule string) *scope {
	s := newScope(nil)
	for _, builtin := range []string{"pront", "prontulate", "girth", "slotify", "Strangs"} {
		s.define(builtin)
	}
	if currentModule != "" {
		if moduleSet, exists := a.moduleArities[currentModule]; exists {
			for name := range moduleSet {
				s.define(name)
			}
		}
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
