package parser

import (
	"fmt"
	"strconv"
	"walrus/frontend/ast"
	"walrus/frontend/lexer"
	"walrus/helpers"
)

// parseBinaryExpr parses a binary expression, given the left-hand side expression
// and the current binding power. It advances the parser to the next token,
// parses the right-hand side expression, and returns an AST iNode representing
// the binary expression.
func parseBinaryExpr(p *Parser, left ast.Expression, bp BINDING_POWER) ast.Expression {

	start := p.currentToken().StartPos

	operatorToken := p.advance()

	right := parseExpr(p, bp)

	_, end := right.GetPos()

	return ast.BinaryExpr{
		BaseStmt: ast.BaseStmt{
			Kind:     ast.BINARY_EXPRESSION,
			StartPos: start,
			EndPos:   end,
		},
		Operator: operatorToken,
		Left:     left,
		Right:    right,
	}
}

// parses a function call expression, including the function name and its arguments.
// It expects the current token to be an opening parenthesis, and it will parse the arguments
// until it encounters a closing parenthesis. The function returns an ast.FunctionCallExpr
// representing the parsed function call.
func parseCallExpr(p *Parser, left ast.Expression, bp BINDING_POWER) ast.Expression {

	start := p.currentToken().StartPos

	p.expect(lexer.OPEN_PAREN)

	var arguments []ast.Expression

	for p.currentTokenKind() != lexer.CLOSE_PAREN {
		//parse the arguments
		argument := parseExpr(p, DEFAULT_BP)
		arguments = append(arguments, argument)

		if p.currentTokenKind() == lexer.COMMA {
			p.advance()
		}
	}

	end := p.expect(lexer.CLOSE_PAREN).EndPos

	return ast.FunctionCallExpr{
		BaseStmt: ast.BaseStmt{
			Kind:     ast.FUNCTION_CALL_EXPRESSION,
			StartPos: start,
			EndPos:   end,
		},
		Function: left,
		Args:     arguments,
	}
}

func parsePropertyExpr(p *Parser, left ast.Expression, bp BINDING_POWER) ast.Expression {

	p.expect(lexer.DOT)

	identifier := p.expect(lexer.IDENTIFIER)

	property := ast.IdentifierExpr{
		BaseStmt: ast.BaseStmt{
			Kind:     ast.IDENTIFIER,
			StartPos: identifier.StartPos,
			EndPos:   identifier.EndPos,
		},
		Identifier: identifier.Value,
	}

	start := property.StartPos

	return ast.StructPropertyExpr{
		BaseStmt: ast.BaseStmt{
			Kind:     ast.STRUCT_PROPERTY,
			StartPos: start,
			EndPos:   property.EndPos,
		},
		Object:   left,
		Property: property,
	}
}

// parseExpr parses an expression with the given binding power.
// It first parses the NUD (Null Denotation) of the expression,
// then continues to parse the LED (Left Denotation) of the expression
// until the binding power of the current token is less than or equal to the given binding power.
// The parsed expression is returned as an ast.Expression.
func parseExpr(p *Parser, bp BINDING_POWER) ast.Expression {

	// Fist parse the NUD
	token := p.currentToken()

	tokenKind := token.Kind

	if tokenKind == lexer.IDENTIFIER && p.nextToken().Kind == lexer.OPEN_CURLY && (p.previousToken().Kind == lexer.WALRUS || p.previousToken().Kind == lexer.ASSIGNMENT) {
		// Function call
		return parseStructInstantiationExpr(p, parsePrimaryExpr(p))
	}

	nudFunction, exists := nudLookup[tokenKind]

	if !exists {

		var msg string
		if lexer.IsKeyword(tokenKind) {
			msg = fmt.Sprintf("Parser:NUD:Unexpected keyword '%s'\n", tokenKind)
		} else {
			msg = fmt.Sprintf("Parser:NUD:Unexpected token '%s'\n", tokenKind)
		}
		//err := fmt.Sprintf("File: %s:%d:%d: %s\n", p.FilePath, token.StartPos.Line, token.StartPos.Column, msg)

		MakeError(p, token.StartPos.Line, p.FilePath, token.StartPos, token.EndPos, msg).Display()
	}

	left := nudFunction(p)

	for GetBP(p.currentTokenKind()) > bp {

		tokenKind = p.currentTokenKind()

		ledFunction, exists := ledLookup[tokenKind]

		if !exists {
			msg := fmt.Sprintf("Parser:LED:Unexpected token %s\n", tokenKind)
			MakeError(p, p.currentToken().StartPos.Line, p.FilePath, p.currentToken().StartPos, p.currentToken().EndPos, msg).Display()
		}

		left = ledFunction(p, left, GetBP(p.currentTokenKind()))
	}

	return left
}

