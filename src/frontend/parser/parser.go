package parser

import (
	"fmt"
	"os"
	"rexlang/frontend/ast"
	"rexlang/frontend/lexer"
	"strings"
)

type Parser struct {
	tokens 	[]lexer.Token
	pos 	int
	Lines	[]string
}

func Parse(source *string, debugMode bool) ast.BlockStmt {

	Body := make([]ast.Stmt, 0)

	parser := createParser(source, debugMode)

	for parser.hasTokens() {
		Body = append(Body, parse_stmt(parser))
	}

	return ast.BlockStmt{
		Body: Body,
	}
}

func createParser(src *string, debugMode bool) *Parser {
	
	createTokenLookups()
	createTokenTypesLookups()
	
	return &Parser{
		tokens: lexer.Tokenize(src, debugMode),
		pos: 0,
		Lines: strings.Split(*src, "\n"),
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
			err = fmt.Sprintf("At line %d: Expected %s but recieved %s instead\n", token.StartPos.Line, lexer.TokenKindString(expectedKind), lexer.TokenKindString(kind))
		}

		PrintError(p, token, fmt.Sprintf("Expected %s but recieved %s instead\n", lexer.TokenKindString(expectedKind), lexer.TokenKindString(kind)))
	}

	return p.advance()
}

func (p *Parser) expect(expectedKind lexer.TokenKind) lexer.Token {
	return p.expectError(expectedKind, nil)
}

func (p *Parser) expectAny(kinds ...lexer.TokenKind) lexer.Token {
	for _, kind := range kinds {
		if p.currentTokenKind() == kind {
			return p.advance()
		}
	}

	panic(fmt.Sprintf("Expected any of %v but recieved %s instead\n", kinds, lexer.TokenKindString(p.currentTokenKind())))
}


func PrintError(p *Parser, token lexer.Token, errMsg string) {

	// decorate the error with ~~~~~ under the error line

	line := p.Lines[token.StartPos.Line - 1]

	fmt.Printf("\nLine: %d, Column: %d\n", token.StartPos.Line, token.StartPos.Column)
	fmt.Printf("%s\n", line)
	fmt.Printf("%s", strings.Repeat("~", token.StartPos.Column))
	fmt.Printf("%s\n", strings.Repeat("^", token.EndPos.Column - token.StartPos.Column))
	fmt.Printf("\n")
	fmt.Printf("Error: %s\n", errMsg)

	//exit
	os.Exit(1)
}