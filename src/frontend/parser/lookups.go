package parser

import (
	"walrus/frontend/ast"
	"walrus/frontend/lexer"
)

type BINDING_POWER int

// Operator precedence levels. PRIMARY is the highest binding power
const (
	DEFAULT_BP BINDING_POWER = iota
	COMMA
	ASSIGNMENT
	LOGICAL
	RELATIONAL
	ADDITIVE
	MULTIPLICATIVE
	UNARY
	CALL
	MEMBER
	PRIMARY
)

// Statement handler. Expect nothing to the left of the token
type statementHandler func(p *Parser) ast.Statement

// Null denoted. Expect nothing to the left of the token
type NUDHandler func(p *Parser) ast.Expression

// Left denoted. Expect something to the left of the token
type LEDHandler func(p *Parser, left ast.Expression, bp BINDING_POWER) ast.Expression

// lookup table for the different token types
type stmtLookupType map[lexer.TOKEN_KIND]statementHandler
type nudLookupType map[lexer.TOKEN_KIND]NUDHandler
type ledLookupType map[lexer.TOKEN_KIND]LEDHandler
type bpLookupType map[lexer.TOKEN_KIND]BINDING_POWER

var nudLookup = nudLookupType{}
var ledLookup = ledLookupType{}
var stmtLookup = stmtLookupType{}
var bpLookupMap = bpLookupType{}

func GetBP(kind lexer.TOKEN_KIND) BINDING_POWER {
	if bp, ok := bpLookupMap[kind]; ok {
		return bp
	} else {
		return DEFAULT_BP
	}
}

func led(kind lexer.TOKEN_KIND, bp BINDING_POWER, led_fn LEDHandler) {
	bpLookupMap[kind] = bp
	ledLookup[kind] = led_fn
}

func nud(kind lexer.TOKEN_KIND, nud_fn NUDHandler) {
	nudLookup[kind] = nud_fn
}

func stmt(kind lexer.TOKEN_KIND, stmt_fn statementHandler) {
	stmtLookup[kind] = stmt_fn
}

func createTokenLookups() {

	//literals & Identifiers
	nud(lexer.NUMBER, parsePrimaryExpr)
	nud(lexer.STRING, parsePrimaryExpr)
	nud(lexer.IDENTIFIER, parsePrimaryExpr)
	nud(lexer.TRUE, parsePrimaryExpr)
	nud(lexer.FALSE, parsePrimaryExpr)
	nud(lexer.NULL, parsePrimaryExpr)

	nud(lexer.OPEN_PAREN, parseGroupingExpr)

	//unary / prefix
	nud(lexer.MINUS, parsePrefixExpr)
	nud(lexer.PLUS, parsePrefixExpr)
	nud(lexer.PLUS_PLUS, parseUnaryExpr)
	nud(lexer.MINUS_MINUS, parseUnaryExpr)
	nud(lexer.NOT, parseUnaryExpr)
	nud(lexer.OPEN_BRACKET, parseArrayExpr)

	// Assignment
	led(lexer.ASSIGNMENT, ASSIGNMENT, parseVarAssignmentExpr)
	led(lexer.PLUS_EQUALS, ASSIGNMENT, parseVarAssignmentExpr)
	led(lexer.MINUS_EQUALS, ASSIGNMENT, parseVarAssignmentExpr)
	led(lexer.TIMES_EQUALS, ASSIGNMENT, parseVarAssignmentExpr)
	led(lexer.DIVIDE_EQUALS, ASSIGNMENT, parseVarAssignmentExpr)
	led(lexer.MODULO_EQUALS, ASSIGNMENT, parseVarAssignmentExpr)

	// Logical operations
	led(lexer.AND, LOGICAL, parseBinaryExpr)
	led(lexer.OR, LOGICAL, parseBinaryExpr)

	// Range
	led(lexer.DOT_DOT, LOGICAL, parseBinaryExpr)

	// Member
	led(lexer.DOT, MEMBER, parsePropertyExpr)

	// Relational
	led(lexer.LESS, RELATIONAL, parseBinaryExpr)
	led(lexer.LESS_EQUALS, RELATIONAL, parseBinaryExpr)
	led(lexer.GREATER, RELATIONAL, parseBinaryExpr)
	led(lexer.GREATER_EQUALS, RELATIONAL, parseBinaryExpr)
	led(lexer.EQUALS, RELATIONAL, parseBinaryExpr)
	led(lexer.NOT_EQUALS, RELATIONAL, parseBinaryExpr)

	// Additive & Multiplicative
	led(lexer.PLUS, ADDITIVE, parseBinaryExpr)
	led(lexer.MINUS, ADDITIVE, parseBinaryExpr)

	led(lexer.TIMES, MULTIPLICATIVE, parseBinaryExpr)
	led(lexer.DIVIDE, MULTIPLICATIVE, parseBinaryExpr)
	led(lexer.MODULO, MULTIPLICATIVE, parseBinaryExpr)

	//call
	led(lexer.OPEN_PAREN, CALL, parseCallExpr)

	// Statements
	stmt(lexer.CONST, parseVarDeclStmt)
	stmt(lexer.LET, parseVarDeclStmt)

	stmt(lexer.MODULE, parseModuleStmt)
	stmt(lexer.IMPORT, parseImportStmt)
	stmt(lexer.STRUCT, parseStructDeclStmt)
	stmt(lexer.TRAIT, parseTraitDeclStmt)
	stmt(lexer.IMPLEMENT, parseImplementStmt)
	stmt(lexer.OPEN_CURLY, parseBlockStmt)

	//conditionals
	stmt(lexer.IF, parseIfStatement)
	stmt(lexer.SWITCH, parseSwitchCaseStmt)
	//loops
	stmt(lexer.FOR, parseForLoopStmt)
	stmt(lexer.FOREACH, parseForLoopStmt)
	stmt(lexer.WHILE, parseWhileLoopStmt)

	//function
	stmt(lexer.FUNCTION, parseFunctionDeclStmt)
	stmt(lexer.RETURN, parseReturnStmt)

	stmt(lexer.CONTINUE, parseContinueStmt)
	stmt(lexer.BREAK, parseBreakStmt)
}
