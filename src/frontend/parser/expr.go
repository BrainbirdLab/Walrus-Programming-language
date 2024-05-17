package parser

import (
	"fmt"
	"rexlang/frontend/ast"
	"rexlang/helpers"
	"rexlang/frontend/lexer"
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

func parse_call_expr(p *Parser, left ast.Expr, bp binding_power) ast.Expr {

	p.expect(lexer.OPEN_PAREN)

	var arguments []ast.Expr

	for p.currentTokenKind() != lexer.CLOSE_PAREN {
		//parse the arguments
		argument := parse_expr(p, default_bp)
		arguments = append(arguments, argument)
		
		if p.currentTokenKind() == lexer.COMMA {
			p.advance()
		}
	}

	p.expect(lexer.CLOSE_PAREN)

	return ast.FunctionCallExpr{
		Kind: ast.FUNCTION_CALL_EXPRESSION,
		Function: left,
		Args: arguments,
	}
}

func parse_expr(p *Parser, bp binding_power) ast.Expr {
	// Fist parse the NUD
	tokenKind := p.currentTokenKind()
	nud_fn, exists := nud_lu[tokenKind]

	if !exists {
		if lexer.IsKeyword(tokenKind) {
			panic(fmt.Sprintf("NUD handler expected for keyword '%s'\n", lexer.TokenKindString(tokenKind)))
		} else {
			panic(fmt.Sprintf("NUD handler expected for token '%s'\n", lexer.TokenKindString(tokenKind)))
		}
	}

	left := nud_fn(p)

	for bp_lu[p.currentTokenKind()] > bp {

		tokenKind := p.currentTokenKind()

		led_fn, exists := led_lu[tokenKind]

		if !exists {
			panic(fmt.Sprintf("LED handler expected for token %s\n", lexer.TokenKindString(tokenKind)))
		}

		left = led_fn(p, left, bp_lu[p.currentTokenKind()])
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
			Type:  "i8",
		}
	case lexer.STRING:
		return ast.StringLiteral{
			Kind:  ast.STRING_LITERAL,
			Value: p.advance().Value,
			Type:  "str",
		}
	case lexer.IDENTIFIER:
		return ast.SymbolExpr{
			Kind:   ast.IDENTIFIER,
			Symbol: p.advance().Value,
			Type:   "infr",
		}
	case lexer.TRUE:
		p.advance()
		return ast.BooleanLiteral{
			Kind:  ast.BOOLEAN_LITERAL,
			Value: true,
			Type:  "bool",
		}
	case lexer.FALSE:
		p.advance()
		return ast.BooleanLiteral{
			Kind:  ast.BOOLEAN_LITERAL,
			Value: false,
			Type:  "bool",
		}
	case lexer.NULL:
		p.advance()
		return ast.NullLiteral{
			Kind:  ast.NULL_LITERAL,
			Value: "null",
			Type:  "null",
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
	if _, ok := left.(ast.SymbolExpr); !ok {
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

	identifier, ok := left.(ast.SymbolExpr)

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

func parse_struct_instantiation_expr(p *Parser, left ast.Expr, bp binding_power) ast.Expr {
	// Check if left is an Identifier
	structName := helpers.ExpectType[ast.SymbolExpr](left).Symbol

	var properties = map[string]ast.Expr{}
	var methods = map[string]ast.FunctionDeclStmt{}

	p.expect(lexer.OPEN_CURLY)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		var propName = p.expect(lexer.IDENTIFIER).Value
		p.expect(lexer.COLON)
		expr := parse_expr(p, logical)

		properties[propName] = expr

		if p.currentTokenKind() != lexer.CLOSE_CURLY {
			p.expect(lexer.COMMA)
		}
	}

	p.expect(lexer.CLOSE_CURLY)

	return ast.StructInstantiationExpr{
		StructName: structName,
		Properties: properties,
		Methods: methods,
	}
}

func parse_array_expr(p *Parser) ast.Expr {

	p.expect(lexer.OPEN_BRACKET)

	elements := []ast.Expr{}

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_BRACKET {
		elements = append(elements, parse_expr(p, primary))
		if p.currentTokenKind() != lexer.CLOSE_BRACKET {
			p.expect(lexer.COMMA)
		}
	}

	p.expect(lexer.CLOSE_BRACKET)

	return ast.ArrayLiterals{
		Kind: ast.ARRAY_LITERALS,
		Elements: elements,
		Size: uint64(len(elements)),
	}
}