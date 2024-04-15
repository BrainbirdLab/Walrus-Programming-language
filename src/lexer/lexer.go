package lexer

import (
	"fmt"
	"regexp"
	//"github.com/sanity-io/litter"
)

type regexHandler func(lex *Lexer, regex *regexp.Regexp)

type regexPattern struct {
	regex   *regexp.Regexp
	handler regexHandler
}

type Lexer struct {
	patterns []regexPattern
	Tokens   []Token
	source   string
	index      int
}

func Tokenize(source string, debug bool) []Token {
	
	lex := createLexer(source)

	for !lex.at_eof() {

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
			panic(fmt.Sprintf("Lexer::Error at position %d: Unexpected token near %v\n", lex.index, lex.remainingLines()))
		}
	}

	lex.push(NewToken(EOF, "EOF"))

	//litter.Dump(lex.Tokens)
	if (debug) {
		for _, token := range lex.Tokens {
			token.Debug()
		}
	}

	return lex.Tokens
}

func (lex *Lexer) advanceN(n int) {
	lex.index += n
}

func (lex *Lexer) push(token Token) {
	lex.Tokens = append(lex.Tokens, token)
}

func (lex *Lexer) at() byte {
	return lex.source[lex.index]
}

func (lex *Lexer) remainder() string {
	return lex.source[lex.index:]
}

func (lex *Lexer) remainingLines() string {
	//until newline or eof
	rem := lex.remainder()
	for i, c := range rem {
		if c == '\n' {
			return rem[:i]
		}
	}
	return rem
}

func (lex *Lexer) at_eof() bool {
	return lex.index >= len(lex.source)
}

func defaultHandler(kind TokenKind, value string) regexHandler {
	return func(lex *Lexer, regex *regexp.Regexp) {
		// advance the position
		lex.advanceN(len(value))
		lex.push(NewToken(kind, value))
	}
}

func createLexer(source string) *Lexer {
	return &Lexer{
		index:    0,
		source: source,
		Tokens: make([]Token, 0),
		patterns: []regexPattern{
			{regexp.MustCompile(`\s+`), skipHandler}, // whitespace
			{regexp.MustCompile(`\/\/.*`), skipHandler}, // single line comments
			{regexp.MustCompile(`\/\*[\s\S]*?\*\/`), skipHandler}, // multi line comments
			{regexp.MustCompile(`"[^"]*"`), stringHandler}, // string literals
			{regexp.MustCompile(`[0-9.]+`), numberHandler}, // numbers
			{regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`), symbolHandler}, // identifiers
			{regexp.MustCompile(`\[`), defaultHandler(OPEN_BRACKET, "[")},
			{regexp.MustCompile(`\]`), defaultHandler(CLOSE_BRACKET, "]")},
			{regexp.MustCompile(`\{`), defaultHandler(OPEN_CURLY, "{")},
			{regexp.MustCompile(`\}`), defaultHandler(CLOSE_CURLY, "}")},
			{regexp.MustCompile(`\(`), defaultHandler(OPEN_PAREN, "(")},
			{regexp.MustCompile(`\)`), defaultHandler(CLOSE_PAREN, ")")},
			{regexp.MustCompile(`==`), defaultHandler(EQUALS, "==")},
			{regexp.MustCompile(`!=`), defaultHandler(NOT_EQUALS, "!=")},
			{regexp.MustCompile(`=`), defaultHandler(ASSIGNMENT, "=")},
			{regexp.MustCompile(`!`), defaultHandler(NOT, "!")},
			{regexp.MustCompile(`<=`), defaultHandler(LESS_EQUALS, "<=")},
			{regexp.MustCompile(`<`), defaultHandler(LESS, "<")},
			{regexp.MustCompile(`>=`), defaultHandler(GREATER_EQUALS, ">=")},
			{regexp.MustCompile(`>`), defaultHandler(GREATER, ">")},
			{regexp.MustCompile(`\|\|`), defaultHandler(OR, "||")},
			{regexp.MustCompile(`&&`), defaultHandler(AND, "&&")},
			{regexp.MustCompile(`\.\.`), defaultHandler(DOT_DOT, "..")},
			{regexp.MustCompile(`\.`), defaultHandler(DOT, ".")},
			{regexp.MustCompile(`;`), defaultHandler(SEMI_COLON, ";")},
			{regexp.MustCompile(`:`), defaultHandler(COLON, ":")},
			//{regexp.MustCompile(`\?\?=`), defaultHandler(NULLISH_ASSIGNMENT, "??=")},
			{regexp.MustCompile(`\?`), defaultHandler(QUESTION, "?")},
			{regexp.MustCompile(`,`), defaultHandler(COMMA, ",")},
			{regexp.MustCompile(`\+\+`), defaultHandler(PLUS_PLUS, "++")},
			{regexp.MustCompile(`--`), defaultHandler(MINUS_MINUS, "--")},
			{regexp.MustCompile(`\+=`), defaultHandler(PLUS_EQUALS, "+=")},
			{regexp.MustCompile(`-=`), defaultHandler(MINUS_EQUALS, "-=")},
			{regexp.MustCompile(`\*=`), defaultHandler(TIMES_EQUALS, "*=")},
			{regexp.MustCompile(`/=`), defaultHandler(DIVIDE_EQUALS, "/=")},
			{regexp.MustCompile(`%=`), defaultHandler(MODULO_EQUALS, "%=")},
			{regexp.MustCompile(`\+`), defaultHandler(PLUS, "+")},
			{regexp.MustCompile(`-`), defaultHandler(MINUS, "-")},
			{regexp.MustCompile(`/`), defaultHandler(DIVIDE, "/")},
			{regexp.MustCompile(`\*`), defaultHandler(TIMES, "*")},
			{regexp.MustCompile(`%`), defaultHandler(MODULO, "%")},
		},
	}
}

func symbolHandler(lex *Lexer, regex *regexp.Regexp) {
	symbol := regex.FindString(lex.remainder())

	if kind, exists := reserved_lookup[symbol]; exists {
		lex.push(NewToken(kind, symbol))
	} else {
		lex.push(NewToken(IDENTIFIER, symbol))
	}

	lex.advanceN(len(symbol))
}

func numberHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	// if the number has more than 1 dot, it's an error
	count := 0
	for _, c := range match {
		if c == '.' {
			count++
		}
	}
	
	if count > 1 {
		panic(fmt.Sprintf("Lexer::Error at position %d: Multiple decimal operator found near %v\n", lex.index, lex.remainder()))
	}

	lex.push(NewToken(NUMBER, match))
	lex.advanceN(len(match))
}

func stringHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	stringLiteral := lex.remainder()[match[0]+1 : match[1]-1]

	lex.push(NewToken(STRING, stringLiteral))
	lex.advanceN(len(stringLiteral) + 2)
}

func skipHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	lex.advanceN(match[1])
}