// Package parser - contains Parser which transforms tokens into AST
package parser

import (
	"fmt"

	"github.com/technoboom/compiler/ast"
	"github.com/technoboom/compiler/lexer"
	"github.com/technoboom/compiler/token"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Parser - structure for storing lexer and state of parsing
type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string // errors for debugging

	// map of prefix parse funtions associated with tokens types
	prefixParseFns map[token.Type]prefixParseFn
	// map of infix parse funtions associated with tokens types
	infixParseFns map[token.Type]infixParseFn
}

// New - creates new Parser accordingly to the lexer in the args
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// read two tokens to ensure that curToken and peekToken are
	// both set
	p.nextToken()
	p.nextToken()

	return p
}

// Reads next token from the Lexer
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// Errors - returns all errors collected by the parser
func (p *Parser) Errors() []string {
	return p.errors
}

// peekError - adds an error to the parser errors array
func (p *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf("expected next token to be '%s', got '%s' instead",
		t,
		p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// parseStatement - parses the statement to make a decision what kind of
// statement it can be, after calls the appropriate function to perform
// correct action with the statement and return ast.Statement object
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

// parseLetStatement - parses let statement
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseReturnStatement - parses return statement
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	// iterate while it's not a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// curTokenIs - checks if current token type is a given type
func (p *Parser) curTokenIs(t token.Type) bool {
	return p.curToken.Type == t
}

// peekTokenIs - checks if peek token type is a given type
func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

// expectPeek - peeks the token if types match, otherwise, returns false
func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

// ParseProgram - parses root node and produces AST
// Build the root node of the AST
// After this it reads tokens one by one until he reached token.EOF
// On each iteration it does parsing statement, if it's success - this statement
// adds to the Statements of the program
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	// read until we reaced the end of the file
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		// proceed with next token
		p.nextToken()
	}
	return program
}

// registerInfix - registers function for parsing prefix for the token
func (p *Parser) registerPrefix(tokenType token.Type, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// registerInfix - registers function for parsing infix for the token
func (p *Parser) registerInfix(tokenType token.Type, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
