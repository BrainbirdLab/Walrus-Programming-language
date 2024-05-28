package parser

import (
	"fmt"
	"strconv"
	"walrus/frontend/ast"
	"walrus/frontend/lexer"
	"walrus/helpers"
)

// parse_binary_expr parses a binary expression, given the left-hand side expression
// and the current binding power. It advances the parser to the next token,
// parses the right-hand side expression, and returns an AST node representing
// the binary expression.
func parse_binary_expr(p *Parser, left ast.Expression, bp BINDING_POWER) ast.Expression {

	start := p.currentToken().StartPos

	operatorToken := p.advance()

	right := parse_expr(p, bp)

	_, end := right.GetPos()

	return ast.BinaryExpr{
		BaseStmt: ast.BaseStmt{
			Kind:    ast.BINARY_EXPRESSION,
			StartPos: start,
			EndPos:   end,
		},
		Operator: operatorToken,
		Left:     left,
		Right:    right,
	}
}

// parse_call_expr parses a function call expression, including the function name and its arguments.
// It expects the current token to be an opening parenthesis, and it will parse the arguments
// until it encounters a closing parenthesis. The function returns an ast.FunctionCallExpr
// representing the parsed function call.
func parse_call_expr(p *Parser, left ast.Expression, bp BINDING_POWER) ast.Expression {

	start := p.currentToken().StartPos

	p.expect(lexer.OPEN_PAREN)

	var arguments []ast.Expression

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
		BaseStmt: ast.BaseStmt{
			Kind:    ast.FUNCTION_CALL_EXPRESSION,
			StartPos: start,
			EndPos:   end,
		},
		Function: left,
		Args:     arguments,
	}
}

func parse_property_expr(p *Parser, left ast.Expression, bp BINDING_POWER) ast.Expression {

	p.expect(lexer.DOT)

	identifier := p.expect(lexer.IDENTIFIER)

	property := ast.IdentifierExpr{
		BaseStmt: ast.BaseStmt{
			Kind:    ast.IDENTIFIER,
			StartPos: identifier.StartPos,
			EndPos:   identifier.EndPos,
		},
		Identifier: identifier.Value,
		Type:       "infr",
	}

	start := property.StartPos

	return ast.StructPropertyExpr{
		BaseStmt: ast.BaseStmt{
			Kind:    ast.STRUCT_PROPERTY,
			StartPos: start,
			EndPos:   property.EndPos,
		},
		Object:   left,
		Property: property,
	}
}

// parse_expr parses an expression with the given binding power.
// It first parses the NUD (Null Denotation) of the expression,
// then continues to parse the LED (Left Denotation) of the expression
// until the binding power of the current token is less than or equal to the given binding power.
// The parsed expression is returned as an ast.Expression.
func parse_expr(p *Parser, bp BINDING_POWER) ast.Expression {

	// Fist parse the NUD
	token := p.currentToken()

	tokenKind := token.Kind

	if tokenKind == lexer.IDENTIFIER && p.nextToken().Kind == lexer.OPEN_CURLY && (p.previousToken().Kind == lexer.WALRUS || p.previousToken().Kind == lexer.ASSIGNMENT) {
		// Function call
		return parse_struct_instantiation_expr(p, parse_primary_expr(p))
	}

	nud_fn, exists := nudLookup[tokenKind]

	if !exists {

		var msg string
		if lexer.IsKeyword(tokenKind) {
			msg = fmt.Sprintf("Parser:NUD:Unexpected keyword '%s'\n", tokenKind)
		} else {
			msg = fmt.Sprintf("Parser:NUD:Unexpected token '%s'\n", tokenKind)
		}
		//err := fmt.Sprintf("File: %s:%d:%d: %s\n", p.FilePath, token.StartPos.Line, token.StartPos.Column, msg)

		p.MakeError(token.StartPos.Line, p.FilePath, token, msg).Display()
	}


	left := nud_fn(p)

	for GetBP(p.currentTokenKind()) > bp {

		tokenKind = p.currentTokenKind()

		led_fn, exists := ledLookup[tokenKind]

		if !exists {
			msg := fmt.Sprintf("Parser:LED:Unexpected token %s\n", tokenKind)
			p.MakeError(p.currentToken().StartPos.Line, p.FilePath, p.currentToken(), msg).Display()
		}

		left = led_fn(p, left, GetBP(p.currentTokenKind()))
	}

	return left
}