// parsePrimaryExpr parses a primary expression in the input stream.
// It handles numeric literals, string literals, identifiers, boolean literals, and null literals.
// If the current token does not match any of these types, it panics with an error message.
func parsePrimaryExpr(p *Parser) ast.Expression {

	startpos := p.currentToken().StartPos

	endpos := p.currentToken().EndPos

	switch p.currentTokenKind() {
	case lexer.NUMBER:
		number, _ := strconv.ParseFloat(p.advance().Value, 64)
		return ast.NumericLiteral{
			BaseStmt: ast.BaseStmt{
				Kind:     ast.NUMERIC_LITERAL,
				StartPos: startpos,
				EndPos:   endpos,
			},
			Value: number,
		}
	case lexer.STRING:
		return ast.StringLiteral{
			BaseStmt: ast.BaseStmt{
				Kind:     ast.STRING_LITERAL,
				StartPos: startpos,
				EndPos:   endpos,
			},
			Value: p.advance().Value,
		}
	case lexer.CHARACTER:
		return ast.CharacterLiteral{
			BaseStmt: ast.BaseStmt{
				Kind:     ast.CHARACTER_LITERAL,
				StartPos: startpos,
				EndPos:   endpos,
			},
			Value: p.advance().Value,
		}
	case lexer.IDENTIFIER:
		return ast.IdentifierExpr{
			BaseStmt: ast.BaseStmt{
				Kind:     ast.IDENTIFIER,
				StartPos: startpos,
				EndPos:   endpos,
			},
			Identifier: p.advance().Value,
		}
	case lexer.TRUE:
		p.advance()
		return ast.BooleanLiteral{
			BaseStmt: ast.BaseStmt{
				Kind:     ast.BOOLEAN_LITERAL,
				StartPos: startpos,
				EndPos:   endpos,
			},
			Value: true,
		}
	case lexer.FALSE:
		p.advance()
		return ast.BooleanLiteral{
			BaseStmt: ast.BaseStmt{
				Kind:     ast.BOOLEAN_LITERAL,
				StartPos: startpos,
				EndPos:   endpos,
			},
			Value: false,
		}
	case lexer.NULL:
		p.advance()
		return ast.NullLiteral{
			BaseStmt: ast.BaseStmt{
				Kind:     ast.NULL_LITERAL,
				StartPos: startpos,
				EndPos:   endpos,
			},
			Value: "null",
		}

	default:
		panic(fmt.Sprintf("Cannot create primary expression from %s\n", p.currentTokenKind()))
	}
}

// parseGroupingExpr parses a grouping expression, which is an expression
// enclosed in parentheses. It expects the opening parenthesis, parses the
// expression inside, and then expects the closing parenthesis.
func parseGroupingExpr(p *Parser) ast.Expression {
	p.expect(lexer.OPEN_PAREN)
	expression := parseExpr(p, DEFAULT_BP)
	p.expect(lexer.CLOSE_PAREN)
	return expression
}

