package parser

import (
	"fmt"
	"os"
	"rexlang/frontend/ast"
	"rexlang/frontend/lexer"
	"rexlang/utils"
	"strings"
)

type Parser struct {
	tokens []lexer.Token
	pos    int
	Lines  []string
}

func Parse(filepath string, debugMode bool) ast.ProgramStmt {

	bytes, err := os.ReadFile(filepath)

	if err != nil {
		panic(err)
	}

	source := string(bytes)

	tokens := lexer.Tokenize(source, debugMode)

	createTokenLookups()
	createTokenTypesLookups()

	parser := &Parser{
		tokens: tokens,
		pos:    0,
		Lines:  strings.Split(source, "\n"),
	}

	var Contents []ast.Stmt

	for parser.hasTokens() {
		Contents = append(Contents, parse_stmt(parser))
	}

	file := ast.File{
		Filename: filepath,
	}

	return ast.ProgramStmt{
		Contents: Contents,
		File:     file,
	}
}

func (p *Parser) currentTokenKind() lexer.TOKEN_KIND {
	return p.currentToken().Kind
}

// Helper functions
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

func (p *Parser) expectError(expectedKind lexer.TOKEN_KIND, err any) lexer.Token {
	token := p.currentToken()
	kind := token.Kind

	if kind != expectedKind {
		if err == nil {
			panic(MakeErrorStr(p, token, fmt.Sprintf("Expected %s but received %s instead\n", expectedKind, kind)))
		} else {
			if errMsg, ok := err.(string); ok {
				panic(MakeErrorStr(p, token, errMsg))
			} else {
				// Handle error if it's not a string
				panic(MakeErrorStr(p, token, "An unexpected error occurred"))
			}
		}
	}

	return p.advance()
}

func (p *Parser) expect(expectedKind lexer.TOKEN_KIND) lexer.Token {
	return p.expectError(expectedKind, nil)
}

func (p *Parser) expectAny(kinds ...lexer.TOKEN_KIND) lexer.Token {
	for _, kind := range kinds {
		if p.currentTokenKind() == kind {
			return p.advance()
		}
	}

	panic(fmt.Sprintf("Expected any of %v but recieved %s instead\n", kinds, p.currentTokenKind()))
}

func MakeErrorStr(p *Parser, token lexer.Token, errMsg string) string {

	// decorate the error with ~~~~~ under the error line

	var errStr string

	line := p.Lines[token.StartPos.Line-1]

	//fmt.Printf("Token start: %v, end: %v\n", token.StartPos, token.EndPos)
	errStr += fmt.Sprintf("Error at %d:%d\n", token.StartPos.Line, token.StartPos.Column)

	paddind := fmt.Sprintf("%d | ", token.StartPos.Line)

	errStr += fmt.Sprintf("%s%s\n", paddind, line)
	errStr += strings.Repeat(" ", (token.StartPos.Column-1) + len(paddind))
	errStr += fmt.Sprint(utils.Colorize(utils.RED, fmt.Sprintf("%s\n", strings.Repeat("~", token.EndPos.Column-token.StartPos.Column))))
	errStr += fmt.Sprint(utils.Colorize(utils.RED, fmt.Sprintf("Error: %s\n", errMsg)))

	return errStr
}
