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

	p.expect(lexer.OPEN_PAREN_TOKEN)

	var arguments []ast.Expression

	for p.currentTokenKind() != lexer.CLOSE_PAREN_TOKEN {
		//parse the arguments
		argument := parseExpr(p, DEFAULT_BP)
		arguments = append(arguments, argument)

		if p.currentTokenKind() == lexer.COMMA_TOKEN {
			p.advance()
		}
	}

	end := p.expect(lexer.CLOSE_PAREN_TOKEN).EndPos

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

	p.expect(lexer.DOT_TOKEN)

	identifier := p.expect(lexer.IDENTIFIER_TOKEN)

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

	if tokenKind == lexer.IDENTIFIER_TOKEN && p.nextToken().Kind == lexer.OPEN_CURLY_TOKEN && (p.previousToken().Kind == lexer.WALRUS_TOKEN || p.previousToken().Kind == lexer.ASSIGNMENT_TOKEN) {
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
	case lexer.INTEGER_TOKEN:

		// A 32-bit integer can store upto -2,147,483,648 to 2,147,483,647
		// it is nearly 10 digits long
		// the raw value is stored as a string
		// so we need to convert it to a number along with the proper size from the string length and the value

		rawValue := p.advance().Value

		size := uint8(32)

		if len(rawValue) > 10 {
			size = 64
		} else {
			// if number is out of range for 32-bit integer
			// then it is a 64-bit integer
			// But to avoid checking both positive and negative ranges, we just check the positive range by using the absolute value
			// if the absolute value is greater than 2,147,483,647 then it is a 64-bit integer
			number, _ := strconv.ParseInt(rawValue, 10, 32)
			if number < 0 {
				number = -number
			}
			if number > 2147483647 {
				size = 64
			}
		}

		return ast.NumericLiteral{
			BaseStmt: ast.BaseStmt{
				Kind:     ast.INTEGER_LITERAL,
				StartPos: startpos,
				EndPos:   endpos,
			},
			Value: rawValue,
			BitSize: size,
		}
	case lexer.FLOATING_TOKEN:
		rawValue := p.advance().Value

		size := uint8(32)
		
		number, _ := strconv.ParseFloat(rawValue, 64)

		if number < 0 {
			number = -number
		}

		// check the floating point decimal size

		decimal := int64(number)

		//max size of a 32-bit floating point number is 7 digits
		if decimal > 9999999 {
			size = 64
		}

		return ast.NumericLiteral{
			BaseStmt: ast.BaseStmt{
				Kind:     ast.FLOAT_LITERAL,
				StartPos: startpos,
				EndPos:   endpos,
			},
			Value: rawValue,
			BitSize: size,
		}

	case lexer.STRING_TOKEN:
		return ast.StringLiteral{
			BaseStmt: ast.BaseStmt{
				Kind:     ast.STRING_LITERAL,
				StartPos: startpos,
				EndPos:   endpos,
			},
			Value: p.advance().Value,
		}
	case lexer.CHARACTER_TOKEN:
		return ast.CharacterLiteral{
			BaseStmt: ast.BaseStmt{
				Kind:     ast.CHARACTER_LITERAL,
				StartPos: startpos,
				EndPos:   endpos,
			},
			Value: p.advance().Value,
		}
	case lexer.IDENTIFIER_TOKEN:
		return ast.IdentifierExpr{
			BaseStmt: ast.BaseStmt{
				Kind:     ast.IDENTIFIER,
				StartPos: startpos,
				EndPos:   endpos,
			},
			Identifier: p.advance().Value,
		}
	case lexer.TRUE_TOKEN:
		p.advance()
		return ast.BooleanLiteral{
			BaseStmt: ast.BaseStmt{
				Kind:     ast.BOOLEAN_LITERAL,
				StartPos: startpos,
				EndPos:   endpos,
			},
			Value: true,
		}
	case lexer.FALSE_TOKEN:
		p.advance()
		return ast.BooleanLiteral{
			BaseStmt: ast.BaseStmt{
				Kind:     ast.BOOLEAN_LITERAL,
				StartPos: startpos,
				EndPos:   endpos,
			},
			Value: false,
		}
	case lexer.NULL_TOKEN:
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
	p.expect(lexer.OPEN_PAREN_TOKEN)
	expression := parseExpr(p, DEFAULT_BP)
	p.expect(lexer.CLOSE_PAREN_TOKEN)
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

	p.expect(lexer.OPEN_CURLY_TOKEN)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY_TOKEN {
		var propName = p.expect(lexer.IDENTIFIER_TOKEN).Value
		p.expect(lexer.COLON_TOKEN)
		expr := parseExpr(p, LOGICAL)

		properties[propName] = expr

		if p.currentTokenKind() != lexer.CLOSE_CURLY_TOKEN {
			p.expect(lexer.COMMA_TOKEN)
		}
	}

	end := p.expect(lexer.CLOSE_CURLY_TOKEN).EndPos

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

	p.expect(lexer.OPEN_BRACKET_TOKEN)

	elements := []ast.Expression{}

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_BRACKET_TOKEN {
		elements = append(elements, parseExpr(p, PRIMARY))
		if p.currentTokenKind() != lexer.CLOSE_BRACKET_TOKEN {
			p.expect(lexer.COMMA_TOKEN)
		}
	}

	end := p.expect(lexer.CLOSE_BRACKET_TOKEN).EndPos

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
