package parser

import (
	"rexlang/ast"
	"rexlang/lexer"
)

type binding_power int

//Operator precedence levels. primary is the highest binding power
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

func nud(kind lexer.TokenKind, bp binding_power, nud_fn nud_handler) {
	bp_lu[kind] = primary
	nud_lu[kind] = nud_fn
}

func stmt(kind lexer.TokenKind, stmt_fn stmt_handler) {
	bp_lu[kind] = default_bp
	stmt_lu[kind] = stmt_fn
}

func createTokenLookups() {

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


	//literals & symbols
	nud(lexer.NUMBER, primary, parse_primary_expr)
	nud(lexer.STRING, primary, parse_primary_expr)
	nud(lexer.IDENTIFIER, primary, parse_primary_expr)
	
	nud(lexer.OPEN_PAREN, primary, parse_primary_expr)
	nud(lexer.MINUS, unary, parse_unary_expr)


	// Statements
	stmt((lexer.CONST), parse_var_decl_stmt)
	stmt((lexer.LET), parse_var_decl_stmt)
}
