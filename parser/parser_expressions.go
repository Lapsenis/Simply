package parser

import (
	"Simply/ast"
	"Simply/lexer"
	"fmt"
	"strconv"
)

// Expression priority
const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

type (
	prefixParseFunc func() ast.Node
	infixParseFunc  func(ast.Node) ast.Node
)

var precedences = map[lexer.TokenType]int{
	lexer.EQ:       EQUALS,
	lexer.NOT_EQ:   EQUALS,
	lexer.LT:       LESSGREATER,
	lexer.GT:       LESSGREATER,
	lexer.PLUS:     SUM,
	lexer.MINUS:    SUM,
	lexer.SLASH:    PRODUCT,
	lexer.ASTERISK: PRODUCT,
	lexer.LPAREN:   CALL,
}

func (p *Parser) nextPrecedence() int {
	if p, ok := precedences[p.lookAheadToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {

		return p
	}
	return LOWEST
}

func (p *Parser) registerPrefixParsers() {
	p.prefixParseFuncMap[lexer.IDENTIFIER] = p.parseIdentifier
	p.prefixParseFuncMap[lexer.FUNCTION] = p.parseFunctionLiteral
	p.prefixParseFuncMap[lexer.IF] = p.parseIfExpression

	p.prefixParseFuncMap[lexer.INT] = p.parseIntegerLiteral
	p.prefixParseFuncMap[lexer.TRUE] = p.parseBooleanLiteral
	p.prefixParseFuncMap[lexer.FALSE] = p.parseBooleanLiteral
	p.prefixParseFuncMap[lexer.STRING] = p.parseStringLiteral

	p.prefixParseFuncMap[lexer.MINUS] = p.parsePrefixExpression
	p.prefixParseFuncMap[lexer.BANG] = p.parsePrefixExpression

	p.prefixParseFuncMap[lexer.LPAREN] = p.parseGroupedExpression
}

func (p *Parser) registerInfixParsers() {
	p.infixParseFuncMap[lexer.PLUS] = p.parseInfixExpression
	p.infixParseFuncMap[lexer.MINUS] = p.parseInfixExpression
	p.infixParseFuncMap[lexer.SLASH] = p.parseInfixExpression
	p.infixParseFuncMap[lexer.ASTERISK] = p.parseInfixExpression
	p.infixParseFuncMap[lexer.EQ] = p.parseInfixExpression
	p.infixParseFuncMap[lexer.NOT_EQ] = p.parseInfixExpression
	p.infixParseFuncMap[lexer.LT] = p.parseInfixExpression
	p.infixParseFuncMap[lexer.GT] = p.parseInfixExpression
	p.infixParseFuncMap[lexer.LPAREN] = p.parseCallExpression
}

func (p *Parser) logParseError(msg string, args ...string) {
	p.Errors = append(p.Errors, fmt.Sprintf(msg, args))
}

func (p *Parser) parseExpression(precedence int) ast.Node {
	prefixFunc, ok := p.prefixParseFuncMap[p.currentToken.Type]
	if !ok {
		p.logParseError("Missing prefix parser for %s", string(p.currentToken.Type))
		return nil
	}

	leftExpression := prefixFunc()

	for !p.nextTokenIs(lexer.SEMICOLON) && precedence < p.nextPrecedence() {
		infixFunc := p.infixParseFuncMap[p.lookAheadToken.Type]

		if infixFunc == nil {
			return leftExpression
		}

		p.nextToken()

		leftExpression = infixFunc(leftExpression)
	}

	return leftExpression
}

func (p *Parser) parseIdentifier() ast.Node {
	return &ast.Identifier{Value: p.currentToken.Literal}
}

func (p *Parser) parseFunctionLiteral() ast.Node {
	fl := &ast.FunctionLiteral{}

	if !p.assertToken(lexer.LPAREN) {
		return nil
	}

	fl.Parameters = p.parseFunctionParameters()

	if !p.assertToken(lexer.LBRACE) {
		return nil
	}

	fl.Body = p.parseCodeBlock()

	return fl
}

func (p *Parser) parseIfExpression() ast.Node {
	if !p.assertToken(lexer.LPAREN) {
		return nil
	}
	p.nextToken()

	c := &ast.ConditionalExpression{}

	c.Condition = p.parseExpression(LOWEST)
	if !p.assertToken(lexer.RPAREN) {
		return nil
	}
	if !p.assertToken(lexer.LBRACE) {
		return nil
	}
	c.True = p.parseCodeBlock()

	if p.nextTokenIs(lexer.ELSE) {
		p.nextToken()
		if !p.assertToken(lexer.LBRACE) {
			return nil
		}
		c.False = p.parseCodeBlock()
	}

	return c
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}
	if p.nextTokenIs(lexer.RPAREN) {
		p.nextToken()
		return identifiers
	}
	p.nextToken()
	ident := &ast.Identifier{Value: p.currentToken.Literal}
	identifiers = append(identifiers, ident)
	for p.nextTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Value: p.currentToken.Literal}
		identifiers = append(identifiers, ident)
	}
	if !p.assertToken(lexer.RPAREN) {
		return nil
	}
	return identifiers
}

func (p *Parser) parseCodeBlock() *ast.CodeBlock {
	block := &ast.CodeBlock{}
	block.Statements = []ast.Node{}
	p.nextToken()
	for !p.currentTokenIs(lexer.RBRACE) && !p.currentTokenIs(lexer.EOF) {
		s := p.parseStatement()
		block.Statements = append(block.Statements, s)
		p.nextToken()
	}
	return block
}

func (p *Parser) parseIntegerLiteral() ast.Node {
	il := ast.IntLiteral{}

	v, e := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	if e != nil {
		p.logParseError("Failed to parse intiger: %s", p.currentToken.Literal)
	}

	il.Value = v

	return &il
}

func (p *Parser) parseBooleanLiteral() ast.Node {
	return &ast.BoolLiteral{Value: p.currentTokenIs(lexer.TRUE)}
}

func (p *Parser) parseStringLiteral() ast.Node {
	return &ast.StringLiteral{Value: p.currentToken.Literal}
}

func (p *Parser) parsePrefixExpression() ast.Node {
	e := &ast.PrefixExpression{Prefix: p.currentToken.Literal}

	p.nextToken()

	e.Expression = p.parseExpression(PREFIX)

	return e
}

func (p *Parser) parseGroupedExpression() ast.Node {
	p.nextToken()

	e := p.parseExpression(LOWEST)

	if !p.assertToken(lexer.RPAREN) {
		return nil
	}

	return e
}

func (p *Parser) parseInfixExpression(node ast.Node) ast.Node {
	i := &ast.InfixExpression{Left: node, Operator: p.currentToken.Literal}
	precedence := p.currentPrecedence()
	p.nextToken()
	i.Right = p.parseExpression(precedence)
	return i
}

func (p *Parser) parseCallExpression(node ast.Node) ast.Node {
	exp := &ast.CallExpression{Function: node}
	exp.Arguments = p.parseCallArguments()

	return exp
}

func (p *Parser) parseCallArguments() []ast.Node {
	args := []ast.Node{}
	if p.nextTokenIs(lexer.RPAREN) {
		p.nextToken()
		return args
	}
	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))
	for p.nextTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}
	if !p.assertToken(lexer.RPAREN) {
		return nil
	}
	return args
}
