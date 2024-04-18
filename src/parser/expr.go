package parser

import (
	"fmt"
	"rexlang/ast"
	"rexlang/lexer"
	"strconv"
)

func parse_binary_expr(p *Parser, left ast.Expr, bp binding_power) ast.Expr {

	operatorToken := p.advance()

	right := parse_expr(p, bp)

	return ast.BinaryExpr{
		Kind:     ast.BINARY_EXPRESSION,
		Operator: operatorToken,
		Left:     left,
		Right:    right,
	}
}

func parse_expr(p *Parser, bp binding_power) ast.Expr {
	// Fist parse the NUD
	tokenKind := p.currentTokenKind()
	fmt.Printf("Parsing expression with token: %s\n", lexer.TokenKindString(tokenKind))

	nud_fn, exists := nud_lu[tokenKind]

	if !exists {
		fmt.Printf("Current token: %s on pos: %d\n", lexer.TokenKindString(tokenKind), p.pos)
		panic(fmt.Sprintf("NUD handler expected for token %s\n", lexer.TokenKindString(tokenKind)))
	}

	left := nud_fn(p)

	// While we have a led and the current bp is < bp of the current token
	//continue parsing the left hand side of the expression

	for bp_lu[p.currentTokenKind()] > bp {

		tokenKind := p.currentTokenKind()

		led_fn, exists := led_lu[tokenKind]

		if !exists {
			panic(fmt.Sprintf("LED handler expected for token %s\n", lexer.TokenKindString(tokenKind)))
		}

		left = led_fn(p, left, bp_lu[tokenKind])
	}

	return left

}

func parse_primary_expr(p *Parser) ast.Expr {
	switch p.currentTokenKind() {
	case lexer.NUMBER:
		number, _ := strconv.ParseFloat(p.advance().Value, 64)
		return ast.NumericLiteral{
			Kind:  ast.NUMERIC_LITERAL,
			Value: number,
			Type: "i8",
		}
	case lexer.STRING:
		return ast.StringLiteral{
			Kind:  ast.STRING_LITERAL,
			Value: p.advance().Value,
			Type: "str",
		}
	case lexer.IDENTIFIER:
		return ast.Identifier{
			Kind:   ast.IDENTIFIER,
			Symbol: p.advance().Value,
			Type:   "infr",
		}
	case lexer.TRUE:
		p.advance()
		return ast.BooleanLiteral{
			Kind:  ast.BOOLEAN_LITERAL,
			Value: true,
			Type: "bool",
		}
	case lexer.FALSE:
		p.advance()
		return ast.BooleanLiteral{
			Kind:  ast.BOOLEAN_LITERAL,
			Value: false,
			Type: "bool",
		}
	case lexer.NULL:
		p.advance()
		return ast.NullLiteral{
			Kind:  ast.NULL_LITERAL,
			Value: "null",
			Type: "null",
		}

	default:
		panic(fmt.Sprintf("Cannot create primary expression from %s\n", lexer.TokenKindString(p.currentTokenKind())))
	}
}

func parse_grouping_expr(p *Parser) ast.Expr {
	p.expect(lexer.OPEN_PAREN)
	expression := parse_expr(p, default_bp)
	p.expect(lexer.CLOSE_PAREN)
	return expression
}

func parse_prefix_expr(p *Parser) ast.Expr {

	operator := p.advance()

	expr := parse_expr(p, unary)

	return ast.UnaryExpr{
		Kind:     ast.UNARY_EXPRESSION,
		Operator: operator,
		Argument: expr,
	}
}

func parse_postfix_expr(p *Parser, left ast.Expr) ast.Expr {
	
	// a++ 
	// a should be a lvalue
	// a LValue is something that can be assigned to

	// Check if left is an Identifier
	if _, ok := left.(ast.Identifier); !ok {
		panic("Cannot increment or decrement value: Expected an identifier")
	}

	operator := p.advance()

	return ast.UnaryExpr{
		Kind:     ast.UNARY_EXPRESSION,
		Operator: operator,
		Argument: left,
	}
}

func parse_unary_expr(p *Parser) ast.Expr {
	return parse_prefix_expr(p)
}

func parse_var_assignment_expr(p *Parser, left ast.Expr, bp binding_power) ast.Expr {
	// Check if left is an Identifier

	identifier, ok := left.(ast.Identifier)

	if !ok {
		panic("Cannot assign value: Expected an identifier on the left side of the assignment")
	}

	operator := p.advance()

	right := parse_expr(p, bp)

	return ast.AssignmentExpr{
		Kind:     ast.ASSIGNMENT_EXPRESSION,
		Assigne:  identifier,
		Operator: operator,
		Value:    right,
	}
}
