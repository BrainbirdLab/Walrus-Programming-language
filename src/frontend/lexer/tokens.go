package lexer

import (
	"fmt"
	"regexp"
)

// TOKEN_KIND represents the type of token
type TOKEN_KIND string

const (
	// Identifiers
	EOF TOKEN_KIND 			= "eof"

	NUMBER                	= "number"
	STRING                	= "string"
	CHARACTER             	= "charecter"
	IDENTIFIER            	= "identifier"
	RETURN                	= "return"

	// Delimiters
	OPEN_BRACKET  			= "["
	CLOSE_BRACKET 			= "]"
	OPEN_CURLY    			= "{"
	CLOSE_CURLY   			= "}"
	OPEN_PAREN    			= "("
	CLOSE_PAREN   			= ")"

	// Assignment operators
	ASSIGNMENT 				= "="
	WALRUS     				= ":="
	ARROW      				= "->"

	// Comparison operators
	EQUALS         			= "=="
	NOT_EQUALS     			= "!="
	LESS           			= "<"
	LESS_EQUALS    			= "<="
	GREATER        			= ">"
	GREATER_EQUALS 			= ">="

	// Logical operators
	OR  					= "||"
	AND 					= "&&"

	// Literals
	DOT        				= "."
	DOT_DOT    				= ".."
	SEMI_COLON 				= ";"
	COLON      				= ":"
	QUESTION   				= "?"
	COMMA      				= ","

	// Unary operators
	NOT           			= "!"
	PLUS_PLUS     			= "++"
	MINUS_MINUS   			= "--"
	PLUS_EQUALS   			= "+="
	MINUS_EQUALS  			= "-="
	TIMES_EQUALS  			= "*="
	DIVIDE_EQUALS 			= "/="
	MODULO_EQUALS 			= "%="

	// Binary operators
	PLUS   					= "+"
	MINUS  					= "-"
	TIMES  					= "*"
	DIVIDE 					= "/"
	MODULO 					= "%"
	POWER  					= "^"

	// Keywords
	LET      				= "let"
	CONST    				= "const"
	NEW      				= "new"
	MODULE					= "module"
	IMPORT   				= "import"
	FROM     				= "from"
	FUNCTION 				= "fn"

	SWITCH					= "switch"
	CASE					= "case"
	DEFAULT					= "default"

	BREAK					= "break"
	CONTINUE				= "continue"

	IF       				= "if"
	ELSEIF   				= "elf"
	ELSE     				= "els"
	FOREACH  				= "foreach"
	WHERE					= "where"
	WHILE    				= "while"
	FOR      				= "for"
	EXPORT   				= "export"
	TYPEOF   				= "typeof"
	IN       				= "in"

	// Special constants
	NULL  					= "null"
	TRUE  					= "true"
	FALSE 					= "false"

	// Other
	STRUCT   				= "struct"
	EMBED	  				= "embed"
	TRAIT	  				= "trait"
	IMPLEMENT 				= "implement"
	OVERRIDE 				= "override"
	STATIC   				= "static"
	ACCESS   				= "access modifier"
	READONLY 				= "readonly"
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

	if (token.Value == DOT_DOT) {
		fmt.Println("Found..................")
	}

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