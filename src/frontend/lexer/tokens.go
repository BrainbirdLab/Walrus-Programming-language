package lexer

import "fmt"

// TokenKind represents the type of token
type TokenKind int

const (
	// Identifiers
	EOF TokenKind = iota
	NUMBER
	STRING
	CHARACTER
	IDENTIFIER
	RETURN

	// Delimiters
	OPEN_BRACKET
	CLOSE_BRACKET
	OPEN_CURLY
	CLOSE_CURLY
	OPEN_PAREN
	CLOSE_PAREN

	// Assignment operators
	ASSIGNMENT
	WALRUS

	ARROW


	// Comparison operators
	EQUALS
	NOT_EQUALS
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
	//ENDLINE
	COLON
	QUESTION
	COMMA
	
	// Unary operators
	NOT
	PLUS_PLUS
	MINUS_MINUS
	PLUS_EQUALS
	MINUS_EQUALS
	TIMES_EQUALS
	DIVIDE_EQUALS
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
	NEW
	IMPORT
	FROM
	FUNCTION
	IF
	ELSEIF
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

	STRUCT
	STATIC
	ACCESS_MODIFIER
	READONLY
)

// Reserved keywords
var reserved_lookup map[string]TokenKind = map[string]TokenKind{
	"let":     		LET,
	"const":   		CONST,
	"new":     		NEW,
	"import":		IMPORT,
	"from":    		FROM,
	"fn":      		FUNCTION,
	"if":      		IF,
	"elf":     		ELSEIF,
	"els":     		ELSE,
	"foreach": 		FOREACH,
	"while":   		WHILE,
	"for":     		FOR,
	"export":  		EXPORT,
	"typeof":  		TYPEOF,
	"in":      		IN,
	"null":    		NULL,
	"true":    		TRUE,
	"false":   		FALSE,
	"struct":  		STRUCT,
	"static":  		STATIC,
	"pub":    		ACCESS_MODIFIER,
	"priv":    		ACCESS_MODIFIER,
	"readonly":		READONLY,
	"ret":    		RETURN,
}

// Location
type Location struct {
	Line   int
	Column int
	Endline int
	Endcolumn int
}

// Define token types
type Token struct {
	Kind  		TokenKind
	Value 		string
	StartPos 	Position
	EndPos   	Position
	//LineNumber 	int
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

func NewToken(kind TokenKind, value string, start Position, end Position) Token {
	return Token{
		kind, value, start, end,
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
		return "Open Bracket ["
	case CLOSE_BRACKET:
		return "Close Bracket ]"
	case OPEN_CURLY:
		return "Open curly {"
	case CLOSE_CURLY:
		return "Close curly }"
	case OPEN_PAREN:
		return "Open paren ("
	case CLOSE_PAREN:
		return "Close paren )"
	case ASSIGNMENT:
		return "Assignment ="
	case WALRUS:
		return "Walrus :="
	case EQUALS:
		return "is =="
	case NOT:
		return "Not !"
	case NOT_EQUALS:
		return "Not equal !="
	case LESS:
		return "Less than <"
	case LESS_EQUALS:
		return "Less than or equal <="
	case GREATER:
		return "Greater than >"
	case GREATER_EQUALS:
		return "Greater than or equal >="
	case OR:
		return "Or ||"
	case AND:
		return "And &&"
	case DOT:
		return "Dot ."
	case DOT_DOT:
		return "Range .."
	case SEMI_COLON:
		return "semi colon ;"
	//case ENDLINE:
	//	return "Endline"
	case COLON:
		return "Colon :"
	case QUESTION:
		return " Question ?"
	case COMMA:
		return "Comma ,"
	case PLUS_PLUS:
		return "Increment ++"
	case MINUS_MINUS:
		return "Decrement --"
	case PLUS_EQUALS:
		return "Incremental assignment +="
	case MINUS_EQUALS:
		return "Decremental assignment -="
	case TIMES_EQUALS:
		return "Multiplicative assignment *="
	case DIVIDE_EQUALS:
		return "Division assignment /="
	case MODULO_EQUALS:
		return "Modulo assignment %="
	case PLUS:
		return "Add +"
	case MINUS:
		return "Subtract -"
	case TIMES:
		return "Multiply *"
	case DIVIDE:
		return "Divide /"
	case MODULO:
		return "Modulo %"
	case LET:
		return "let"
	case CONST:
		return "const"
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
	case ELSEIF:
		return "elseif"
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
	case NULL:
		return "null"
	case TRUE:
		return "true"
	case FALSE:
		return "false"
	case STRUCT:
		return "struct"
	case STATIC:
		return "static"
	case ARROW:
		return "->"
	case ACCESS_MODIFIER:
		return "access modifier"
	case READONLY:
		return "readonly"
	case RETURN:
		return "return"
	default:
		return "unknown"
	}
}