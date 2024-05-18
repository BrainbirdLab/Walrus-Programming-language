package parser

import (
	"fmt"
	"rexlang/frontend/ast"
	"rexlang/frontend/lexer"
	"rexlang/helpers"
	"strconv"
)

// parse_binary_expr parses a binary expression, given the left-hand side expression
// and the current binding power. It advances the parser to the next token,
// parses the right-hand side expression, and returns an AST node representing
// the binary expression.
func parse_binary_expr(p *Parser, left ast.Expr, bp BINDING_POWER) ast.Expr {

	start := p.currentToken().StartPos

	operatorToken := p.advance()

	right := parse_expr(p, bp)

	_, end := right.GetPos()

	return ast.BinaryExpr{
		Kind:     ast.BINARY_EXPRESSION,
		Operator: operatorToken,
		Left:     left,
		Right:    right,
		StartPos: start,
		EndPos:   end,
	}
}

// parse_call_expr parses a function call expression, including the function name and its arguments.
// It expects the current token to be an opening parenthesis, and it will parse the arguments
// until it encounters a closing parenthesis. The function returns an ast.FunctionCallExpr
// representing the parsed function call.
func parse_call_expr(p *Parser, left ast.Expr, bp BINDING_POWER) ast.Expr {

	start := p.currentToken().StartPos

	p.expect(lexer.OPEN_PAREN)

	var arguments []ast.Expr

	for p.currentTokenKind() != lexer.CLOSE_PAREN {
		//parse the arguments
		argument := parse_expr(p, DEFAULT_BP)
		arguments = append(arguments, argument)

		if p.currentTokenKind() == lexer.COMMA {
			p.advance()
		}
	}

	end := p.expect(lexer.CLOSE_PAREN).EndPos

	return ast.FunctionCallExpr{
		Kind:     ast.FUNCTION_CALL_EXPRESSION,
		Function: left,
		Args:     arguments,
		StartPos: start,
		EndPos:   end,
	}
}

// parse_expr parses an expression with the given binding power.
// It first parses the NUD (Null Denotation) of the expression,
// then continues to parse the LED (Left Denotation) of the expression
// until the binding power of the current token is less than or equal to the given binding power.
// The parsed expression is returned as an ast.Expr.
func parse_expr(p *Parser, bp BINDING_POWER) ast.Expr {

	// Fist parse the NUD
	tokenKind := p.currentTokenKind()
	nud_fn, exists := nudLookup[tokenKind]

	if !exists {
		if lexer.IsKeyword(tokenKind) {
			panic(fmt.Sprintf("NUD handler expected for keyword '%s'\n", tokenKind))
		} else {
			panic(fmt.Sprintf("NUD handler expected for token '%s'\n", tokenKind))
		}
	}

	left := nud_fn(p)

	for bpLookup[p.currentTokenKind()] > bp {

		tokenKind := p.currentTokenKind()

		led_fn, exists := ledLookup[tokenKind]

		if !exists {
			panic(fmt.Sprintf("LED handler expected for token %s\n", tokenKind))
		}

		left = led_fn(p, left, bpLookup[p.currentTokenKind()])
	}

	return left
}

// parse_primary_expr parses a primary expression in the input stream.
// It handles numeric literals, string literals, identifiers, boolean literals, and null literals.
// If the current token does not match any of these types, it panics with an error message.
func parse_primary_expr(p *Parser) ast.Expr {

	startpos := p.currentToken().StartPos

	switch p.currentTokenKind() {
	case lexer.NUMBER:
		number, _ := strconv.ParseFloat(p.advance().Value, 64)
		return ast.NumericLiteral{
			Kind:     ast.NUMERIC_LITERAL,
			Value:    number,
			Type:     "i8",
			StartPos: startpos,
			EndPos:   p.currentToken().EndPos,
		}
	case lexer.STRING:
		return ast.StringLiteral{
			Kind:     ast.STRING_LITERAL,
			Value:    p.advance().Value,
			Type:     "str",
			StartPos: startpos,
			EndPos:   p.currentToken().EndPos,
		}
	case lexer.IDENTIFIER:
		return ast.SymbolExpr{
			Kind:     ast.IDENTIFIER,
			Symbol:   p.advance().Value,
			Type:     "infr",
			StartPos: startpos,
			EndPos:   p.currentToken().EndPos,
		}
	case lexer.TRUE:
		p.advance()
		return ast.BooleanLiteral{
			Kind:     ast.BOOLEAN_LITERAL,
			Value:    true,
			Type:     "bool",
			StartPos: startpos,
			EndPos:   p.currentToken().EndPos,
		}
	case lexer.FALSE:
		p.advance()
		return ast.BooleanLiteral{
			Kind:     ast.BOOLEAN_LITERAL,
			Value:    false,
			Type:     "bool",
			StartPos: startpos,
			EndPos:   p.currentToken().EndPos,
		}
	case lexer.NULL:
		p.advance()
		return ast.NullLiteral{
			Kind:     ast.NULL_LITERAL,
			Value:    "null",
			Type:     "null",
			StartPos: startpos,
			EndPos:   p.currentToken().EndPos,
		}

	default:
		panic(fmt.Sprintf("Cannot create primary expression from %s\n", p.currentTokenKind()))
	}
}

