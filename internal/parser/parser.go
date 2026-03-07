package parser

import (
	"fmt"
	"strings"

	"github.com/asciifaceman/jave/internal/ast"
	"github.com/asciifaceman/jave/internal/diagnostics"
	"github.com/asciifaceman/jave/internal/token"
)

// Parse builds an AST from lexer tokens.
func Parse(tokens []token.Token) (*ast.Program, []diagnostics.Diagnostic) {
	p := &parser{tokens: tokens}
	return p.parseProgram(), p.diags
}

type parser struct {
	tokens []token.Token
	idx    int
	diags  []diagnostics.Diagnostic
}

func (p *parser) parseProgram() *ast.Program {
	program := &ast.Program{}
	for !p.at(token.EOF) {
		switch p.current().Kind {
		case token.Install:
			imp := p.parseImport()
			if imp != nil {
				program.Imports = append(program.Imports, *imp)
			}
		case token.Outy, token.Inny:
			seq := p.parseSequence()
			if seq != nil {
				program.Sequences = append(program.Sequences, *seq)
			}
		default:
			p.errorHere("expected top-level 'install' or sequence declaration")
			p.synchronizeTopLevel()
		}
	}
	return program
}

func (p *parser) parseImport() *ast.ImportDecl {
	installTok := p.advance() // install
	nameTok, ok := p.expect(token.Identifier, "expected import name after 'install'")
	if !ok {
		return nil
	}
	if _, ok := p.expect(token.From, "expected 'from' in import declaration"); !ok {
		return nil
	}

	parts := make([]string, 0)
	for !p.at(token.StmtEnd) && !p.at(token.EOF) {
		parts = append(parts, p.advance().Lexeme)
	}
	if _, ok := p.expect(token.StmtEnd, "expected ';;' after import declaration"); !ok {
		return nil
	}

	return &ast.ImportDecl{Pos: installTok.Pos, Name: nameTok.Lexeme, From: strings.Join(parts, "")}
}

func (p *parser) parseSequence() *ast.SequenceDecl {
	visTok := p.current()
	if !(p.match(token.Outy) || p.match(token.Inny)) {
		p.errorHere("expected visibility keyword 'outy' or 'inny'")
		return nil
	}

	if !(p.match(token.Seq) || p.match(token.Sequence)) {
		p.errorHere("expected sequence declaration after visibility")
		return nil
	}

	nameTok, ok := p.expect(token.Identifier, "expected sequence name")
	if !ok {
		return nil
	}

	if _, ok := p.expect(token.LAngle, "expected '<' to start sequence type arguments"); !ok {
		return nil
	}
	if !p.consumeBalancedAngles() {
		p.errorHere("unterminated sequence type arguments")
		return nil
	}

	if _, ok := p.expect(token.ReturnArrow, "expected '-->' after sequence signature"); !ok {
		return nil
	}
	if _, ok := p.expect(token.DoubleLAngle, "expected '<<' before return type"); !ok {
		return nil
	}
	retTok, ok := p.expectOneOf([]token.Kind{token.Identifier, token.TypeExact, token.TypeVag, token.TypeTruther, token.TypeStrang, token.TypeNada, token.TypeNaw}, "expected return type name")
	if !ok {
		return nil
	}
	if _, ok := p.expect(token.DoubleRAngle, "expected '>>' after return type"); !ok {
		return nil
	}

	body, ok := p.parseBlockContents()
	if !ok {
		return nil
	}

	vis := visTok.Lexeme
	if vis == "" {
		vis = string(visTok.Kind)
	}

	return &ast.SequenceDecl{
		Pos:        visTok.Pos,
		Visibility: vis,
		Name:       nameTok.Lexeme,
		ReturnType: retTok.Lexeme,
		Body:       body,
	}
}

func (p *parser) parseStatement() ast.Stmt {
	switch p.current().Kind {
	case token.Give:
		return p.parseGive()
	case token.Allow:
		return p.parseVarDecl()
	case token.Maybe:
		return p.parseIfStmt()
	case token.Given:
		return p.parseGivenStmt()
	case token.Identifier:
		if p.peek(1).Kind == token.Assign {
			return p.parseAssignment()
		}
		fallthrough
	default:
		expr := p.parseExpr()
		if expr == nil {
			return nil
		}
		if _, ok := p.expect(token.StmtEnd, "expected ';;' after expression statement"); !ok {
			return nil
		}
		return ast.ExprStmt{Pos: exprPos(expr), Expr: expr}
	}
}

