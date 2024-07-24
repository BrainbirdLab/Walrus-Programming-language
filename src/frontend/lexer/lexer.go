package lexer

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"walrus/utils"
)

type regexHandler func(lex *Lexer, regex *regexp.Regexp)

type regexPattern struct {
	regex   *regexp.Regexp
	handler regexHandler
}

type Position struct {
	Line   int
	Column int
	Index  int
}

func (p *Position) advance(toSkip string) *Position {

	for _, char := range toSkip { // a rune is an alias for int32
		if char == '\n' {
			p.Line++
			p.Column = 1
		} else {
			p.Column++
		}
		p.Index++
	}

	return p
}

type Lexer struct {
	patterns []regexPattern
	Tokens   []Token
	Lines    []string
	source   *string
	Pos      Position
	FilePath string
}

func Tokenize(source, file string, debug bool) ([]Token, *[]string) {

	lex := createLexer(&source)
	lex.FilePath = file
	lex.Lines = strings.Split(source, "\n")

	for !lex.atEOF() {

		matched := false

		for _, pattern := range lex.patterns {

			loc := pattern.regex.FindStringIndex(lex.remainder())

			if loc != nil && loc[0] == 0 {
				pattern.handler(lex, pattern.regex)
				matched = true
				break
			}
		}

		if !matched {

			//line is from lex.Pos.Index to
			padding := fmt.Sprintf("%d | ", lex.Pos.Line)

			errStr := fmt.Sprintf("\n%s:%d:%d\n", lex.FilePath, lex.Pos.Line, lex.Pos.Column)
			errStr += utils.Colorize(utils.GREY, padding) + Highlight(lex.Lines[lex.Pos.Line]) + "\n"
			errStr += utils.Colorize(utils.BOLD_RED, (strings.Repeat(" ", (lex.Pos.Column-1)+len(padding)) + "^\n"))
			errStr += fmt.Sprintf("At line %d: Unexpected character: '%c'", lex.Pos.Line, lex.at())
			fmt.Println(errStr)

			os.Exit(-1)
		}
	}

	lex.push(NewToken(EOF_TOKEN, "End of file", lex.Pos, lex.Pos))

	//litter.Dump(lex.Tokens)
	if debug {
		for _, token := range lex.Tokens {
			token.Debug()
		}
	}

	return lex.Tokens, &lex.Lines
}

func (lex *Lexer) advanceN(match string) {
	//ascii value of match-
	lex.Pos.advance(match)
}

func (lex *Lexer) push(token Token) {
	lex.Tokens = append(lex.Tokens, token)
}

func (lex *Lexer) at() byte {
	return (*(lex.source))[lex.Pos.Index]
}

func (lex *Lexer) remainder() string {
	return (*(lex.source))[lex.Pos.Index:]
}

func (lex *Lexer) atEOF() bool {
	return lex.Pos.Index >= len(*(lex.source))
}

func defaultHandler(kind TOKEN_KIND, value string) regexHandler {

	return func(lex *Lexer, _ *regexp.Regexp) {

		start := lex.Pos
		lex.advanceN(value)
		end := lex.Pos

		lex.push(NewToken(kind, value, start, end))
	}
}

