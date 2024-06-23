package lexer

import (
	"fmt"
	"regexp"
)

// TOKEN_KIND represents the type of token
type TOKEN_KIND string

const (
	// Identifiers
	EOF_TOKEN TOKEN_KIND = "eof"

	// Literals
	INTEGER_TOKEN   TOKEN_KIND = "integer"
	FLOATING_TOKEN  TOKEN_KIND = "float"
	STRING_TOKEN    TOKEN_KIND = "string"
	CHARACTER_TOKEN TOKEN_KIND = "charecter"
	BOOLEAN_TOKEN   TOKEN_KIND = "boolean"

	IDENTIFIER_TOKEN TOKEN_KIND = "identifier"
	RETURN_TOKEN     TOKEN_KIND = "return"

	// Delimiters
	OPEN_BRACKET_TOKEN  TOKEN_KIND = "["
	CLOSE_BRACKET_TOKEN TOKEN_KIND = "]"
	OPEN_CURLY_TOKEN    TOKEN_KIND = "{"
	CLOSE_CURLY_TOKEN   TOKEN_KIND = "}"
	OPEN_PAREN_TOKEN    TOKEN_KIND = "("
	CLOSE_PAREN_TOKEN   TOKEN_KIND = ")"

	// Assignment operators
	ASSIGNMENT_TOKEN TOKEN_KIND = "="
	WALRUS_TOKEN     TOKEN_KIND = ":="
	ARROW_TOKEN      TOKEN_KIND = "->"

	// Comparison operators
	EQUALS_TOKEN         TOKEN_KIND = "=="
	NOT_EQUALS_TOKEN     TOKEN_KIND = "!="
	LESS_TOKEN           TOKEN_KIND = "<"
	LESS_EQUALS_TOKEN    TOKEN_KIND = "<="
	GREATER_TOKEN        TOKEN_KIND = ">"
	GREATER_EQUALS_TOKEN TOKEN_KIND = ">="

	// Logical operators
	OR_TOKEN  TOKEN_KIND = "||"
	AND_TOKEN TOKEN_KIND = "&&"

	// Literals
	DOT_TOKEN        TOKEN_KIND = "."
	DOT_DOT_TOKEN    TOKEN_KIND = ".."
	SEMI_COLON_TOKEN TOKEN_KIND = ";"
	COLON_TOKEN      TOKEN_KIND = ":"
	QUESTION_TOKEN   TOKEN_KIND = "?"
	COMMA_TOKEN      TOKEN_KIND = ","

	// Unary operators
	NOT_TOKEN           TOKEN_KIND = "!"
	PLUS_PLUS_TOKEN     TOKEN_KIND = "++"
	MINUS_MINUS_TOKEN   TOKEN_KIND = "--"
	PLUS_EQUALS_TOKEN   TOKEN_KIND = "+="
	MINUS_EQUALS_TOKEN  TOKEN_KIND = "-="
	TIMES_EQUALS_TOKEN  TOKEN_KIND = "*="
	DIVIDE_EQUALS_TOKEN TOKEN_KIND = "/="
	MODULO_EQUALS_TOKEN TOKEN_KIND = "%="
	POWER_EQUALS_TOKEN  TOKEN_KIND = "^="

	// Binary operators
	PLUS_TOKEN   TOKEN_KIND = "+"
	MINUS_TOKEN  TOKEN_KIND = "-"
	TIMES_TOKEN  TOKEN_KIND = "*"
	DIVIDE_TOKEN TOKEN_KIND = "/"
	MODULO_TOKEN TOKEN_KIND = "%"
	POWER_TOKEN  TOKEN_KIND = "^"

	// Keywords
	LET_TOKEN      TOKEN_KIND = "let"
	CONST_TOKEN    TOKEN_KIND = "const"
	NEW_TOKEN      TOKEN_KIND = "new"
	MODULE_TOKEN   TOKEN_KIND = "module"
	IMPORT_TOKEN   TOKEN_KIND = "import"
	FROM_TOKEN     TOKEN_KIND = "from"
	FUNCTION_TOKEN TOKEN_KIND = "fn"

	SWITCH_TOKEN  TOKEN_KIND = "switch"
	CASE_TOKEN    TOKEN_KIND = "case"
	DEFAULT_TOKEN TOKEN_KIND = "default"

	BREAK_TOKEN    TOKEN_KIND = "break"
	CONTINUE_TOKEN TOKEN_KIND = "continue"

	IF_TOKEN      TOKEN_KIND = "if"
	ELSEIF_TOKEN  TOKEN_KIND = "elf"
	ELSE_TOKEN    TOKEN_KIND = "els"
	FOREACH_TOKEN TOKEN_KIND = "foreach"
	WHERE_TOKEN   TOKEN_KIND = "where"
	WHILE_TOKEN   TOKEN_KIND = "while"
	FOR_TOKEN     TOKEN_KIND = "for"
	EXPORT_TOKEN  TOKEN_KIND = "export"
	TYPEOF_TOKEN  TOKEN_KIND = "typeof"
	IN_TOKEN      TOKEN_KIND = "in"

	// Special constants
	NULL_TOKEN  TOKEN_KIND = "null"
	TRUE_TOKEN  TOKEN_KIND = "true"
	FALSE_TOKEN TOKEN_KIND = "false"

	// Other
	STRUCT_TOKEN    TOKEN_KIND = "struct"
	EMBED_TOKEN     TOKEN_KIND = "embed"
	TRAIT_TOKEN     TOKEN_KIND = "trait"
	IMPLEMENT_TOKEN TOKEN_KIND = "implement"
	OVERRIDE_TOKEN  TOKEN_KIND = "override"
	STATIC_TOKEN    TOKEN_KIND = "static"
	ACCESS_TOKEN    TOKEN_KIND = "access modifier"
	READONLY_TOKEN  TOKEN_KIND = "readonly"
)

