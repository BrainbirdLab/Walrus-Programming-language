package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"rexlang/frontend/ast"
	"rexlang/frontend/lexer"
	"rexlang/utils"
	"strings"
)

type Parser struct {
	tokens 		[]lexer.Token
	pos    		int
	Lines  		[]string
	FilePath 	string
}

func Parse(fileSrc string, debugMode bool) ast.ProgramStmt {

	//read file and file data

	bytes, err := os.ReadFile(fileSrc)

	if err!= nil {
		panic(err)
	}

	source := string(bytes)


	filePath := filepath.Base(fileSrc)


	tokens := lexer.Tokenize(source, debugMode)

	createTokenLookups()
	createTokenTypesLookups()

	parser := &Parser{
		tokens: tokens,
		pos:    0,
		Lines:  strings.Split(source, "\n"),
		FilePath: filePath,
	}

	var Contents []ast.Stmt

	for parser.hasTokens() {
		Contents = append(Contents, parse_stmt(parser))
	}

	end := tokens[len(tokens) - 1].EndPos

	return ast.ProgramStmt{
		Contents: Contents,
		FileName: filePath,
		StartPos: lexer.Position{
			Line:   1,
			Column: 1,
			Index:  0,
		},
		EndPos: end,
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

func MakeErrorStr(p *Parser, token lexer.Token, errMsg string) string {

	// decorate the error with ~~~~~ under the error line

	var errStr string

	line := p.Lines[token.StartPos.Line-1]

	//fmt.Printf("Token start: %v, end: %v\n", token.StartPos, token.EndPos)
	errStr += fmt.Sprintf("\n%s:%d:%d\n", p.FilePath, token.StartPos.Line, token.StartPos.Column)

	padding := fmt.Sprintf("%d | ", token.StartPos.Line)

	errStr += utils.Colorize(utils.GREY, padding) + Highlight(line[0:token.StartPos.Column - 1]) + utils.Colorize(utils.RED, line[token.StartPos.Column - 1:token.EndPos.Column - 1]) + Highlight(line[token.EndPos.Column - 1:]) + "\n"
	errStr += strings.Repeat(" ", (token.StartPos.Column - 1) + len(padding))
	errStr += fmt.Sprint(utils.Colorize(utils.BOLD_RED, fmt.Sprintf("%s%s\n", "^", strings.Repeat("~", token.EndPos.Column-token.StartPos.Column - 1))))
	errStr += fmt.Sprint(utils.Colorize(utils.RED, fmt.Sprintf("Error: %s\n", errMsg)))

	return errStr
}

func Highlight(line string) string {
	words := strings.Split(line, " ")

	for i, word := range words {
		if lexer.IsKeyword(lexer.TOKEN_KIND(word)) {
			words[i] = utils.Colorize(utils.PURPLE, word)
		} else if lexer.IsNumber(lexer.TOKEN_KIND(word)) {
			words[i] = utils.Colorize(utils.ORANGE, word)
		}
	}

	return strings.Join(words, " ")
}