package parser

import (
	"rexlang/frontend/ast"
	"rexlang/frontend/lexer"
)

type binding_power int

// Operator precedence levels. primary is the highest binding power
const (
	default_bp binding_power = iota
	comma
	assignment
	logical
	relational
	additive
	multiplicative
	unary
	call
	member
	primary
)

// Statement handler. Expect nothing to the left of the token
type stmt_handler func(p *Parser) ast.Stmt

// Null denoted. Expect nothing to the left of the token
type nud_handler func(p *Parser) ast.Expr

// Left denoted. Expect something to the left of the token
type led_handler func(p *Parser, left ast.Expr, bp binding_power) ast.Expr

// lookup table for the different token types
type stmt_lookup map[lexer.TokenKind]stmt_handler
type nud_lookup map[lexer.TokenKind]nud_handler
type led_lookup map[lexer.TokenKind]led_handler
type bp_lookup map[lexer.TokenKind]binding_power


var bp_lu = bp_lookup{}
var nud_lu = nud_lookup{}
var led_lu = led_lookup{}
var stmt_lu = stmt_lookup{}


func led(kind lexer.TokenKind, bp binding_power, led_fn led_handler) {
	bp_lu[kind] = bp
	led_lu[kind] = led_fn
}

func nud(kind lexer.TokenKind, nud_fn nud_handler) {
	nud_lu[kind] = nud_fn
}

func stmt(kind lexer.TokenKind, stmt_fn stmt_handler) {
	bp_lu[kind] = default_bp
	stmt_lu[kind] = stmt_fn
}

func createTokenLookups() {

	// Assignment
	led(lexer.ASSIGNMENT, assignment, parse_var_assignment_expr)
	led(lexer.PLUS_EQUALS, assignment, parse_var_assignment_expr)
	led(lexer.MINUS_EQUALS, assignment, parse_var_assignment_expr)
	led(lexer.TIMES_EQUALS, assignment, parse_var_assignment_expr)
	led(lexer.DIVIDE_EQUALS, assignment, parse_var_assignment_expr)
	led(lexer.MODULO_EQUALS, assignment, parse_var_assignment_expr)

	// Logical operations
	led(lexer.AND, logical, parse_binary_expr)
	led(lexer.OR, logical, parse_binary_expr)
	led(lexer.DOT_DOT, logical, parse_binary_expr)

	// Relational
	led(lexer.LESS, relational, parse_binary_expr)
	led(lexer.LESS_EQUALS, relational, parse_binary_expr)
	led(lexer.GREATER, relational, parse_binary_expr)
	led(lexer.GREATER_EQUALS, relational, parse_binary_expr)
	led(lexer.EQUALS, relational, parse_binary_expr)
	led(lexer.NOT_EQUALS, relational, parse_binary_expr)

	// Additive & Multiplicative
	led(lexer.PLUS, additive, parse_binary_expr)
	led(lexer.MINUS, additive, parse_binary_expr)

	led(lexer.TIMES, multiplicative, parse_binary_expr)
	led(lexer.DIVIDE, multiplicative, parse_binary_expr)
	led(lexer.MODULO, multiplicative, parse_binary_expr)

	//call
	led(lexer.OPEN_PAREN, call, parse_call_expr)

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
	led(lexer.OPEN_CURLY, call, parse_struct_instantiation_expr)

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
