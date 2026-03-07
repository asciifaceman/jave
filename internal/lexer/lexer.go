package lexer

import (
	"strings"
	"unicode"

	"github.com/asciifaceman/jave/internal/diagnostics"
	"github.com/asciifaceman/jave/internal/token"
)

// Lex tokenizes Jave source code and returns tokens with diagnostics.
func Lex(src string) ([]token.Token, []diagnostics.Diagnostic) {
	l := &lexer{
		src:    []rune(src),
		line:   1,
		column: 1,
	}
	return l.lex()
}

type lexer struct {
	src         []rune
	idx         int
	line        int
	column      int
	tokens      []token.Token
	diagnostics []diagnostics.Diagnostic
}

func (l *lexer) lex() ([]token.Token, []diagnostics.Diagnostic) {
	for !l.atEnd() {
		l.skipWhitespace()
		if l.atEnd() {
			break
		}

		start := l.position()
		ch := l.peek()

		switch {
		case ch == '>' && l.wouldMatch(">>|"):
			l.skipLineComment()
			continue
		case ch == '=' && l.wouldMatch("=["):
			l.skipBlockComment(start)
			continue
		case ch == 'd' && l.wouldMatch("doc<"):
			l.lexDocstring(start)
			continue
		case ch == '2' && l.wouldMatch("2b=2"):
			l.lexSymbol(start)
		case isIdentifierStart(ch):
			l.lexIdentifier(start)
		case unicode.IsDigit(ch):
			l.lexNumber(start)
		default:
			l.lexSymbol(start)
		}
	}

	l.emit(token.EOF, "", l.position())
	return l.tokens, l.diagnostics
}

func (l *lexer) skipLineComment() {
	_ = l.matchString(">>|")
	for !l.atEnd() && l.peek() != '\n' {
		l.advance()
	}
}

func (l *lexer) skipBlockComment(start token.Position) {
	_ = l.matchString("=[")
	for !l.atEnd() {
		if l.matchString("]=") {
			return
		}
		l.advance()
	}
	l.diagnostics = append(l.diagnostics, diagnostics.Diagnostic{
		Severity: diagnostics.SeverityError,
		Message:  "unterminated block comment",
		Pos:      diagnostics.Position{Line: start.Line, Column: start.Column},
	})
}

func (l *lexer) lexDocstring(start token.Position) {
	_ = l.matchString("doc<")
	if !l.atEnd() && l.peek() != '\n' && l.peek() != '\r' {
		l.diagnostics = append(l.diagnostics, diagnostics.Diagnostic{
			Severity: diagnostics.SeverityError,
			Message:  "docstring must be multiline",
			Pos:      diagnostics.Position{Line: start.Line, Column: start.Column},
		})
	}

	var content strings.Builder
	for !l.atEnd() {
		lineStart := l.idx
		for !l.atEnd() && l.peek() != '\n' {
			l.advance()
		}
		line := string(l.src[lineStart:l.idx])
		if strings.TrimSpace(line) == ">" {
			if !l.atEnd() && l.peek() == '\n' {
				l.advance()
			}
			l.emit(token.Docstring, content.String(), start)
			return
		}
		content.WriteString(line)
		if !l.atEnd() && l.peek() == '\n' {
			content.WriteRune('\n')
			l.advance()
		}
	}

	l.diagnostics = append(l.diagnostics, diagnostics.Diagnostic{
		Severity: diagnostics.SeverityError,
		Message:  "unterminated docstring",
		Pos:      diagnostics.Position{Line: start.Line, Column: start.Column},
	})
	l.emit(token.Docstring, content.String(), start)
}

func (l *lexer) lexIdentifier(start token.Position) {
	begin := l.idx
	for !l.atEnd() && isIdentifierPart(l.peek()) {
		l.advance()
	}
	lit := string(l.src[begin:l.idx])
	kind := token.LookupIdentifier(lit)
	l.emit(kind, lit, start)
}

func (l *lexer) lexNumber(start token.Position) {
	begin := l.idx
	for !l.atEnd() && unicode.IsDigit(l.peek()) {
		l.advance()
	}

	kind := token.Integer
	if !l.atEnd() && l.peek() == '.' && l.hasNextDigit() {
		kind = token.Float
		l.advance()
		for !l.atEnd() && unicode.IsDigit(l.peek()) {
			l.advance()
		}
	}

	lit := string(l.src[begin:l.idx])
	l.emit(kind, lit, start)
}