// Reserved keywords
var reservedLookup map[string]TOKEN_KIND = map[string]TOKEN_KIND{
	"let":      LET_TOKEN,
	"const":    CONST_TOKEN,
	"new":      NEW_TOKEN,
	"mod":      MODULE_TOKEN,
	"import":   IMPORT_TOKEN,
	"from":     FROM_TOKEN,
	"fn":       FUNCTION_TOKEN,
	"switch":   SWITCH_TOKEN,
	"case":     CASE_TOKEN,
	"default":  DEFAULT_TOKEN,
	"break":    BREAK_TOKEN,
	"continue": CONTINUE_TOKEN,
	"if":       IF_TOKEN,
	"elf":      ELSEIF_TOKEN,
	"els":      ELSE_TOKEN,
	"foreach":  FOREACH_TOKEN,
	"where":    WHERE_TOKEN,
	"while":    WHILE_TOKEN,
	"for":      FOR_TOKEN,
	"export":   EXPORT_TOKEN,
	"typeof":   TYPEOF_TOKEN,
	"in":       IN_TOKEN,
	"null":     NULL_TOKEN,
	"true":     TRUE_TOKEN,
	"false":    FALSE_TOKEN,
	"struct":   STRUCT_TOKEN,
	"embed":    EMBED_TOKEN,
	"trait":    TRAIT_TOKEN,
	"impl":     IMPLEMENT_TOKEN,
	"override": OVERRIDE_TOKEN,
	"static":   STATIC_TOKEN,
	"pub":      ACCESS_TOKEN,
	"priv":     ACCESS_TOKEN,
	"readonly": READONLY_TOKEN,
	"ret":      RETURN_TOKEN,
}

func IsKeyword(tokenKind TOKEN_KIND) bool {
	_, ok := reservedLookup[string(tokenKind)]
	return ok
}

func IsNumber(tokenKind TOKEN_KIND) bool {
	regexp := regexp.MustCompile(`[0-9]+(?:\.[0-9]*)?`)
	return regexp.MatchString(string(tokenKind))
}

func IsBuiltInType(tokenKind TOKEN_KIND) bool {
	regexp := regexp.MustCompile(`i8|i16|i32|i64|i128|f32|f64|bool|str|chr`)
	return regexp.MatchString(string(tokenKind))
}

// Define token types
type Token struct {
	Kind     TOKEN_KIND
	Value    string
	StartPos Position
	EndPos   Position
	//LineNumber 	int
}

func (token Token) isOneOfMany(expectedTokens ...TOKEN_KIND) bool {
	for _, expected := range expectedTokens {
		if expected == token.Kind {
			return true
		}
	}

	return false
}

// Debug prints a debug representation of the token
func (token Token) Debug() {
	if token.isOneOfMany(IDENTIFIER_TOKEN, INTEGER_TOKEN, FLOATING_TOKEN, STRING_TOKEN) {
		fmt.Printf("%s (%s)\n", token.Kind, token.Value)
	} else {
		fmt.Printf("%s ()\n", token.Kind)
	}
}

func NewToken(kind TOKEN_KIND, value string, start Position, end Position) Token {
	return Token{
		kind, value, start, end,
	}
}
