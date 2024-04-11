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
	DEVIDE
	MODULO

	// Keywords
	LET
	CONST
	CLASS
	NEW
	IMPORT
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
var reserved_lookup map[string]TokenKind = map[string]TokenKind {
	"let": LET,
	"const": CONST,
	"class": CLASS,
	"new": NEW,
	"import": IMPORT,
	"from": FROM,
	"fn": FUNCTION,
	"if": IF,
	"else": ELSE,
	"foreach": FOREACH,
	"while": WHILE,
	"for": FOR,
	"export": EXPORT,
	"typeof": TYPEOF,
	"in": IN,
	"null": NULL,
	"true": TRUE,
	"false": FALSE,
	"pi": PI,
	"e": E,
}

// Define token types
type Token struct {
	Kind  TokenKind
	Value string
}

func (token Token) isOneOfMany (expectedTokens ...TokenKind) bool {
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

func NewToken (kind TokenKind, value string) Token {
	return Token {
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
		return "open_bracket"
	case CLOSE_BRACKET:
		return "close_bracket"
	case OPEN_CURLY:
		return "open_curly"
	case CLOSE_CURLY:
		return "close_curly"
	case OPEN_PAREN:
		return "open_paren"
	case CLOSE_PAREN:
		return "close_paren"
	case ASSIGNMENT:
		return "assignment"
	case EQUALS:
		return "equals"
	case NOT:
		return "not"
	case NOT_EQUALS:
		return "not_equals"
	case LESS:
		return "less"
	case LESS_EQUALS:
		return "less_equals"
	case GREATER:
		return "greater"
	case GREATER_EQUALS:
		return "greater_equals"
	case OR:
		return "or"
	case AND:
		return "and"
	case DOT:
		return "dot"
	case DOT_DOT:
		return "dot_dot"
	case SEMI_COLON:
		return "semi_colon"
	case COLON:
		return "colon"
	case QUESTION:
		return "question"
	case COMMA:
		return "comma"
	case PLUS_PLUS:
		return "plus_plus"
	case MINUS_MINUS:
		return "minus_minus"
	case PLUS_EQUALS:
		return "plus_equals"
	case MINUS_EQUALS:
		return "minus_equals"
	case TIMES_EQUALS:
		return "times_equals"
	case DEVIDE_EQUALS:
		return "devide_equals"
	case MODULO_EQUALS:
		return "modulo_equals"
	case PLUS:
		return "plus"
	case MINUS:
		return "minus"
	case TIMES:
		return "times"
	case DEVIDE:
		return "devide"
	case MODULO:
		return "modulo"
	case LET:
		return "let"
	case CONST:
		return "const"
	case CLASS:
		return "class"
	case NEW:
		return "new"
	case IMPORT:
		return "import"
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