// parse_grouping_expr parses a grouping expression, which is an expression
// enclosed in parentheses. It expects the opening parenthesis, parses the
// expression inside, and then expects the closing parenthesis.
func parse_grouping_expr(p *Parser) ast.Expr {
	p.expect(lexer.OPEN_PAREN)
	expression := parse_expr(p, DEFAULT_BP)
	p.expect(lexer.CLOSE_PAREN)
	return expression
}

// parse_prefix_expr parses a prefix expression, which consists of a unary operator
// followed by an expression. It returns an ast.UnaryExpr representing the parsed
// prefix expression.
func parse_prefix_expr(p *Parser) ast.Expr {

	startpos := p.currentToken().StartPos

	operator := p.advance()

	expr := parse_expr(p, UNARY)

	_, end := expr.GetPos()

	return ast.UnaryExpr{
		Kind:     ast.UNARY_EXPRESSION,
		Operator: operator,
		Argument: expr,
		StartPos: startpos,
		EndPos:   end,
	}
}

// parse_postfix_expr parses a postfix expression, which can be an increment or
// decrement operation on an identifier. It checks that the left-hand side is a
// valid identifier, and returns an AST node representing the unary expression.
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

// parse_unary_expr parses a unary expression from the input stream.
// It returns the parsed expression as an ast.Expr.
func parse_unary_expr(p *Parser) ast.Expr {
	return parse_prefix_expr(p)
}

// parse_var_assignment_expr parses a variable assignment expression. It takes a Parser, a left-hand side expression, and a binding power.
// If the left-hand side is an identifier, it creates an AssignmentExpr with the identifier, the assignment operator, and the right-hand side expression.
// If the left-hand side is not an identifier, it panics with an error message.
func parse_var_assignment_expr(p *Parser, left ast.Expr, bp BINDING_POWER) ast.Expr {
	// Check if left is an Identifier

	start := p.currentToken().StartPos

	identifier, ok := left.(ast.SymbolExpr)

	if !ok {
		panic("Cannot assign value: Expected an identifier on the left side of the assignment")
	}

	operator := p.advance()

	right := parse_expr(p, bp)

	_, end := right.GetPos()

	return ast.AssignmentExpr{
		Kind:     ast.ASSIGNMENT_EXPRESSION,
		Assigne:  identifier,
		Operator: operator,
		Value:    right,
		StartPos: start,
		EndPos:   end,
	}
}

// parse_struct_instantiation_expr parses a struct instantiation expression, which creates a new instance of a struct.
// It expects the left-hand side to be an identifier representing the struct type, followed by a block of property assignments
// enclosed in curly braces. The function returns an ast.StructInstantiationExpr representing the parsed expression.
func parse_struct_instantiation_expr(p *Parser, left ast.Expr, bp BINDING_POWER) ast.Expr {

	start := p.currentToken().StartPos

	// Check if left is an Identifier
	structName := helpers.ExpectType[ast.SymbolExpr](left).Symbol

	var properties = map[string]ast.Expr{}
	var methods = map[string]ast.FunctionDeclStmt{}

	p.expect(lexer.OPEN_CURLY)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		var propName = p.expect(lexer.IDENTIFIER).Value
		p.expect(lexer.COLON)
		expr := parse_expr(p, LOGICAL)

		properties[propName] = expr

		if p.currentTokenKind() != lexer.CLOSE_CURLY {
			p.expect(lexer.COMMA)
		}
	}

	end := p.expect(lexer.CLOSE_CURLY).EndPos

	return ast.StructInstantiationExpr{
		StructName: structName,
		Properties: properties,
		Methods:    methods,
		StartPos:   start,
		EndPos:     end,
	}
}

// parse_array_expr parses an array expression in the input stream.
// It expects the opening '[' bracket, parses the array elements,
// and returns an ast.ArrayLiterals node representing the array.
func parse_array_expr(p *Parser) ast.Expr {

	start := p.currentToken().StartPos

	p.expect(lexer.OPEN_BRACKET)

	elements := []ast.Expr{}

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_BRACKET {
		elements = append(elements, parse_expr(p, PRIMARY))
		if p.currentTokenKind() != lexer.CLOSE_BRACKET {
			p.expect(lexer.COMMA)
		}
	}

	end := p.expect(lexer.CLOSE_BRACKET).EndPos

	return ast.ArrayLiterals{
		Kind:     ast.ARRAY_LITERALS,
		Elements: elements,
		Size:     uint64(len(elements)),
		StartPos: start,
		EndPos:   end,
	}
}
