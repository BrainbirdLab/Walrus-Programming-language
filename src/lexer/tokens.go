package lexer

import "fmt"

// TokenKind represents the type of token
type TokenKind int

const (
	// Identifiers
	EOF TokenKind = iota
	NUMBER
	STRING
	IDENTIFIER

	// Delimiters
	OPEN_BRACKET
	CLOSE_BRACKET
	OPEN_CURLY
	CLOSE_CURLY
	OPEN_PAREN
	CLOSE_PAREN

	// Assignment operators
	ASSIGNMENT
	EQUALS
	NOT
	NOT_EQUALS

	// Comparison operators
	LESS
	LESS_EQUALS
	GREATER
	GREATER_EQUALS

	// Logical operators
	OR
	AND

	// Literals
	DOT
	DOT_DOT
	SEMI_COLON
	COLON
	QUESTION
	COMMA

	// Unary operators
	PLUS_PLUS
	MINUS_MINUS
	PLUS_EQUALS
	MINUS_EQUALS
	TIMES_EQUALS
	DEVIDE_EQUALS
	MODULO_EQUALS

	// Binary operators
	PLUS
	MINUS
	TIMES
	DIVIDE
	MODULO

	// Keywords
	LET
	CONST
	CLASS
	NEW
	USE
	FROM
	FUNCTION
	IF
	ELSE
	FOREACH
	WHILE
	FOR
	EXPORT
	TYPEOF
	IN

	// Special constants
	NULL
	TRUE
	FALSE
	PI
	E
)

// Reserved keywords
var reserved_lookup map[string]TokenKind = map[string]TokenKind{
	"let":     LET,
	"const":   CONST,
	"class":   CLASS,
	"new":     NEW,
	"use":  USE,
	"from":    FROM,
	"fn":      FUNCTION,
	"if":      IF,
	"else":    ELSE,
	"foreach": FOREACH,
	"while":   WHILE,
	"for":     FOR,
	"export":  EXPORT,
	"typeof":  TYPEOF,
	"in":      IN,
	"null":    NULL,
	"true":    TRUE,
	"false":   FALSE,
	"pi":      PI,
	"e":       E,
}

// Define token types
type Token struct {
	Kind  TokenKind
	Value string
}

func (token Token) isOneOfMany(expectedTokens ...TokenKind) bool {
	for _, expected := range expectedTokens {
		if expected == token.Kind {
			return true
		}
	}

	return false
}

// Debug prints a debug representation of the token
func (token Token) Debug() {
	if token.isOneOfMany(IDENTIFIER, NUMBER, STRING) {
		fmt.Printf("%s (%s)\n", TokenKindString(token.Kind), token.Value)
	} else {
		fmt.Printf("%s ()\n", TokenKindString(token.Kind))
	}
}

func NewToken(kind TokenKind, value string) Token {
	return Token{
		kind, value,
	}
}

// TokenKindString returns the string representation of a token kind
func TokenKindString(kind TokenKind) string {
	switch kind {
	case EOF:
		return "eof"
	case NUMBER:
		return "number"
	case STRING:
		return "string"
	case IDENTIFIER:
		return "identifier"
	case OPEN_BRACKET:
		return "["
	case CLOSE_BRACKET:
		return "]"
	case OPEN_CURLY:
		return "{"
	case CLOSE_CURLY:
		return "}"
	case OPEN_PAREN:
		return "("
	case CLOSE_PAREN:
		return ")"
	case ASSIGNMENT:
		return "="
	case EQUALS:
		return "is"
	case NOT:
		return "!"
	case NOT_EQUALS:
		return "!="
	case LESS:
		return "<"
	case LESS_EQUALS:
		return "<="
	case GREATER:
		return ">"
	case GREATER_EQUALS:
		return ">="
	case OR:
		return "or"
	case AND:
		return "and"
	case DOT:
		return "."
	case DOT_DOT:
		return ".."
	case SEMI_COLON:
		return ";"
	case COLON:
		return ":"
	case QUESTION:
		return "?"
	case COMMA:
		return ","
	case PLUS_PLUS:
		return "++"
	case MINUS_MINUS:
		return "--"
	case PLUS_EQUALS:
		return "+="
	case MINUS_EQUALS:
		return "-="
	case TIMES_EQUALS:
		return "*="
	case DEVIDE_EQUALS:
		return "/="
	case MODULO_EQUALS:
		return "%="
	case PLUS:
		return "+"
	case MINUS:
		return "-"
	case TIMES:
		return "*"
	case DIVIDE:
		return "/"
	case MODULO:
		return "%"
	case LET:
		return "let"
	case CONST:
		return "const"
	case CLASS:
		return "class"
	case NEW:
		return "new"
	case USE:
		return "use"
	case FROM:
		return "from"
	case FUNCTION:
		return "function"
	case IF:
		return "if"
	case ELSE:
		return "else"
	case FOREACH:
		return "foreach"
	case WHILE:
		return "while"
	case FOR:
		return "for"
	case EXPORT:
		return "export"
	case TYPEOF:
		return "typeof"
	case IN:
		return "in"
	default:
		return "unknown"
	}
}