func (l *lexer) lexSymbol(start token.Position) {
	ch := l.peek()

	if ch == '2' && l.matchString("2b=2") {
		l.emit(token.Assign, "2b=2", start)
		return
	}
	if l.matchString(";;") {
		l.emit(token.StmtEnd, ";;", start)
		return
	}
	if l.matchString("-->") {
		l.emit(token.ReturnArrow, "-->", start)
		return
	}
	if l.matchString("->") {
		l.emit(token.ThinArrow, "->", start)
		return
	}
	if l.matchString("<<") {
		l.emit(token.DoubleLAngle, "<<", start)
		return
	}
	if l.matchString(">>") {
		l.emit(token.DoubleRAngle, ">>", start)
		return
	}

	switch ch {
	case '(':
		l.advance()
		l.emit(token.LParen, "(", start)
	case ')':
		l.advance()
		l.emit(token.RParen, ")", start)
	case '{':
		l.advance()
		l.emit(token.LBrace, "{", start)
	case '}':
		l.advance()
		l.emit(token.RBrace, "}", start)
	case '[':
		l.advance()
		l.emit(token.LBracket, "[", start)
	case ']':
		l.advance()
		l.emit(token.RBracket, "]", start)
	case '<':
		l.advance()
		l.emit(token.LAngle, "<", start)
	case '>':
		l.advance()
		l.emit(token.RAngle, ">", start)
	case ',':
		l.advance()
		l.emit(token.Comma, ",", start)
	case '.':
		l.advance()
		l.emit(token.Dot, ".", start)
	case ':':
		l.advance()
		l.emit(token.Colon, ":", start)
	case '+':
		l.advance()
		l.emit(token.Plus, "+", start)
	case '-':
		l.advance()
		l.emit(token.Minus, "-", start)
	case '*':
		l.advance()
		l.emit(token.Star, "*", start)
	case '/':
		l.advance()
		l.emit(token.Slash, "/", start)
	case '%':
		l.advance()
		l.emit(token.Percent, "%", start)
	case '"':
		l.lexString(start)
	default:
		l.advance()
		l.emit(token.Illegal, string(ch), start)
		l.diagnostics = append(l.diagnostics, diagnostics.Diagnostic{
			Severity: diagnostics.SeverityError,
			Message:  "unrecognized character",
			Pos: diagnostics.Position{
				Line:   start.Line,
				Column: start.Column,
			},
		})
	}
}

func (l *lexer) lexString(start token.Position) {
	l.advance() // opening quote
	begin := l.idx

	for !l.atEnd() {
		if l.peek() == '"' {
			lit := string(l.src[begin:l.idx])
			l.advance()
			l.emit(token.String, lit, start)
			return
		}
		if l.peek() == '\n' {
			l.diagnostics = append(l.diagnostics, diagnostics.Diagnostic{
				Severity: diagnostics.SeverityError,
				Message:  "unterminated string literal",
				Pos:      diagnostics.Position{Line: start.Line, Column: start.Column},
			})
			l.emit(token.Illegal, "", start)
			return
		}
		if l.peek() == '\\' && !l.atEnd() {
			l.advance()
			if !l.atEnd() {
				l.advance()
			}
			continue
		}
		l.advance()
	}

	l.diagnostics = append(l.diagnostics, diagnostics.Diagnostic{
		Severity: diagnostics.SeverityError,
		Message:  "unterminated string literal",
		Pos:      diagnostics.Position{Line: start.Line, Column: start.Column},
	})
	l.emit(token.Illegal, "", start)
}

func (l *lexer) skipWhitespace() {
	for !l.atEnd() {
		r := l.peek()
		if r == ' ' || r == '\t' || r == '\r' || r == '\n' {
			l.advance()
			continue
		}
		break
	}
}

func (l *lexer) emit(kind token.Kind, lexeme string, pos token.Position) {
	l.tokens = append(l.tokens, token.Token{
		Kind:   kind,
		Lexeme: lexeme,
		Pos:    pos,
	})
}

func (l *lexer) matchString(s string) bool {
	runes := []rune(s)
	if l.idx+len(runes) > len(l.src) {
		return false
	}
	for i := range runes {
		if l.src[l.idx+i] != runes[i] {
			return false
		}
	}
	for range runes {
		l.advance()
	}
	return true
}

func (l *lexer) wouldMatch(s string) bool {
	runes := []rune(s)
	if l.idx+len(runes) > len(l.src) {
		return false
	}
	for i := range runes {
		if l.src[l.idx+i] != runes[i] {
			return false
		}
	}
	return true
}

func (l *lexer) hasNextDigit() bool {
	next := l.idx + 1
	return next < len(l.src) && unicode.IsDigit(l.src[next])
}

func (l *lexer) atEnd() bool {
	return l.idx >= len(l.src)
}

func (l *lexer) peek() rune {
	if l.atEnd() {
		return 0
	}
	return l.src[l.idx]
}

func (l *lexer) advance() rune {
	if l.atEnd() {
		return 0
	}
	r := l.src[l.idx]
	l.idx++
	if r == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
	return r
}

func (l *lexer) position() token.Position {
	return token.Position{Line: l.line, Column: l.column}
}

func isIdentifierStart(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func isIdentifierPart(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}
