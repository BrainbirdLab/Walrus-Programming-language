package parser

import (
	"fmt"
	"rexlang/ast"
	"rexlang/lexer"
	"strconv"
)

func parse_expr(p *Parser, bp binding_power) ast.Expr {
	// Fist parse the NUD
	tokenKind := p.currentTokenKind()

	nud_fn, exists := nud_lu[tokenKind]

	if !exists {
		panic(fmt.Sprintf("NUD handler expected for token %s\n", lexer.TokenKindString(tokenKind)))
	}

	left := nud_fn(p)

	// While we have a led and the current bp is < bp of the current token
	//continue parsing the left hand side of the expression

	for bp_lu[p.currentTokenKind()] > bp {
		tokenKind = p.currentTokenKind()
		led_fn, exists := led_lu[tokenKind]

		if !exists {
			panic(fmt.Sprintf("LED handler expected for token %s\n", lexer.TokenKindString(tokenKind)))
		}

		left = led_fn(p, left, bp)
	}

	return left

}

func parse_primary_expr(p *Parser) ast.Expr {
	switch p.currentTokenKind() {
	case lexer.NUMBER:
		number, _ := strconv.ParseFloat(p.advance().Value, 64)
		return ast.NumericLiteral{
			Value: number,
		}
	case lexer.STRING:
		return ast.StringLiteral{
			Value: p.advance().Value,
		}
	case lexer.IDENTIFIER:
		return ast.SymbolExpr{
			Value: p.advance().Value,
		}
	case lexer.OPEN_PAREN:
		p.advance() // Consume the opening parenthesis
		expr := parse_expr(p, default_bp) // Parse expression inside parentheses
		p.expect(lexer.CLOSE_PAREN) // Expect closing parenthesis
		return expr
	default:
		panic(fmt.Sprintf("Cannot create primary expression from %s\n", lexer.TokenKindString(p.currentTokenKind())))
	}
}

func parse_binary_expr(p *Parser, left ast.Expr, bp binding_power) ast.Expr {
	operatorToken := p.advance()
	right := parse_expr(p, bp)

	return ast.BinaryExpr{
		Left:     left,
		Operator: operatorToken,
		Right:    right,
	}
}
