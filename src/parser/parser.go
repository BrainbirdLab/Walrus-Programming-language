package parser

import (
	"fmt"
	"rexlang/ast"
	"rexlang/lexer"
)

type Parser struct {
	//!ToDo: error
	tokens 	[]lexer.Token
	pos 	int
}

func Parse(tokens []lexer.Token) ast.BlockStmt {

	Body := make([]ast.Stmt, 0)

	parser := createParser(tokens)

	for parser.hasTokens() {
		Body = append(Body, parse_stmt(parser))
	}

	return ast.BlockStmt{
		Body: Body,
	}
}

func createParser(tokens []lexer.Token) *Parser {
	
	createTokenLookups()
	
	return &Parser{
		tokens: tokens,
		pos: 0,
	}
}

func (p *Parser) currentTokenKind() lexer.TokenKind {
	return p.currentToken().Kind
}

//Helper functions
func (p *Parser) currentToken() lexer.Token {
	return p.tokens[p.pos]
}

func (p *Parser) advance() lexer.Token {
	token := p.currentToken()
	p.pos++
	return token
}

func (p *Parser) hasTokens() bool {
	return p.pos < len(p.tokens) && p.currentTokenKind() != lexer.EOF
}

func (p *Parser) expectError(expectedKind lexer.TokenKind, err any) lexer.Token {
	token := p.currentToken()
	kind := token.Kind

	if kind != expectedKind {
		if err == nil {
			err = fmt.Sprintf("Expected %s but recieved %s instead\n", lexer.TokenKindString(expectedKind), lexer.TokenKindString(kind))
		}

		panic(err)
	}

	return p.advance()
}

func (p *Parser) expect(expectedKind lexer.TokenKind) lexer.Token {
	return p.expectError(expectedKind, nil)
}