func (p *parser) parseGive() ast.Stmt {
	giveTok := p.advance() // give
	if p.match(token.Up) {
		if _, ok := p.expect(token.StmtEnd, "expected ';;' after give up"); !ok {
			return nil
		}
		return ast.GiveStmt{Pos: giveTok.Pos}
	}

	val := p.parseExpr()
	if val == nil {
		return nil
	}
	if _, ok := p.expect(token.Up, "expected 'up' after return value"); !ok {
		return nil
	}
	if _, ok := p.expect(token.StmtEnd, "expected ';;' after give statement"); !ok {
		return nil
	}
	return ast.GiveStmt{Pos: giveTok.Pos, Value: val}
}

func (p *parser) parseVarDecl() ast.Stmt {
	allowTok := p.advance() // allow
	typeName, ok := p.parseTypeName()
	if !ok {
		return nil
	}
	nameTok, ok := p.expect(token.Identifier, "expected variable name")
	if !ok {
		return nil
	}
	if _, ok := p.expect(token.Assign, "expected '2b=2' in declaration"); !ok {
		return nil
	}
	value := p.parseExpr()
	if value == nil {
		return nil
	}
	if _, ok := p.expect(token.StmtEnd, "expected ';;' after declaration"); !ok {
		return nil
	}
	return ast.VarDeclStmt{Pos: allowTok.Pos, TypeName: typeName, Name: nameTok.Lexeme, Value: value}
}

func (p *parser) parseExpr() ast.Expr {
	return p.parseLogicalOr()
}

func (p *parser) parseLogicalOr() ast.Expr {
	left := p.parseLogicalAnd()
	for p.matchWord("orelse") {
		right := p.parseLogicalAnd()
		if right == nil {
			return left
		}
		left = ast.BinaryExpr{Left: left, Op: "orelse", Right: right}
	}
	return left
}

func (p *parser) parseLogicalAnd() ast.Expr {
	left := p.parseComparison()
	for p.matchWord("plusalso") {
		right := p.parseComparison()
		if right == nil {
			return left
		}
		left = ast.BinaryExpr{Left: left, Op: "plusalso", Right: right}
	}
	return left
}

func (p *parser) parseComparison() ast.Expr {
	left := p.parseAdditive()
	for p.isComparisonWord(p.current()) {
		op := p.advance().Lexeme
		right := p.parseAdditive()
		if right == nil {
			return left
		}
		left = ast.BinaryExpr{Left: left, Op: op, Right: right}
	}
	return left
}

func (p *parser) parseAdditive() ast.Expr {
	left := p.parseMultiplicative()
	for p.at(token.Plus) || p.at(token.Minus) {
		op := p.advance().Lexeme
		right := p.parseMultiplicative()
		if right == nil {
			return left
		}
		left = ast.BinaryExpr{Left: left, Op: op, Right: right}
	}
	return left
}

func (p *parser) parseMultiplicative() ast.Expr {
	left := p.parsePostfix()
	for p.at(token.Star) || p.at(token.Slash) {
		op := p.advance().Lexeme
		right := p.parsePostfix()
		if right == nil {
			return left
		}
		left = ast.BinaryExpr{Left: left, Op: op, Right: right}
	}
	return left
}

func (p *parser) parsePostfix() ast.Expr {
	expr := p.parsePrimary()
	if expr == nil {
		return nil
	}

	for {
		switch {
		case p.match(token.Dot):
			nameTok, ok := p.expect(token.Identifier, "expected member name after '.'")
			if !ok {
				return expr
			}
			expr = ast.MemberExpr{Target: expr, Name: nameTok.Lexeme}
		case p.match(token.LBracket):
			idxExpr := p.parseExpr()
			if idxExpr == nil {
				return expr
			}
			if _, ok := p.expect(token.RBracket, "expected ']' after index expression"); !ok {
				return expr
			}
			expr = ast.IndexExpr{Target: expr, Index: idxExpr}
		case p.at(token.LParen):
			args, ok := p.parseCallArgs(token.LParen, token.RParen)
			if !ok {
				return expr
			}
			expr = ast.CallExpr{Callee: expr, Args: args, Style: "paren"}
		case p.at(token.LAngle):
			args, ok := p.parseCallArgs(token.LAngle, token.RAngle)
			if !ok {
				return expr
			}
			expr = ast.CallExpr{Callee: expr, Args: args, Style: "angle"}
		default:
			return expr
		}
	}
}

