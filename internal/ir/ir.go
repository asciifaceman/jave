package ir

import (
	"github.com/asciifaceman/jave/internal/ast"
	"github.com/asciifaceman/jave/internal/token"
)

// ProgramIR is a lowered representation ready for runtime execution.
type ProgramIR struct {
	Forewards               []SequenceIR
	Foremost                SequenceIR
	Sequences               map[string]SequenceIR
	ModuleSequences         map[string]map[string]SequenceIR
	SequenceOverloads       map[string]map[int]SequenceIR
	ModuleSequenceOverloads map[string]map[string]map[int]SequenceIR
	SequenceVariadics       map[string]SequenceIR
	ModuleSequenceVariadics map[string]map[string]SequenceIR
}

// SequenceIR contains executable instructions for one sequence.
type SequenceIR struct {
	Name         string
	Module       string
	Params       []string
	Variadic     bool
	FixedParams  int
	ReturnType   string
	Instructions []Instruction
}

// Instruction is an executable IR instruction.
type Instruction interface {
	isInstruction()
}

// VarDeclInstr declares and initializes a variable.
type VarDeclInstr struct {
	Pos   token.Position
	Name  string
	Value ast.Expr
}

func (VarDeclInstr) isInstruction() {}

// AssignInstr updates an existing variable.
type AssignInstr struct {
	Pos   token.Position
	Name  string
	Value ast.Expr
}

func (AssignInstr) isInstruction() {}

// ExprInstr evaluates an expression for side effects.
type ExprInstr struct {
	Pos  token.Position
	Expr ast.Expr
}

func (ExprInstr) isInstruction() {}

// ReturnInstr exits sequence execution.
type ReturnInstr struct {
	Pos   token.Position
	Value ast.Expr
}

func (ReturnInstr) isInstruction() {}

// IfBranchIR contains a branch condition and its body.
type IfBranchIR struct {
	Condition ast.Expr
	Body      []Instruction
}

// IfInstr models maybe/furthermore/otherwise control flow.
type IfInstr struct {
	Pos      token.Position
	Branches []IfBranchIR
	ElseBody []Instruction
}

func (IfInstr) isInstruction() {}

// WhileInstr models the while-ish given loop form.
type WhileInstr struct {
	Pos       token.Position
	Condition ast.Expr
	Body      []Instruction
}

func (WhileInstr) isInstruction() {}

// ForInstr models the for-ish given loop form.
type ForInstr struct {
	Pos       token.Position
	Init      []Instruction
	Condition ast.Expr
	Step      []Instruction
	Body      []Instruction
}

func (ForInstr) isInstruction() {}

// WithinInstr models the collection iteration form `given (<Name within Expr>)`.
type WithinInstr struct {
	Pos      token.Position
	VarName  string
	Iterable ast.Expr
	Body     []Instruction
}

func (WithinInstr) isInstruction() {}
