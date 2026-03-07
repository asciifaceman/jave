package ast

import "github.com/asciifaceman/jave/internal/token"

// Program is the root AST node.
type Program struct {
	Imports   []ImportDecl
	Sequences []SequenceDecl
}

// ImportDecl represents `install X from path;;`.
type ImportDecl struct {
	Pos  token.Position
	Name string
	From string
}

// SequenceDecl represents a top-level sequence declaration.
type SequenceDecl struct {
	Pos          token.Position
	Visibility   string
	SourceModule string
	Name         string
	Params       []SequenceParam
	ReturnType   string
	Doc          *Docstring
	Body         []Stmt
}

// Docstring is structured documentation metadata attached to declarations.
type Docstring struct {
	Pos      token.Position
	Raw      string
	Sections []DocSection
}

// DocSection is a labeled block inside a docstring.
type DocSection struct {
	Label string
	Body  string
}

// SequenceParam represents one typed parameter in a sequence signature.
type SequenceParam struct {
	TypeName string
	Name     string
	Variadic bool
}

// Stmt is a statement node.
type Stmt interface {
	isStmt()
}

// Expr is an expression node.
type Expr interface {
	isExpr()
}

// GiveStmt models `give ... up;;`.
type GiveStmt struct {
	Pos   token.Position
	Value Expr // nil for bare `give up;;`
}

func (GiveStmt) isStmt() {}

// VarDeclStmt models `allow T Name 2b=2 Expr;;`.
type VarDeclStmt struct {
	Pos      token.Position
	TypeName string
	Name     string
	Value    Expr
}

func (VarDeclStmt) isStmt() {}

// ExprStmt wraps a standalone expression statement.
type ExprStmt struct {
	Pos  token.Position
	Expr Expr
}

func (ExprStmt) isStmt() {}

// AssignmentStmt models `Name 2b=2 Expr;;`.
type AssignmentStmt struct {
	Pos   token.Position
	Name  string
	Value Expr
}

func (AssignmentStmt) isStmt() {}

// IfBranch captures one maybe/furthermore branch.
type IfBranch struct {
	Condition Expr
	Body      []Stmt
}

// IfStmt models `maybe`/`furthermore`/`otherwise`.
type IfStmt struct {
	Pos      token.Position
	Branches []IfBranch
	ElseBody []Stmt
}

func (IfStmt) isStmt() {}

// GivenStmt models all `given` loop forms.
type GivenStmt struct {
	Pos    token.Position
	Mode   string // while-ish, for-ish, within
	Header string
	Cond   Expr
	Init   Stmt
	Step   Stmt
	Var    string
	In     Expr
	Body   []Stmt
}

func (GivenStmt) isStmt() {}

// IdentifierExpr is an identifier reference.
type IdentifierExpr struct {
	Pos  token.Position
	Name string
}

func (IdentifierExpr) isExpr() {}

// NumberExpr is an integer/float literal.
type NumberExpr struct {
	Value string
}

func (NumberExpr) isExpr() {}

// StringExpr is a string literal.
type StringExpr struct {
	Value string
}

func (StringExpr) isExpr() {}

// BoolExpr is a `yee` or `nee` literal.
type BoolExpr struct {
	Value bool
}

func (BoolExpr) isExpr() {}

// CollectionLiteralExpr models list-like and map-like literals.
type CollectionLiteralExpr struct {
	Form  string // table, enumeration, lexis
	Items []Expr
	Pairs []KeyValueExpr
}

func (CollectionLiteralExpr) isExpr() {}

// KeyValueExpr models a key/value entry inside lexis literals.
type KeyValueExpr struct {
	Key   Expr
	Value Expr
}

// MemberExpr models `Target.Name`.
type MemberExpr struct {
	Target Expr
	Name   string
}

func (MemberExpr) isExpr() {}

// IndexExpr models `Target[Index]`.
type IndexExpr struct {
	Target Expr
	Index  Expr
}

func (IndexExpr) isExpr() {}

// CallExpr models both paren and angle argument calls.
type CallExpr struct {
	Callee Expr
	Args   []Expr
	Style  string // "paren" or "angle"
}

func (CallExpr) isExpr() {}

// BinaryExpr models infix arithmetic.
type BinaryExpr struct {
	Left  Expr
	Op    string
	Right Expr
}

func (BinaryExpr) isExpr() {}
