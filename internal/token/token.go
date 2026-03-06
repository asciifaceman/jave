package token

// Kind identifies a lexical token kind.
type Kind string

const (
	Illegal Kind = "Illegal"
	EOF     Kind = "EOF"

	Identifier Kind = "Identifier"
	Integer    Kind = "Integer"
	Float      Kind = "Float"
	String     Kind = "String"

	Outy        Kind = "Outy"
	Inny        Kind = "Inny"
	Sequence    Kind = "Sequence"
	Seq         Kind = "Seq"
	Give        Kind = "Give"
	Up          Kind = "Up"
	Maybe       Kind = "Maybe"
	Furthermore Kind = "Furthermore"
	Otherwise   Kind = "Otherwise"
	Given       Kind = "Given"
	Again       Kind = "Again"
	Allow       Kind = "Allow"
	Install     Kind = "Install"
	From        Kind = "From"
	Within      Kind = "Within"
	Yee         Kind = "Yee"
	Nee         Kind = "Nee"

	TypeExact   Kind = "TypeExact"
	TypeVag     Kind = "TypeVag"
	TypeTruther Kind = "TypeTruther"
	TypeStrang  Kind = "TypeStrang"
	TypeNada    Kind = "TypeNada"
	TypeNaw     Kind = "TypeNaw"

	Assign         Kind = "Assign"
	StmtEnd        Kind = "StmtEnd"
	ReturnArrow    Kind = "ReturnArrow"
	ThinArrow      Kind = "ThinArrow"
	DoubleLAngle   Kind = "DoubleLAngle"
	DoubleRAngle   Kind = "DoubleRAngle"
	LParen         Kind = "LParen"
	RParen         Kind = "RParen"
	LBrace         Kind = "LBrace"
	RBrace         Kind = "RBrace"
	LBracket       Kind = "LBracket"
	RBracket       Kind = "RBracket"
	LAngle         Kind = "LAngle"
	RAngle         Kind = "RAngle"
	Comma          Kind = "Comma"
	Dot            Kind = "Dot"
	Colon          Kind = "Colon"
	Plus           Kind = "Plus"
	Minus          Kind = "Minus"
	Star           Kind = "Star"
	Slash          Kind = "Slash"
	Percent        Kind = "Percent"
)

// Position represents a 1-based line/column source location.
type Position struct {
	Line   int
	Column int
}

// Token is a single lexical unit.
type Token struct {
	Kind   Kind
	Lexeme string
	Pos    Position
}
