package parser

import (
	"Simply/ast"
	"Simply/lexer"
	"fmt"
)

type Parser struct {
	t              *lexer.Tokenizer
	currentToken   lexer.Token
	lookAheadToken lexer.Token
	Errors         []string

	prefixParseFuncMap map[lexer.TokenType]prefixParseFunc
	infixParseFuncMap  map[lexer.TokenType]infixParseFunc
}

func NewParser(t *lexer.Tokenizer) *Parser {
	p := &Parser{t: t, Errors: []string{}}

	p.nextToken()
	p.nextToken()

	p.prefixParseFuncMap = make(map[lexer.TokenType]prefixParseFunc)
	p.infixParseFuncMap = make(map[lexer.TokenType]infixParseFunc)

	p.registerPrefixParsers()
	p.registerInfixParsers()

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{Statements: []ast.Node{}}

	for p.currentToken.Type != lexer.EOF {

		s := p.parseStatement()
		program.Statements = append(program.Statements, s)

		p.nextToken()
	}

	return program
}

func (p *Parser) nextToken() {
	p.currentToken = p.lookAheadToken
	p.lookAheadToken = p.t.NextToken()
}

func (p *Parser) currentTokenIs(t lexer.TokenType) bool {
	return p.currentToken.Type == t
}

func (p *Parser) nextTokenIs(t lexer.TokenType) bool {
	return p.lookAheadToken.Type == t
}

func (p *Parser) assertToken(t lexer.TokenType) bool {
	if p.nextTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.logInvalidToken(t)
		return false
	}
}

func (p *Parser) logInvalidToken(t lexer.TokenType) {
	p.Errors = append(p.Errors,
		fmt.Sprintf("Invalid token: expected %s, got %s",
			t,
			p.lookAheadToken.Type,
		),
	)
}

func (p *Parser) parseStatement() ast.Node {
	switch p.currentToken.Type {
	case lexer.LET:
		return p.parseDeclarativeStatement()
	case lexer.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseDeclarativeStatement() *ast.DeclarativeStatement {
	s := &ast.DeclarativeStatement{}

	if !p.assertToken(lexer.IDENTIFIER) {
		return nil
	}

	s.Name = ast.Identifier{Value: p.currentToken.Literal}

	if !p.assertToken(lexer.ASSIGN) {
		return nil
	}

	p.nextToken()

	s.Value = p.parseExpression(LOWEST)

	if p.nextTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return s
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	s := &ast.ReturnStatement{}
	p.nextToken()

	s.Value = p.parseExpression((LOWEST))

	for !p.currentTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}
	return s
}

func (p *Parser) parseExpressionStatement() ast.Node {
	e := &ast.ExpressionStatement{}

	e.Expression = p.parseExpression(LOWEST)

	if p.nextTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return e
}