func (p *parser) parseCallArgs(open token.Kind, close token.Kind) ([]ast.Expr, bool) {
	if _, ok := p.expect(open, fmt.Sprintf("expected %q", open)); !ok {
		return nil, false
	}
	args := make([]ast.Expr, 0)
	if p.match(close) {
		return args, true
	}
	for {
		arg := p.parseExpr()
		if arg == nil {
			return nil, false
		}
		args = append(args, arg)
		if p.match(close) {
			return args, true
		}
		if _, ok := p.expect(token.Comma, "expected ',' between call arguments"); !ok {
			return nil, false
		}
	}
}

func (p *parser) parsePrimary() ast.Expr {
	tok := p.current()
	switch tok.Kind {
	case token.Identifier:
		p.advance()
		return ast.IdentifierExpr{Pos: tok.Pos, Name: tok.Lexeme}
	case token.Integer, token.Float:
		p.advance()
		return ast.NumberExpr{Value: tok.Lexeme}
	case token.String:
		p.advance()
		return ast.StringExpr{Value: tok.Lexeme}
	case token.Yee:
		p.advance()
		return ast.BoolExpr{Value: true}
	case token.Nee:
		p.advance()
		return ast.BoolExpr{Value: false}
	case token.LBracket:
		items, ok := p.parseExprList(token.LBracket, token.RBracket)
		if !ok {
			return nil
		}
		return ast.CollectionLiteralExpr{Form: "table", Items: items}
	case token.LAngle:
		items, ok := p.parseExprList(token.LAngle, token.RAngle)
		if !ok {
			return nil
		}
		return ast.CollectionLiteralExpr{Form: "enumeration", Items: items}
	case token.LBrace:
		pairs, ok := p.parseMapLiteral()
		if !ok {
			return nil
		}
		return ast.CollectionLiteralExpr{Form: "lexis", Pairs: pairs}
	case token.LParen:
		p.advance()
		expr := p.parseExpr()
		if _, ok := p.expect(token.RParen, "expected ')' after expression"); !ok {
			return expr
		}
		return expr
	default:
		p.errorHere("expected expression")
		return nil
	}
}

func (p *parser) parseAssignment() ast.Stmt {
	nameTok := p.advance()
	if _, ok := p.expect(token.Assign, "expected '2b=2' in assignment"); !ok {
		return nil
	}
	value := p.parseExpr()
	if value == nil {
		return nil
	}
	if _, ok := p.expect(token.StmtEnd, "expected ';;' after assignment"); !ok {
		return nil
	}
	return ast.AssignmentStmt{Pos: nameTok.Pos, Name: nameTok.Lexeme, Value: value}
}

func (p *parser) parseIfStmt() ast.Stmt {
	maybeTok := p.advance() // maybe
	cond := p.parseConditionWrapper()
	if cond == nil {
		return nil
	}
	if _, ok := p.expect(token.ThinArrow, "expected '->' after maybe condition"); !ok {
		return nil
	}
	body, ok := p.parseBlockContents()
	if !ok {
		return nil
	}

	stmt := ast.IfStmt{Pos: maybeTok.Pos, Branches: []ast.IfBranch{{Condition: cond, Body: body}}}
	for p.match(token.Furthermore) {
		branchCond := p.parseConditionWrapper()
		if branchCond == nil {
			return nil
		}
		if _, ok := p.expect(token.ThinArrow, "expected '->' after furthermore condition"); !ok {
			return nil
		}
		branchBody, ok := p.parseBlockContents()
		if !ok {
			return nil
		}
		stmt.Branches = append(stmt.Branches, ast.IfBranch{Condition: branchCond, Body: branchBody})
	}

	if p.match(token.Otherwise) {
		if _, ok := p.expect(token.ThinArrow, "expected '->' after otherwise"); !ok {
			return nil
		}
		elseBody, ok := p.parseBlockContents()
		if !ok {
			return nil
		}
		stmt.ElseBody = elseBody
	}

	return stmt
}

