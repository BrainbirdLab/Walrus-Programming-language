package lexer

import (
	"fmt"
	"regexp"
)

// TOKEN_KIND represents the type of token
type TOKEN_KIND string

const (
	// Identifiers
	EOF 					TOKEN_KIND = "eof"

	NUMBER                	TOKEN_KIND = "number"
	STRING                	TOKEN_KIND = "string"
	CHARACTER             	TOKEN_KIND = "charecter"
	IDENTIFIER            	TOKEN_KIND = "identifier"
	RETURN                	TOKEN_KIND = "return"

	// Delimiters
	OPEN_BRACKET  			TOKEN_KIND = "["
	CLOSE_BRACKET 			TOKEN_KIND = "]"
	OPEN_CURLY    			TOKEN_KIND = "{"
	CLOSE_CURLY   			TOKEN_KIND = "}"
	OPEN_PAREN    			TOKEN_KIND = "("
	CLOSE_PAREN   			TOKEN_KIND = ")"

	// Assignment operators
	ASSIGNMENT 				TOKEN_KIND = "="
	WALRUS     				TOKEN_KIND = ":="
	ARROW      				TOKEN_KIND = "->"

	// Comparison operators
	EQUALS         			TOKEN_KIND = "=="
	NOT_EQUALS     			TOKEN_KIND = "!="
	LESS           			TOKEN_KIND = "<"
	LESS_EQUALS    			TOKEN_KIND = "<="
	GREATER        			TOKEN_KIND = ">"
	GREATER_EQUALS 			TOKEN_KIND = ">="

	// Logical operators
	OR  					TOKEN_KIND = "||"
	AND 					TOKEN_KIND = "&&"

	// Literals
	DOT        				TOKEN_KIND = "."
	DOT_DOT    				TOKEN_KIND = ".."
	SEMI_COLON 				TOKEN_KIND = ";"
	COLON      				TOKEN_KIND = ":"
	QUESTION   				TOKEN_KIND = "?"
	COMMA      				TOKEN_KIND = ","

	// Unary operators
	NOT           			TOKEN_KIND = "!"
	PLUS_PLUS     			TOKEN_KIND = "++"
	MINUS_MINUS   			TOKEN_KIND = "--"
	PLUS_EQUALS   			TOKEN_KIND = "+="
	MINUS_EQUALS  			TOKEN_KIND = "-="
	TIMES_EQUALS  			TOKEN_KIND = "*="
	DIVIDE_EQUALS 			TOKEN_KIND = "/="
	MODULO_EQUALS 			TOKEN_KIND = "%="

	// Binary operators
	PLUS   					TOKEN_KIND = "+"
	MINUS  					TOKEN_KIND = "-"
	TIMES  					TOKEN_KIND = "*"
	DIVIDE 					TOKEN_KIND = "/"
	MODULO 					TOKEN_KIND = "%"
	POWER  					TOKEN_KIND = "^"

	// Keywords
	LET      				TOKEN_KIND = "let"
	CONST    				TOKEN_KIND = "const"
	NEW      				TOKEN_KIND = "new"
	MODULE					TOKEN_KIND = "module"
	IMPORT   				TOKEN_KIND = "import"
	FROM     				TOKEN_KIND = "from"
	FUNCTION 				TOKEN_KIND = "fn"

	SWITCH					TOKEN_KIND = "switch"
	CASE					TOKEN_KIND = "case"
	DEFAULT					TOKEN_KIND = "default"

	BREAK					TOKEN_KIND = "break"
	CONTINUE				TOKEN_KIND = "continue"

	IF       				TOKEN_KIND = "if"
	ELSEIF   				TOKEN_KIND = "elf"
	ELSE     				TOKEN_KIND = "els"
	FOREACH  				TOKEN_KIND = "foreach"
	WHERE					TOKEN_KIND = "where"
	WHILE    				TOKEN_KIND = "while"
	FOR      				TOKEN_KIND = "for"
	EXPORT   				TOKEN_KIND = "export"
	TYPEOF   				TOKEN_KIND = "typeof"
	IN       				TOKEN_KIND = "in"

	// Special constants
	NULL  					TOKEN_KIND = "null"
	TRUE  					TOKEN_KIND = "true"
	FALSE 					TOKEN_KIND = "false"

	// Other
	STRUCT   				TOKEN_KIND = "struct"
	EMBED	  				TOKEN_KIND = "embed"
	TRAIT	  				TOKEN_KIND = "trait"
	IMPLEMENT 				TOKEN_KIND = "implement"
	OVERRIDE 				TOKEN_KIND = "override"
	STATIC   				TOKEN_KIND = "static"
	ACCESS   				TOKEN_KIND = "access modifier"
	READONLY 				TOKEN_KIND = "readonly"
)

// Reserved keywords
var reservedLookup map[string]TOKEN_KIND = map[string]TOKEN_KIND{
	"let":      	LET,
	"const":    	CONST,
	"new":      	NEW,
	"mod":   		MODULE,
	"import":   	IMPORT,
	"from":     	FROM,
	"fn":       	FUNCTION,
	"switch":		SWITCH,
	"case":			CASE,
	"default":		DEFAULT,
	"break":		BREAK,
	"continue":		CONTINUE,
	"if":       	IF,
	"elf":      	ELSEIF,
	"els":      	ELSE,
	"foreach":  	FOREACH,
	"where":    	WHERE,
	"while":    	WHILE,
	"for":      	FOR,
	"export":   	EXPORT,
	"typeof":   	TYPEOF,
	"in":       	IN,
	"null":     	NULL,
	"true":     	TRUE,
	"false":    	FALSE,
	"struct":   	STRUCT,
	"embed":  		EMBED,
	"trait":		TRAIT,
	"impl":			IMPLEMENT,
	"override": 	OVERRIDE,
	"static":   	STATIC,
	"pub":      	ACCESS,
	"priv":     	ACCESS,
	"readonly": 	READONLY,
	"ret":      	RETURN,
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

	if token.isOneOfMany(IDENTIFIER, NUMBER, STRING) {
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