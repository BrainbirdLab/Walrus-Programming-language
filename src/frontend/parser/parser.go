package parser

import (
	"fmt"
	"os"
	"walrus/frontend/ast"
	"walrus/frontend/lexer"
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
		stmt := parseNode(parser)

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
				Index:  0,
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
			p.MakeError(p.currentToken().StartPos.Line, p.FilePath, token, fmt.Sprintf("Unexpected '%s' at line %d", token.Value, token.StartPos.Line)).AddHint(fmt.Sprintf("How about trying '%s' instead?", expectedKind), TEXT_HINT).Display()
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