func (p *parser) parseGivenStmt() ast.Stmt {
	givenTok := p.advance() // given
	headerToks, ok := p.parseGivenHeaderTokens()
	if !ok {
		return nil
	}

	mode := detectGivenMode(headerToks)
	var condExpr ast.Expr
	var initStmt ast.Stmt
	var stepStmt ast.Stmt
	var iterVar string
	var iterExpr ast.Expr
	if mode == "while-ish" {
		condExpr = p.parseExprFromTokens(headerToks)
		if condExpr == nil {
			p.errorHere("unable to parse while-ish given condition")
			return nil
		}
		if _, ok := p.expect(token.Again, "expected 'again' in while-ish given loop"); !ok {
			return nil
		}
	} else if mode == "for-ish" {
		initStmt, condExpr, stepStmt, ok = p.parseForHeader(headerToks)
		if !ok {
			return nil
		}
	} else if mode == "within" {
		iterVar, iterExpr, ok = p.parseWithinHeader(headerToks)
		if !ok {
			return nil
		}
	}

	if _, ok := p.expect(token.ThinArrow, "expected '->' after given header"); !ok {
		return nil
	}
	body, ok := p.parseBlockContents()
	if !ok {
		return nil
	}

	return ast.GivenStmt{
		Pos:    givenTok.Pos,
		Mode:   mode,
		Header: tokensText(headerToks),
		Cond:   condExpr,
		Init:   initStmt,
		Step:   stepStmt,
		Var:    iterVar,
		In:     iterExpr,
		Body:   body,
	}
}

func (p *parser) parseConditionWrapper() ast.Expr {
	if _, ok := p.expect(token.LParen, "expected '(' before condition wrapper"); !ok {
		return nil
	}
	if _, ok := p.expect(token.LAngle, "expected '<' before condition expression"); !ok {
		return nil
	}
	cond := p.parseExpr()
	if cond == nil {
		return nil
	}
	if _, ok := p.expect(token.RAngle, "expected '>' after condition expression"); !ok {
		return nil
	}
	if _, ok := p.expect(token.RParen, "expected ')' after condition wrapper"); !ok {
		return nil
	}
	return cond
}

func (p *parser) parseGivenHeaderTokens() ([]token.Token, bool) {
	if _, ok := p.expect(token.LParen, "expected '(' after 'given'"); !ok {
		return nil, false
	}
	if _, ok := p.expect(token.LAngle, "expected '<' after '(' in given header"); !ok {
		return nil, false
	}

	header := make([]token.Token, 0)
	depth := 1
	for !p.at(token.EOF) {
		tok := p.advance()
		if tok.Kind == token.LAngle {
			depth++
		}
		if tok.Kind == token.RAngle {
			depth--
			if depth == 0 {
				break
			}
		}
		header = append(header, tok)
	}
	if depth != 0 {
		p.errorHere("unterminated given header")
		return nil, false
	}
	if _, ok := p.expect(token.RParen, "expected ')' after given header"); !ok {
		return nil, false
	}
	return header, true
}

func (p *parser) parseBlockContents() ([]ast.Stmt, bool) {
	if _, ok := p.expect(token.LBrace, "expected '{' to start block"); !ok {
		return nil, false
	}
	body := make([]ast.Stmt, 0)
	for !p.at(token.RBrace) && !p.at(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			body = append(body, stmt)
			continue
		}
		p.synchronizeStatement()
	}
	if _, ok := p.expect(token.RBrace, "expected '}' to close block"); !ok {
		return nil, false
	}
	return body, true
}

func (p *parser) parseTypeName() (string, bool) {
	base, ok := p.expectOneOf([]token.Kind{token.Identifier, token.TypeExact, token.TypeVag, token.TypeTruther, token.TypeStrang, token.TypeNada, token.TypeNaw, token.Given}, "expected type after 'allow'")
	if !ok {
		return "", false
	}
	name := base.Lexeme
	if !p.at(token.LAngle) {
		return name, true
	}
	suffix, ok := p.consumeTypeSuffix()
	if !ok {
		return "", false
	}
	return name + suffix, true
}

