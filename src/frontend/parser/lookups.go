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

func led(kind lexer.TOKEN_KIND, bp BINDING_POWER, handleLED LEDHandler) {
	bpLookupMap[kind] = bp
	ledLookup[kind] = handleLED
}

func nud(kind lexer.TOKEN_KIND, handleNullDenotation NUDHandler) {
	nudLookup[kind] = handleNullDenotation
}

func stmt(kind lexer.TOKEN_KIND, handleStatement statementHandler) {
	stmtLookup[kind] = handleStatement
}

func createTokenLookups() {

	//literals & Identifiers
	nud(lexer.INTEGER_TOKEN, parsePrimaryExpr)
	nud(lexer.FLOATING_TOKEN, parsePrimaryExpr)
	nud(lexer.STRING_TOKEN, parsePrimaryExpr)
	nud(lexer.CHARACTER_TOKEN, parsePrimaryExpr)
	nud(lexer.IDENTIFIER_TOKEN, parsePrimaryExpr)
	nud(lexer.TRUE_TOKEN, parsePrimaryExpr)
	nud(lexer.FALSE_TOKEN, parsePrimaryExpr)
	nud(lexer.NULL_TOKEN, parsePrimaryExpr)

	nud(lexer.OPEN_PAREN_TOKEN, parseGroupingExpr)

	//unary / prefix
	nud(lexer.MINUS_TOKEN, parsePrefixExpr)
	nud(lexer.PLUS_TOKEN, parsePrefixExpr)
	nud(lexer.PLUS_PLUS_TOKEN, parseUnaryExpr)
	nud(lexer.MINUS_MINUS_TOKEN, parseUnaryExpr)
	nud(lexer.NOT_TOKEN, parseUnaryExpr)
	nud(lexer.OPEN_BRACKET_TOKEN, parseArrayExpr)

	// Assignment
	led(lexer.ASSIGNMENT_TOKEN, ASSIGNMENT, parseVarAssignmentExpr)
	led(lexer.PLUS_EQUALS_TOKEN, ASSIGNMENT, parseVarAssignmentExpr)
	led(lexer.MINUS_EQUALS_TOKEN, ASSIGNMENT, parseVarAssignmentExpr)
	led(lexer.TIMES_EQUALS_TOKEN, ASSIGNMENT, parseVarAssignmentExpr)
	led(lexer.DIVIDE_EQUALS_TOKEN, ASSIGNMENT, parseVarAssignmentExpr)
	led(lexer.MODULO_EQUALS_TOKEN, ASSIGNMENT, parseVarAssignmentExpr)

	// Logical operations
	led(lexer.AND_TOKEN, LOGICAL, parseBinaryExpr)
	led(lexer.OR_TOKEN, LOGICAL, parseBinaryExpr)

	// Range
	led(lexer.DOT_DOT_TOKEN, LOGICAL, parseBinaryExpr)

	// Member
	led(lexer.DOT_TOKEN, MEMBER, parsePropertyExpr)

	// Relational
	led(lexer.LESS_TOKEN, RELATIONAL, parseBinaryExpr)
	led(lexer.LESS_EQUALS_TOKEN, RELATIONAL, parseBinaryExpr)
	led(lexer.GREATER_TOKEN, RELATIONAL, parseBinaryExpr)
	led(lexer.GREATER_EQUALS_TOKEN, RELATIONAL, parseBinaryExpr)
	led(lexer.EQUALS_TOKEN, RELATIONAL, parseBinaryExpr)
	led(lexer.NOT_EQUALS_TOKEN, RELATIONAL, parseBinaryExpr)

	// Additive & Multiplicative
	led(lexer.PLUS_TOKEN, ADDITIVE, parseBinaryExpr)
	led(lexer.MINUS_TOKEN, ADDITIVE, parseBinaryExpr)

	led(lexer.TIMES_TOKEN, MULTIPLICATIVE, parseBinaryExpr)
	led(lexer.DIVIDE_TOKEN, MULTIPLICATIVE, parseBinaryExpr)
	led(lexer.MODULO_TOKEN, MULTIPLICATIVE, parseBinaryExpr)
	led(lexer.POWER_TOKEN, MULTIPLICATIVE, parseBinaryExpr)

	//call
	led(lexer.OPEN_PAREN_TOKEN, CALL, parseCallExpr)

	// Statements
	stmt(lexer.CONST_TOKEN, parseVarDeclStmt)
	stmt(lexer.LET_TOKEN, parseVarDeclStmt)

	stmt(lexer.MODULE_TOKEN, parseModuleStmt)
	stmt(lexer.IMPORT_TOKEN, parseImportStmt)
	stmt(lexer.STRUCT_TOKEN, parseStructDeclStmt)
	stmt(lexer.TRAIT_TOKEN, parseTraitDeclStmt)
	stmt(lexer.IMPLEMENT_TOKEN, parseImplementStmt)
	stmt(lexer.OPEN_CURLY_TOKEN, parseBlockStmt)

	//conditionals
	stmt(lexer.IF_TOKEN, parseIfStatement)
	stmt(lexer.SWITCH_TOKEN, parseSwitchCaseStmt)
	//loops
	stmt(lexer.FOR_TOKEN, parseForLoopStmt)
	stmt(lexer.FOREACH_TOKEN, parseForLoopStmt)
	stmt(lexer.WHILE_TOKEN, parseWhileLoopStmt)

	//function
	stmt(lexer.FUNCTION_TOKEN, parseFunctionDeclStmt)
	stmt(lexer.RETURN_TOKEN, parseReturnStmt)

	stmt(lexer.CONTINUE_TOKEN, parseContinueStmt)
	stmt(lexer.BREAK_TOKEN, parseBreakStmt)
}