// parsePrefixExpr parses a prefix expression, which consists of a unary operator
// followed by an expression. It returns an ast.UnaryExpr representing the parsed
// prefix expression.
func parsePrefixExpr(p *Parser) ast.Expression {

	startpos := p.currentToken().StartPos

	operator := p.advance()

	expr := parseExpr(p, UNARY)

	_, endpos := expr.GetPos()

	return ast.UnaryExpr{
		BaseStmt: ast.BaseStmt{
			Kind:     ast.UNARY_EXPRESSION,
			StartPos: startpos,
			EndPos:   endpos,
		},
		Operator: operator,
		Argument: expr,
	}
}

// parseUnaryExpr parses a unary expression from the input stream.
// It returns the parsed expression as an ast.Expression.
func parseUnaryExpr(p *Parser) ast.Expression {
	return parsePrefixExpr(p)
}

// parseVarAssignmentExpr parses a variable assignment expression. It takes a Parser, a left-hand side expression, and a binding power.
// If the left-hand side is an identifier, it creates an AssignmentExpr with the identifier, the assignment operator, and the right-hand side expression.
// If the left-hand side is not an identifier, it panics with an error message.
func parseVarAssignmentExpr(p *Parser, left ast.Expression, bp BINDING_POWER) ast.Expression {
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
		MakeError(p, start.Line, p.FilePath, p.previousToken().StartPos, p.previousToken().EndPos, errMsg).AddHint("Expected an identifier", TEXT_HINT).Display()
	}

	operator := p.advance()

	right := parseExpr(p, bp)

	_, end := right.GetPos()

	return ast.AssignmentExpr{
		BaseStmt: ast.BaseStmt{
			Kind:     ast.ASSIGNMENT_EXPRESSION,
			StartPos: start,
			EndPos:   end,
		},
		Assigne:  identifier,
		Operator: operator,
		Value:    right,
	}
}

// parseStructInstantiationExpr parses a struct instantiation expression, which creates a new instance of a struct.
// It expects the left-hand side to be an identifier representing the struct type, followed by a block of property assignments
// enclosed in curly braces. The function returns an ast.StructInstantiationExpr representing the parsed expression.
func parseStructInstantiationExpr(p *Parser, left ast.Expression) ast.Expression {

	start := p.currentToken().StartPos

	// Check if left is an Identifier
	structName := helpers.ExpectType[ast.IdentifierExpr](left).Identifier

	var properties = map[string]ast.Expression{}
	var methods = map[string]ast.FunctionDeclStmt{}

	p.expect(lexer.OPEN_CURLY)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		var propName = p.expect(lexer.IDENTIFIER).Value
		p.expect(lexer.COLON)
		expr := parseExpr(p, LOGICAL)

		properties[propName] = expr

		if p.currentTokenKind() != lexer.CLOSE_CURLY {
			p.expect(lexer.COMMA)
		}
	}

	end := p.expect(lexer.CLOSE_CURLY).EndPos

	return ast.StructInstantiationExpr{
		BaseStmt: ast.BaseStmt{
			Kind:     ast.STRUCT_LITERAL,
			StartPos: start,
			EndPos:   end,
		},
		StructName: structName,
		Properties: properties,
		Methods:    methods,
	}
}

// parseArrayExpr parses an array expression in the input stream.
// It expects the opening '[' bracket, parses the array elements,
// and returns an ast.ArrayLiterals iNode representing the array.
func parseArrayExpr(p *Parser) ast.Expression {

	start := p.currentToken().StartPos

	p.expect(lexer.OPEN_BRACKET)

	elements := []ast.Expression{}

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_BRACKET {
		elements = append(elements, parseExpr(p, PRIMARY))
		if p.currentTokenKind() != lexer.CLOSE_BRACKET {
			p.expect(lexer.COMMA)
		}
	}

	end := p.expect(lexer.CLOSE_BRACKET).EndPos

	return ast.ArrayLiterals{
		BaseStmt: ast.BaseStmt{
			Kind:     ast.ARRAY_LITERALS,
			StartPos: start,
			EndPos:   end,
		},
		Elements: elements,
		Size:     uint64(len(elements)),
	}
}