func (p *parser) consumeTypeSuffix() (string, bool) {
	parts := make([]string, 0)
	depth := 0
	for !p.at(token.EOF) {
		tok := p.advance()
		parts = append(parts, tok.Lexeme)
		if tok.Kind == token.LAngle {
			depth++
		}
		if tok.Kind == token.RAngle {
			depth--
			if depth == 0 {
				return strings.Join(parts, ""), true
			}
		}
		if tok.Kind == token.DoubleRAngle {
			depth -= 2
			if depth <= 0 {
				return strings.Join(parts, ""), true
			}
		}
	}
	p.errorHere("unterminated generic type")
	return "", false
}

func (p *parser) consumeBalancedAngles() bool {
	depth := 1
	for !p.at(token.EOF) {
		tok := p.advance()
		if tok.Kind == token.LAngle {
			depth++
		}
		if tok.Kind == token.RAngle {
			depth--
			if depth == 0 {
				return true
			}
		}
		if tok.Kind == token.DoubleRAngle {
			depth -= 2
			if depth <= 0 {
				return true
			}
		}
	}
	return false
}

func (p *parser) parseExprList(open token.Kind, close token.Kind) ([]ast.Expr, bool) {
	if _, ok := p.expect(open, fmt.Sprintf("expected %q", open)); !ok {
		return nil, false
	}
	items := make([]ast.Expr, 0)
	if p.match(close) {
		return items, true
	}
	for {
		item := p.parseExpr()
		if item == nil {
			return nil, false
		}
		items = append(items, item)
		if p.match(close) {
			return items, true
		}
		if _, ok := p.expect(token.Comma, "expected ',' between literal items"); !ok {
			return nil, false
		}
	}
}

func (p *parser) parseMapLiteral() ([]ast.KeyValueExpr, bool) {
	if _, ok := p.expect(token.LBrace, "expected '{'"); !ok {
		return nil, false
	}
	pairs := make([]ast.KeyValueExpr, 0)
	if p.match(token.RBrace) {
		return pairs, true
	}
	for {
		key := p.parseExpr()
		if key == nil {
			return nil, false
		}
		if _, ok := p.expect(token.Colon, "expected ':' between key and value"); !ok {
			return nil, false
		}
		val := p.parseExpr()
		if val == nil {
			return nil, false
		}
		pairs = append(pairs, ast.KeyValueExpr{Key: key, Value: val})
		if p.match(token.RBrace) {
			return pairs, true
		}
		if _, ok := p.expect(token.Comma, "expected ',' between map entries"); !ok {
			return nil, false
		}
	}
}

func (p *parser) synchronizeTopLevel() {
	for !p.at(token.EOF) {
		if p.at(token.Outy) || p.at(token.Inny) || p.at(token.Install) {
			return
		}
		p.advance()
	}
}

func (p *parser) synchronizeStatement() {
	for !p.at(token.EOF) && !p.at(token.RBrace) {
		if p.match(token.StmtEnd) || p.at(token.Maybe) || p.at(token.Given) {
			return
		}
		p.advance()
	}
}

func (p *parser) expect(kind token.Kind, msg string) (token.Token, bool) {
	if p.at(kind) {
		return p.advance(), true
	}
	p.errorHere(msg)
	return token.Token{}, false
}

func (p *parser) expectOneOf(kinds []token.Kind, msg string) (token.Token, bool) {
	for _, k := range kinds {
		if p.at(k) {
			return p.advance(), true
		}
	}
	p.errorHere(msg)
	return token.Token{}, false
}

func (p *parser) match(kind token.Kind) bool {
	if p.at(kind) {
		p.advance()
		return true
	}
	return false
}

func (p *parser) at(kind token.Kind) bool {
	return p.current().Kind == kind
}

func (p *parser) current() token.Token {
	if p.idx >= len(p.tokens) {
		return token.Token{Kind: token.EOF}
	}
	return p.tokens[p.idx]
}

func (p *parser) peek(offset int) token.Token {
	idx := p.idx + offset
	if idx < 0 || idx >= len(p.tokens) {
		return token.Token{Kind: token.EOF}
	}
	return p.tokens[idx]
}

func (p *parser) advance() token.Token {
	tok := p.current()
	if p.idx < len(p.tokens) {
		p.idx++
	}
	return tok
}

func (p *parser) errorHere(message string) {
	tok := p.current()
	p.diags = append(p.diags, diagnostics.Diagnostic{
		Severity: diagnostics.SeverityError,
		Message:  message,
		Pos: diagnostics.Position{
			Line:   tok.Pos.Line,
			Column: tok.Pos.Column,
		},
	})
}

