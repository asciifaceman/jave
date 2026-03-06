package diagnostics

// Severity represents diagnostic severity.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
)

// Position represents a source position.
type Position struct {
	Line   int
	Column int
}

// Diagnostic is a user-facing compiler message.
type Diagnostic struct {
	Severity Severity
	Message  string
	Pos      Position
}
