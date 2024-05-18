package parser

import (
	"rexlang/frontend/ast"
	"rexlang/frontend/lexer"
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
type statementHandler func(p *Parser) ast.Stmt

// Null denoted. Expect nothing to the left of the token
type nudHandler func(p *Parser) ast.Expr

// Left denoted. Expect something to the left of the token
type LEDHandler func(p *Parser, left ast.Expr, bp BINDING_POWER) ast.Expr

// lookup table for the different token types
type stmtLookupType map[lexer.TOKEN_KIND]statementHandler
type nudLookupType map[lexer.TOKEN_KIND]nudHandler
type ledLookupType map[lexer.TOKEN_KIND]LEDHandler
type bpLookupType map[lexer.TOKEN_KIND]BINDING_POWER

var nudLookup = nudLookupType{}
var ledLookup = ledLookupType{}
var stmtLookup = stmtLookupType{}
var bpLookup = bpLookupType{}

func led(kind lexer.TOKEN_KIND, bp BINDING_POWER, led_fn LEDHandler) {
	bpLookup[kind] = bp
	ledLookup[kind] = led_fn
}

func nud(kind lexer.TOKEN_KIND, nud_fn nudHandler) {
	nudLookup[kind] = nud_fn
}

func stmt(kind lexer.TOKEN_KIND, stmt_fn statementHandler) {
	bpLookup[kind] = DEFAULT_BP
	stmtLookup[kind] = stmt_fn
}

func createTokenLookups() {

	// Assignment
	led(lexer.ASSIGNMENT, ASSIGNMENT, parse_var_assignment_expr)
	led(lexer.PLUS_EQUALS, ASSIGNMENT, parse_var_assignment_expr)
	led(lexer.MINUS_EQUALS, ASSIGNMENT, parse_var_assignment_expr)
	led(lexer.TIMES_EQUALS, ASSIGNMENT, parse_var_assignment_expr)
	led(lexer.DIVIDE_EQUALS, ASSIGNMENT, parse_var_assignment_expr)
	led(lexer.MODULO_EQUALS, ASSIGNMENT, parse_var_assignment_expr)

	// Logical operations
	led(lexer.AND, LOGICAL, parse_binary_expr)
	led(lexer.OR, LOGICAL, parse_binary_expr)
	led(lexer.DOT_DOT, LOGICAL, parse_binary_expr)

	// Relational
	led(lexer.LESS, RELATIONAL, parse_binary_expr)
	led(lexer.LESS_EQUALS, RELATIONAL, parse_binary_expr)
	led(lexer.GREATER, RELATIONAL, parse_binary_expr)
	led(lexer.GREATER_EQUALS, RELATIONAL, parse_binary_expr)
	led(lexer.EQUALS, RELATIONAL, parse_binary_expr)
	led(lexer.NOT_EQUALS, RELATIONAL, parse_binary_expr)

	// Additive & Multiplicative
	led(lexer.PLUS, ADDITIVE, parse_binary_expr)
	led(lexer.MINUS, ADDITIVE, parse_binary_expr)

	led(lexer.TIMES, MULTIPLICATIVE, parse_binary_expr)
	led(lexer.DIVIDE, MULTIPLICATIVE, parse_binary_expr)
	led(lexer.MODULO, MULTIPLICATIVE, parse_binary_expr)

	//call
	led(lexer.OPEN_PAREN, CALL, parse_call_expr)

	//literals & symbols
	nud(lexer.NUMBER, parse_primary_expr)
	nud(lexer.STRING, parse_primary_expr)
	nud(lexer.IDENTIFIER, parse_primary_expr)
	nud(lexer.TRUE, parse_primary_expr)
	nud(lexer.FALSE, parse_primary_expr)
	nud(lexer.NULL, parse_primary_expr)

	nud(lexer.OPEN_PAREN, parse_grouping_expr)

	//unary / prefix
	nud(lexer.MINUS, parse_prefix_expr)
	nud(lexer.PLUS, parse_prefix_expr)
	nud(lexer.PLUS_PLUS, parse_unary_expr)
	nud(lexer.MINUS_MINUS, parse_unary_expr)
	nud(lexer.NOT, parse_unary_expr)
	nud(lexer.OPEN_BRACKET, parse_array_expr)

	//call/member
	led(lexer.OPEN_CURLY, CALL, parse_struct_instantiation_expr)

	// Statements
	stmt(lexer.CONST, parse_var_decl_stmt)
	stmt(lexer.LET, parse_var_decl_stmt)
	stmt(lexer.IF, parse_if_statement)
	stmt(lexer.STRUCT, parse_struct_decl_stmt)
	//function
	stmt(lexer.FUNCTION, parse_function_decl_stmt)
	stmt(lexer.RETURN, parse_return_stmt)
	//stmt(lexer.ELSEIF, parse_if_statement)
}