func (p *parser) matchWord(word string) bool {
	tok := p.current()
	if tok.Kind == token.Identifier && tok.Lexeme == word {
		p.advance()
		return true
	}
	return false
}

func (p *parser) isComparisonWord(tok token.Token) bool {
	if tok.Kind != token.Identifier {
		return false
	}
	switch tok.Lexeme {
	case "samewise", "notsamewise", "bigly", "lessly", "biglysame", "lesslysame":
		return true
	default:
		return false
	}
}

func detectGivenMode(header []token.Token) string {
	for _, tok := range header {
		if tok.Kind == token.Within {
			return "within"
		}
		if tok.Kind == token.StmtEnd {
			return "for-ish"
		}
	}
	return "while-ish"
}

func tokensText(tokens []token.Token) string {
	parts := make([]string, 0, len(tokens))
	for _, tok := range tokens {
		parts = append(parts, tok.Lexeme)
	}
	return strings.Join(parts, " ")
}

func (p *parser) parseExprFromTokens(tokensIn []token.Token) ast.Expr {
	toks := make([]token.Token, 0, len(tokensIn)+1)
	toks = append(toks, tokensIn...)
	toks = append(toks, token.Token{Kind: token.EOF})
	sub := &parser{tokens: toks}
	expr := sub.parseExpr()
	if len(sub.diags) > 0 || !sub.at(token.EOF) {
		return nil
	}
	return expr
}

func (p *parser) parseForHeader(header []token.Token) (ast.Stmt, ast.Expr, ast.Stmt, bool) {
	segments := splitByStmtEnd(header)
	if len(segments) < 3 {
		p.errorHere("for-ish given header must contain init, condition, and step")
		return nil, nil, nil, false
	}
	init := p.parseStmtFromHeaderTokens(segments[0])
	if init == nil {
		p.errorHere("unable to parse for-ish init statement")
		return nil, nil, nil, false
	}
	cond := p.parseExprFromTokens(segments[1])
	if cond == nil {
		p.errorHere("unable to parse for-ish condition expression")
		return nil, nil, nil, false
	}
	step := p.parseStmtFromHeaderTokens(segments[2])
	if step == nil {
		p.errorHere("unable to parse for-ish step statement")
		return nil, nil, nil, false
	}
	return init, cond, step, true
}

func (p *parser) parseWithinHeader(header []token.Token) (string, ast.Expr, bool) {
	withinIdx := -1
	for i, tok := range header {
		if tok.Kind == token.Within {
			withinIdx = i
			break
		}
	}
	if withinIdx <= 0 || withinIdx >= len(header)-1 {
		p.errorHere("within header must be '<Name within Expr>'")
		return "", nil, false
	}
	varTok := header[0]
	if varTok.Kind != token.Identifier {
		p.errorHere("within loop variable must be an identifier")
		return "", nil, false
	}
	expr := p.parseExprFromTokens(header[withinIdx+1:])
	if expr == nil {
		p.errorHere("unable to parse within iterable expression")
		return "", nil, false
	}
	return varTok.Lexeme, expr, true
}

func (p *parser) parseStmtFromHeaderTokens(tokensIn []token.Token) ast.Stmt {
	toks := make([]token.Token, 0, len(tokensIn)+2)
	toks = append(toks, tokensIn...)
	toks = append(toks, token.Token{Kind: token.StmtEnd, Lexeme: ";;"})
	toks = append(toks, token.Token{Kind: token.EOF})
	sub := &parser{tokens: toks}
	stmt := sub.parseStatement()
	if stmt == nil || len(sub.diags) > 0 {
		return nil
	}
	return stmt
}

func splitByStmtEnd(tokensIn []token.Token) [][]token.Token {
	out := make([][]token.Token, 0)
	current := make([]token.Token, 0)
	for _, tok := range tokensIn {
		if tok.Kind == token.StmtEnd {
			if len(current) > 0 {
				out = append(out, current)
				current = make([]token.Token, 0)
			}
			continue
		}
		current = append(current, tok)
	}
	if len(current) > 0 {
		out = append(out, current)
	}
	return out
}

func exprPos(expr ast.Expr) token.Position {
	switch e := expr.(type) {
	case ast.IdentifierExpr:
		return e.Pos
	default:
		return token.Position{Line: 1, Column: 1}
	}
}