func createLexer(source *string) *Lexer {
	lex := &Lexer{
		source: source,
		Tokens: make([]Token, 0),
		Pos: Position{
			Line:   1,
			Column: 1,
			Index:  0,
		},
		patterns: []regexPattern{
			//{regexp.MustCompile(`\n`), skipHandler}, // newlines
			{regexp.MustCompile(`\s+`), skipHandler},                          // whitespace
			{regexp.MustCompile(`\/\/.*`), skipHandler},                       // single line comments
			{regexp.MustCompile(`\/\*[\s\S]*?\*\/`), skipHandler},             // multi line comments
			{regexp.MustCompile(`"[^"]*"`), stringHandler},                    // string literals
			{regexp.MustCompile(`'[^']'`), characterHandler},                  // character literals
			{regexp.MustCompile(`[0-9]+(?:\.[0-9]+)?`), numberHandler},        // decimal numbers
			{regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`), identifierHandler}, // identifiers
			{regexp.MustCompile(`\[`), defaultHandler(OPEN_BRACKET_TOKEN, "[")},
			{regexp.MustCompile(`\]`), defaultHandler(CLOSE_BRACKET_TOKEN, "]")},
			{regexp.MustCompile(`\{`), defaultHandler(OPEN_CURLY_TOKEN, "{")},
			{regexp.MustCompile(`\}`), defaultHandler(CLOSE_CURLY_TOKEN, "}")},
			{regexp.MustCompile(`\(`), defaultHandler(OPEN_PAREN_TOKEN, "(")},
			{regexp.MustCompile(`\)`), defaultHandler(CLOSE_PAREN_TOKEN, ")")},
			{regexp.MustCompile(`==`), defaultHandler(EQUALS_TOKEN, "==")},
			{regexp.MustCompile(`!=`), defaultHandler(NOT_EQUALS_TOKEN, "!=")},
			{regexp.MustCompile(`=`), defaultHandler(ASSIGNMENT_TOKEN, "=")},
			{regexp.MustCompile(`:=`), defaultHandler(WALRUS_TOKEN, ":=")},
			{regexp.MustCompile(`!`), defaultHandler(NOT_TOKEN, "!")},
			{regexp.MustCompile(`<=`), defaultHandler(LESS_EQUALS_TOKEN, "<=")},
			{regexp.MustCompile(`<`), defaultHandler(LESS_TOKEN, "<")},
			{regexp.MustCompile(`>=`), defaultHandler(GREATER_EQUALS_TOKEN, ">=")},
			{regexp.MustCompile(`>`), defaultHandler(GREATER_TOKEN, ">")},
			{regexp.MustCompile(`\|\|`), defaultHandler(OR_TOKEN, "||")},
			{regexp.MustCompile(`&&`), defaultHandler(AND_TOKEN, "&&")},
			{regexp.MustCompile(`\.\.`), defaultHandler(DOT_DOT_TOKEN, "..")},
			{regexp.MustCompile(`\.`), defaultHandler(DOT_TOKEN, ".")},
			{regexp.MustCompile(`;`), defaultHandler(SEMI_COLON_TOKEN, ";")},
			{regexp.MustCompile(`:`), defaultHandler(COLON_TOKEN, ":")},
			//{regexp.MustCompile(`\?\?=`), defaultHandler(NULLISH_ASSIGNMENT, "??=")},
			{regexp.MustCompile(`->`), defaultHandler(ARROW_TOKEN, "->")},
			{regexp.MustCompile(`\?`), defaultHandler(QUESTION_TOKEN, "?")},
			{regexp.MustCompile(`,`), defaultHandler(COMMA_TOKEN, ",")},
			{regexp.MustCompile(`\+\+`), defaultHandler(PLUS_PLUS_TOKEN, "++")},
			{regexp.MustCompile(`--`), defaultHandler(MINUS_MINUS_TOKEN, "--")},
			{regexp.MustCompile(`\+=`), defaultHandler(PLUS_EQUALS_TOKEN, "+=")},
			{regexp.MustCompile(`-=`), defaultHandler(MINUS_EQUALS_TOKEN, "-=")},
			{regexp.MustCompile(`\*=`), defaultHandler(TIMES_EQUALS_TOKEN, "*=")},
			{regexp.MustCompile(`/=`), defaultHandler(DIVIDE_EQUALS_TOKEN, "/=")},
			{regexp.MustCompile(`%=`), defaultHandler(MODULO_EQUALS_TOKEN, "%=")},
			{regexp.MustCompile(`\^=`), defaultHandler(POWER_EQUALS_TOKEN, "^=")},
			{regexp.MustCompile(`\+`), defaultHandler(PLUS_TOKEN, "+")},
			{regexp.MustCompile(`-`), defaultHandler(MINUS_TOKEN, "-")},
			{regexp.MustCompile(`/`), defaultHandler(DIVIDE_TOKEN, "/")},
			{regexp.MustCompile(`\*`), defaultHandler(TIMES_TOKEN, "*")},
			{regexp.MustCompile(`%`), defaultHandler(MODULO_TOKEN, "%")},
			{regexp.MustCompile(`^`), defaultHandler(POWER_TOKEN, "^")},
		},
	}
	return lex
}

func identifierHandler(lex *Lexer, regex *regexp.Regexp) {

	identifier := regex.FindString(lex.remainder())

	start := lex.Pos
	lex.advanceN(identifier)
	end := lex.Pos

	if kind, exists := reservedLookup[identifier]; exists {
		lex.push(NewToken(kind, identifier, start, end))
	} else {
		lex.push(NewToken(IDENTIFIER_TOKEN, identifier, start, end))
	}

}

func numberHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())

	start := lex.Pos

	lex.advanceN(match)

	end := lex.Pos

	//find the number is a float or an integer
	if strings.Contains(match, ".") {
		lex.push(NewToken(FLOATING_TOKEN, match, start, end))
	} else {
		lex.push(NewToken(INTEGER_TOKEN, match, start, end))
	}
}

func stringHandler(lex *Lexer, regex *regexp.Regexp) {

	match := regex.FindString(lex.remainder())

	//exclude the quotes
	stringLiteral := match[1 : len(match)-1]

	start := lex.Pos
	lex.advanceN(match)
	end := lex.Pos

	lex.push(NewToken(STRING_TOKEN, stringLiteral, start, end))
}

func characterHandler(lex *Lexer, regex *regexp.Regexp) {

	match := regex.FindString(lex.remainder())

	//exclude the quotes
	characterLiteral := match[1 : len(match)-1]

	start := lex.Pos
	lex.advanceN(match)
	end := lex.Pos

	lex.push(NewToken(CHARACTER_TOKEN, characterLiteral, start, end))
}

func skipHandler(lex *Lexer, regex *regexp.Regexp) {

	match := regex.FindString(lex.remainder())

	lex.advanceN(match)
}
