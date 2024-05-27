package parser

import (
	"fmt"
	"os"

	//"path/filepath"
	"walrus/frontend/ast"
	"walrus/frontend/lexer"
	"walrus/utils"
	"strings"
)

type Parser struct {
	tokens   []lexer.Token
	pos      int
	Lines    *[]string
	FilePath string
}

func Parse(fileSrc string, debugMode bool) ast.ProgramStmt {

	//read file and file data

	bytes, err := os.ReadFile(fileSrc)

	if err != nil {
		panic(err)
	}

	source := string(bytes)

	//filePath := filepath.Base(fileSrc)
	filePath := fileSrc

	tokens, lines := lexer.Tokenize(source, filePath, debugMode)

	createTokenLookups()
	createTokenTypesLookups()

	parser := &Parser{
		tokens:   tokens,
		pos:      0,
		Lines:    lines,
		FilePath: filePath,
	}

	var moduleName string
	var imports []ast.ImportStmt
	var contents []ast.Node

	for parser.hasTokens() {
		stmt := parse_node(parser)

		switch v := stmt.(type) {
		case ast.ModuleStmt:
			moduleName = v.ModuleName
		case ast.ImportStmt:
			imports = append(imports, v)
		default:
			contents = append(contents, stmt)
		}

	}

	end := tokens[len(tokens)-1].EndPos

	return ast.ProgramStmt{
		BaseStmt: ast.BaseStmt{
			Kind: ast.PROGRAM,
			StartPos: lexer.Position{
				Line:   1,
				Column: 1,
				Index: 0,
			},
			EndPos: end,
		},
		ModuleName: moduleName,
		Imports:    imports,
		Contents:   contents,
		FileName:   filePath,
	}
}

func (p *Parser) currentTokenKind() lexer.TOKEN_KIND {
	return p.currentToken().Kind
}

// Helper functions
func (p *Parser) currentToken() lexer.Token {
	return p.tokens[p.pos]
}

func (p *Parser) nextToken() lexer.Token {
	if p.pos+1 < len(p.tokens) {
		return p.tokens[p.pos+1]
	}
	return lexer.Token{}
}

func (p *Parser) previousToken() lexer.Token {
	if p.pos-1 >= 0 {
		return p.tokens[p.pos-1]
	}
	return lexer.Token{}
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
			fmt.Print(p.currentToken(), p.FilePath, token, expectedKind, kind)
			p.MakeError(p.currentToken().StartPos.Line, p.FilePath, token, fmt.Sprintf("Expected %s but received %s instead\n", expectedKind, kind)).Display()
		} else {
			if errMsg, ok := err.(string); ok {
				p.MakeError(p.currentToken().StartPos.Line, p.FilePath, token, errMsg).Display()
			} else {
				// Handle error if it's not a string
				p.MakeError(p.currentToken().StartPos.Line, p.FilePath, token, "An unexpected error occurred").Display()
			}
		}
	}

	return p.advance()
}

func (p *Parser) expect(expectedKind lexer.TOKEN_KIND) lexer.Token {
	return p.expectError(expectedKind, nil)
}

type ErrorMessage struct {
	Message string
	hints   []string
}

func (e *ErrorMessage) AddHint(hint string) *ErrorMessage {
	e.hints = append(e.hints, hint)
	return e
}

func (e *ErrorMessage) Display() {
	fmt.Println(e.Message)
	// hints
	for _, hint := range e.hints {
		fmt.Printf("Hint: %s\n", hint)
	}
	os.Exit(1)
}

func (p *Parser) MakeError(lineNo int, filePath string, token lexer.Token, errMsg string) *ErrorMessage {

	// decorate the error with ~~~~~ under the error line

	var errStr string

	var prvLines []string
	var nxtLines []string
	line := (*p.Lines)[lineNo-1]
	maxWidth := len(fmt.Sprintf("%d", len(*p.Lines)))

	if lineNo-1 > 0 {
		// add the padding to each line
		prvLines = (*p.Lines)[lineNo-2 : lineNo-1]

		for i, l := range prvLines {
			prvLines[i] = utils.Colorize(utils.GREY, fmt.Sprintf("%*d | ", maxWidth, lineNo-1+i)) + lexer.Highlight(l)
		}
	}

	if lineNo+1 < len(*p.Lines) {
		nxtLines = (*p.Lines)[lineNo : lineNo+1]

		for i, l := range nxtLines {
			nxtLines[i] = utils.Colorize(utils.GREY, fmt.Sprintf("%*d | ", maxWidth, lineNo+1+i)) + lexer.Highlight(l)
		}
	}

	errStr += fmt.Sprintf("\n%s:%d:%d\n", filePath, token.StartPos.Line, token.StartPos.Column)

	padding := fmt.Sprintf("%*d | ", maxWidth, token.StartPos.Line)

	errStr += strings.Join(prvLines, "\n") + "\n"
	errStr += utils.Colorize(utils.GREY, padding) + lexer.Highlight(line[0:token.StartPos.Column-1]) + utils.Colorize(utils.RED, line[token.StartPos.Column-1:token.EndPos.Column-1]) + lexer.Highlight(line[token.EndPos.Column-1:]) + "\n"
	errStr += strings.Repeat(" ", (token.StartPos.Column-1)+len(padding))
	errStr += fmt.Sprint(utils.Colorize(utils.BOLD_RED, fmt.Sprintf("%s%s\n", "^", strings.Repeat("~", token.EndPos.Column-token.StartPos.Column-1))))
	errStr += strings.Join(nxtLines, "\n") + "\n"
	errStr += fmt.Sprint(utils.Colorize(utils.RED, fmt.Sprintf("Error: %s\n", errMsg)))

	return &ErrorMessage{
		Message: errStr,
	}
}
