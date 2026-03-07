package lowering

import (
	"github.com/asciifaceman/jave/internal/ast"
	"github.com/asciifaceman/jave/internal/diagnostics"
	"github.com/asciifaceman/jave/internal/ir"
)

// Lower converts a parsed+validated AST program into executable IR.
func Lower(program *ast.Program) (*ir.ProgramIR, []diagnostics.Diagnostic) {
	l := &lowerer{}
	return l.lower(program), l.diags
}

type lowerer struct {
	diags []diagnostics.Diagnostic
}

func (l *lowerer) lower(program *ast.Program) *ir.ProgramIR {
	forewards := make([]ast.SequenceDecl, 0)
	var foremost *ast.SequenceDecl
	for i := range program.Sequences {
		if program.Sequences[i].Name == "Foreward" {
			forewards = append(forewards, program.Sequences[i])
		}
		if program.Sequences[i].Name == "Foremost" {
			foremost = &program.Sequences[i]
		}
	}
	if foremost == nil {
		l.errorAt(1, 1, "cannot lower program without Foremost sequence")
		return nil
	}

	out := &ir.ProgramIR{
		Foremost: ir.SequenceIR{
			Name:         foremost.Name,
			Params:       namesFromParams(foremost.Params),
			ReturnType:   foremost.ReturnType,
			Instructions: make([]ir.Instruction, 0, len(foremost.Body)),
		},
		Sequences: map[string]ir.SequenceIR{},
	}
	if len(forewards) > 0 {
		out.Forewards = make([]ir.SequenceIR, 0, len(forewards))
		for _, f := range forewards {
			lowered := ir.SequenceIR{
				Name:         f.Name,
				Params:       namesFromParams(f.Params),
				ReturnType:   f.ReturnType,
				Instructions: l.lowerStatements(f.Body),
			}
			out.Forewards = append(out.Forewards, lowered)
			out.Sequences[f.Name] = lowered
		}
	}

	out.Foremost.Instructions = l.lowerStatements(foremost.Body)
	out.Sequences[out.Foremost.Name] = out.Foremost

	for _, seq := range program.Sequences {
		if _, exists := out.Sequences[seq.Name]; exists {
			continue
		}
		out.Sequences[seq.Name] = ir.SequenceIR{
			Name:         seq.Name,
			Params:       namesFromParams(seq.Params),
			ReturnType:   seq.ReturnType,
			Instructions: l.lowerStatements(seq.Body),
		}
	}

	return out
}

func namesFromParams(params []ast.SequenceParam) []string {
	if len(params) == 0 {
		return nil
	}
	out := make([]string, 0, len(params))
	for _, p := range params {
		out = append(out, p.Name)
	}
	return out
}

func (l *lowerer) lowerStatements(stmts []ast.Stmt) []ir.Instruction {
	out := make([]ir.Instruction, 0, len(stmts))
	for _, stmt := range stmts {
		switch s := stmt.(type) {
		case ast.VarDeclStmt:
			out = append(out, ir.VarDeclInstr{Pos: s.Pos, Name: s.Name, Value: s.Value})
		case ast.AssignmentStmt:
			out = append(out, ir.AssignInstr{Pos: s.Pos, Name: s.Name, Value: s.Value})
		case ast.ExprStmt:
			out = append(out, ir.ExprInstr{Pos: s.Pos, Expr: s.Expr})
		case ast.GiveStmt:
			out = append(out, ir.ReturnInstr{Pos: s.Pos, Value: s.Value})
		case ast.IfStmt:
			branches := make([]ir.IfBranchIR, 0, len(s.Branches))
			for _, b := range s.Branches {
				branches = append(branches, ir.IfBranchIR{
					Condition: b.Condition,
					Body:      l.lowerStatements(b.Body),
				})
			}
			out = append(out, ir.IfInstr{
				Pos:      s.Pos,
				Branches: branches,
				ElseBody: l.lowerStatements(s.ElseBody),
			})
		case ast.GivenStmt:
			switch s.Mode {
			case "while-ish":
				if s.Cond == nil {
					l.errorAt(s.Pos.Line, s.Pos.Column, "while-ish given loop is missing condition")
					continue
				}
				out = append(out, ir.WhileInstr{
					Pos:       s.Pos,
					Condition: s.Cond,
					Body:      l.lowerStatements(s.Body),
				})
			case "for-ish":
				if s.Cond == nil || s.Init == nil || s.Step == nil {
					l.errorAt(s.Pos.Line, s.Pos.Column, "for-ish given loop is missing init/condition/step")
					continue
				}
				out = append(out, ir.ForInstr{
					Pos:       s.Pos,
					Init:      l.lowerStatements([]ast.Stmt{s.Init}),
					Condition: s.Cond,
					Step:      l.lowerStatements([]ast.Stmt{s.Step}),
					Body:      l.lowerStatements(s.Body),
				})
			case "within":
				if s.Var == "" || s.In == nil {
					l.errorAt(s.Pos.Line, s.Pos.Column, "within given loop is missing variable/iterable")
					continue
				}
				out = append(out, ir.WithinInstr{
					Pos:      s.Pos,
					VarName:  s.Var,
					Iterable: s.In,
					Body:     l.lowerStatements(s.Body),
				})
			default:
				l.errorAt(s.Pos.Line, s.Pos.Column, "unsupported given loop mode: "+s.Mode)
			}
		default:
			pos := stmtPos(stmt)
			l.errorAt(pos.Line, pos.Column, "lowering does not support this statement yet")
		}
	}
	return out
}

func (l *lowerer) errorAt(line, column int, message string) {
	l.diags = append(l.diags, diagnostics.Diagnostic{
		Severity: diagnostics.SeverityError,
		Message:  message,
		Pos: diagnostics.Position{
			Line:   line,
			Column: column,
		},
	})
}

func stmtPos(stmt ast.Stmt) (pos struct{ Line, Column int }) {
	pos.Line = 1
	pos.Column = 1
	switch s := stmt.(type) {
	case ast.GiveStmt:
		pos.Line, pos.Column = s.Pos.Line, s.Pos.Column
	case ast.VarDeclStmt:
		pos.Line, pos.Column = s.Pos.Line, s.Pos.Column
	case ast.AssignmentStmt:
		pos.Line, pos.Column = s.Pos.Line, s.Pos.Column
	case ast.ExprStmt:
		pos.Line, pos.Column = s.Pos.Line, s.Pos.Column
	case ast.IfStmt:
		pos.Line, pos.Column = s.Pos.Line, s.Pos.Column
	case ast.GivenStmt:
		pos.Line, pos.Column = s.Pos.Line, s.Pos.Column
	}
	return pos
}