// parse_primary_expr parses a primary expression in the input stream.
// It handles numeric literals, string literals, identifiers, boolean literals, and null literals.
// If the current token does not match any of these types, it panics with an error message.
func parse_primary_expr(p *Parser) ast.Expression {

	startpos := p.currentToken().StartPos

	endpos := p.currentToken().EndPos

	switch p.currentTokenKind() {
	case lexer.NUMBER:
		number, _ := strconv.ParseFloat(p.advance().Value, 64)
		return ast.NumericLiteral{
			BaseStmt: ast.BaseStmt{
				Kind:    ast.NUMERIC_LITERAL,
				StartPos: startpos,
				EndPos:   endpos,
			},
			Value:    number,
			Type:     "i8",
		}
	case lexer.STRING:
		return ast.StringLiteral{
			BaseStmt: ast.BaseStmt{
				Kind:    ast.STRING_LITERAL,
				StartPos: startpos,
				EndPos:   endpos,
			},
			Value:    p.advance().Value,
			Type:     "str",
		}
	case lexer.IDENTIFIER:
		return ast.IdentifierExpr{
			BaseStmt: ast.BaseStmt{
				Kind:    ast.IDENTIFIER,
				StartPos: startpos,
				EndPos:   endpos,
			},
			Identifier: p.advance().Value,
			Type:       "infr",
		}
	case lexer.TRUE:
		p.advance()
		return ast.BooleanLiteral{
			BaseStmt: ast.BaseStmt{
				Kind:    ast.BOOLEAN_LITERAL,
				StartPos: startpos,
				EndPos:   endpos,
			},
			Value:    true,
			Type:     "bool",
		}
	case lexer.FALSE:
		p.advance()
		return ast.BooleanLiteral{
			BaseStmt: ast.BaseStmt{
				Kind:    ast.BOOLEAN_LITERAL,
				StartPos: startpos,
				EndPos:   endpos,
			},
			Value:    false,
			Type:     "bool",
		}
	case lexer.NULL:
		p.advance()
		return ast.NullLiteral{
			BaseStmt: ast.BaseStmt{
				Kind:    ast.NULL_LITERAL,
				StartPos: startpos,
				EndPos:   endpos,
			},
			Value:    "null",
			Type:     "null",
		}

	default:
		panic(fmt.Sprintf("Cannot create primary expression from %s\n", p.currentTokenKind()))
	}
}

// parse_grouping_expr parses a grouping expression, which is an expression
// enclosed in parentheses. It expects the opening parenthesis, parses the
// expression inside, and then expects the closing parenthesis.
func parse_grouping_expr(p *Parser) ast.Expression {
	p.expect(lexer.OPEN_PAREN)
	expression := parse_expr(p, DEFAULT_BP)
	p.expect(lexer.CLOSE_PAREN)
	return expression
}

// parse_prefix_expr parses a prefix expression, which consists of a unary operator
// followed by an expression. It returns an ast.UnaryExpr representing the parsed
// prefix expression.
func parse_prefix_expr(p *Parser) ast.Expression {

	startpos := p.currentToken().StartPos

	operator := p.advance()

	expr := parse_expr(p, UNARY)

	_, endpos := expr.GetPos()

	return ast.UnaryExpr{
		BaseStmt: ast.BaseStmt{
			Kind:    ast.UNARY_EXPRESSION,
			StartPos: startpos,
			EndPos:   endpos,
		},
		Operator: operator,
		Argument: expr,
	}
}


// parse_unary_expr parses a unary expression from the input stream.
// It returns the parsed expression as an ast.Expression.
func parse_unary_expr(p *Parser) ast.Expression {
	return parse_prefix_expr(p)
}

// parse_var_assignment_expr parses a variable assignment expression. It takes a Parser, a left-hand side expression, and a binding power.
// If the left-hand side is an identifier, it creates an AssignmentExpr with the identifier, the assignment operator, and the right-hand side expression.
// If the left-hand side is not an identifier, it panics with an error message.
func parse_var_assignment_expr(p *Parser, left ast.Expression, bp BINDING_POWER) ast.Expression {
	// Check if left is an Identifier

	start := p.currentToken().StartPos

	var identifier ast.IdentifierExpr

	switch assignee := left.(type) {

	case ast.IdentifierExpr:

		identifier = assignee

	case ast.StructPropertyExpr:

		identifier = assignee.Property

	default:
		errMsg := "Cannot assign to a non-identifier\n"
		p.MakeError(start.Line, p.FilePath, p.previousToken(), errMsg).AddHint("Expected an identifier", TEXT_HINT).Display()
	}

	operator := p.advance()

	right := parse_expr(p, bp)

	_, end := right.GetPos()

	return ast.AssignmentExpr{
		BaseStmt: ast.BaseStmt{
			Kind:    ast.ASSIGNMENT_EXPRESSION,
			StartPos: start,
			EndPos:   end,
		},
		Assigne:  identifier,
		Operator: operator,
		Value:    right,
	}
}

// parse_struct_instantiation_expr parses a struct instantiation expression, which creates a new instance of a struct.
// It expects the left-hand side to be an identifier representing the struct type, followed by a block of property assignments
// enclosed in curly braces. The function returns an ast.StructInstantiationExpr representing the parsed expression.
func parse_struct_instantiation_expr(p *Parser, left ast.Expression) ast.Expression {

	start := p.currentToken().StartPos

	// Check if left is an Identifier
	structName := helpers.ExpectType[ast.IdentifierExpr](left).Identifier

	var properties = map[string]ast.Expression{}
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
		BaseStmt: ast.BaseStmt{
			Kind:    ast.STRUCT_LITERAL,
			StartPos: start,
			EndPos:   end,
		},
		StructName: structName,
		Properties: properties,
		Methods:    methods,
	}
}

// parse_array_expr parses an array expression in the input stream.
// It expects the opening '[' bracket, parses the array elements,
// and returns an ast.ArrayLiterals node representing the array.
func parse_array_expr(p *Parser) ast.Expression {

	start := p.currentToken().StartPos

	p.expect(lexer.OPEN_BRACKET)

	elements := []ast.Expression{}

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_BRACKET {
		elements = append(elements, parse_expr(p, PRIMARY))
		if p.currentTokenKind() != lexer.CLOSE_BRACKET {
			p.expect(lexer.COMMA)
		}
	}

	end := p.expect(lexer.CLOSE_BRACKET).EndPos

	return ast.ArrayLiterals{
		BaseStmt: ast.BaseStmt{
			Kind:    ast.ARRAY_LITERALS,
			StartPos: start,
			EndPos:   end,
		},
		Elements: elements,
		Size:     uint64(len(elements)),
	}
}
