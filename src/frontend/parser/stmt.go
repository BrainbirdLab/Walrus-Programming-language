package parser

import (
	"fmt"
	"rexlang/frontend/ast"
	"rexlang/frontend/lexer"
	"rexlang/utils"
)

func parse_stmt(p *Parser) ast.Stmt {

	stmt_fn, exists := stmt_lu[p.currentTokenKind()]

	if exists {
		return stmt_fn(p)
	}

	return parse_expression_stmt(p)
}

func parse_expression_stmt(p *Parser) ast.ExpressionStmt {

	expression := parse_expr(p, default_bp)

	p.expect(lexer.SEMI_COLON)

	return ast.ExpressionStmt{
		Kind:       ast.STATEMENT,
		Expression: expression,
	}
}

func parse_var_decl_stmt(p *Parser) ast.Stmt {

	var explicitType ast.Type
	var assignedValue ast.Expr

	isConstant := p.advance().Kind == lexer.CONST

	//varName := p.expectError(lexer.IDENTIFIER, "Expected identifier after " + (isConstant ? "const" : "let")  ).Value
	errMsg := fmt.Sprintf("Expected identifier after %s", utils.IF(isConstant, "const", "let"))

	varName := p.expectError(lexer.IDENTIFIER, errMsg).Value

	//p.expectError(lexer.COLON, "Expected type or value after variable name")
	if p.currentTokenKind() != lexer.COLON {
		// then we expect wallrus
		p.expect(lexer.WALRUS)
		// then we expect value
		if p.currentTokenKind() == lexer.SEMI_COLON {
			panic("Expected value after := operator")
		}

		assignedValue = parse_expr(p, default_bp)

		if assignedValue == nil {
			panic("Expected value after := operator")
		}
	} else {
		// then we expect type
		p.advance()
		explicitType = parse_type(p, default_bp)
		if p.currentTokenKind() == lexer.ASSIGNMENT {
			// then we expect assignment
			p.advance()
			assignedValue = parse_expr(p, default_bp)
		}
	}

	p.expect(lexer.SEMI_COLON)

	return ast.VariableDclStml{
		Kind:         ast.VARIABLE_DECLARATION_STATEMENT,
		IsConstant:   isConstant,
		Identifier:   varName,
		Value:        assignedValue,
		ExplicitType: explicitType,
	}
}

func parse_block(p *Parser) ast.BlockStmt {

	p.expect(lexer.OPEN_CURLY)

	body := make([]ast.Stmt, 0)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		body = append(body, parse_stmt(p))
	}

	p.expect(lexer.CLOSE_CURLY)

	return ast.BlockStmt{
		Kind: ast.BLOCK_STATEMENT,
		Body: body,
	}
}

func parse_function_decl_stmt(p *Parser) ast.Stmt {
	p.expect(lexer.FUNCTION)

	functionName := p.expect(lexer.IDENTIFIER).Value
	//parse parameters
	params := parse_params(p)

	// if there is a ARROW token, then we have explicit return type. else we have implicit return type of void
	var explicitReturnType ast.Type
	if p.currentTokenKind() == lexer.ARROW {
		p.advance()
		explicitReturnType = parse_type(p, default_bp)
	} else {
		explicitReturnType = ast.VoidType{}
	}

	// parse block
	functionBody := parse_block(p)

	return ast.FunctionDeclStmt{
		Kind:         ast.FN_DECLARATION_STATEMENT,
		FunctionName: functionName,
		Parameters:   params,
		Block:        functionBody,
		ReturnType:   explicitReturnType,
	}
}

func parse_return_stmt(p *Parser) ast.Stmt {
	p.expect(lexer.RETURN)

	var value ast.Expr

	if p.currentTokenKind() != lexer.SEMI_COLON {
		value = parse_expr(p, default_bp)
	} else {
		value = ast.VoidExpr{}
	}

	p.expect(lexer.SEMI_COLON)

	return ast.ReturnStmt{
		Kind:       ast.RETURN_STATEMENT,
		Expression: value,
	}
}

func parse_struct_decl_stmt(p *Parser) ast.Stmt {

	p.expect(lexer.STRUCT)

	properties := map[string]ast.StructProperty{}
	methods := map[string]ast.StructMethod{}
	structName := p.expect(lexer.IDENTIFIER).Value

	p.expect(lexer.OPEN_CURLY)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {

		var IsStatic bool
		var IsPublic bool
		var ReadOnly bool
		var propname string

		//property
		if p.currentTokenKind() == lexer.ACCESS_MODIFIER {

			if p.currentToken().Value == "pub" {
				IsPublic = true
			} else {
				IsPublic = false
			}

			p.advance() //pass the access modifier

			if p.currentTokenKind() == lexer.STATIC {
				IsStatic = true
				p.advance()
			} else {
				IsStatic = false
			}

			if p.currentTokenKind() == lexer.READONLY {
				ReadOnly = true
				p.advance()
			} else {
				ReadOnly = false
			}

			propname = p.expect(lexer.IDENTIFIER).Value

			if p.currentTokenKind() == lexer.COLON {
				//then its a property

				p.advance()

				propertyType := parse_type(p, default_bp)

				p.expect(lexer.SEMI_COLON)

				//check if already exists
				if _, exists := properties[propname]; exists {
					panic(fmt.Sprintf("Property %s already declared", propname))
				}

				properties[propname] = ast.StructProperty{
					IsStatic: IsStatic,
					IsPublic: IsPublic,
					ReadOnly: ReadOnly,
					Type:     propertyType,
					//Value: nil,
				}

			} else if p.currentTokenKind() == lexer.OPEN_PAREN {
				//then its a method
				//parse the params
				params := parse_params(p)

				var returnType ast.Type
				//check if return type is present, else use void
				if p.currentTokenKind() == lexer.ARROW {
					// so then we expect return type
					p.advance()
					returnType = parse_type(p, default_bp)
				} else {
					returnType = ast.NullType{}
				}

				p.expect(lexer.SEMI_COLON)

				//check if already exists
				if _, exists := methods[propname]; exists {
					panic(fmt.Sprintf("Method %s already declared", propname))
				}

				methods[propname] = ast.StructMethod{
					IsStatic:   IsStatic,
					IsPublic:   IsPublic,
					Parameters: params,
					ReturnType: returnType,
				}
			}

			continue
		} else {
			panic("Expected access modifier of property or method")
		}
	}

	p.expect(lexer.CLOSE_CURLY)

	return ast.StructDeclStatement{
		Kind:       ast.STRUCT_DECLARATION_STATEMENT,
		Properties: properties,
		Methods:    methods,
		StructName: structName,
	}
}

func parse_params(p *Parser) map[string]ast.Type {
	params := map[string]ast.Type{}
	//while )
	p.advance()
	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_PAREN {

		paramName := p.expect(lexer.IDENTIFIER).Value

		p.expect(lexer.COLON)

		paramType := parse_type(p, default_bp)

		//add to the map
		params[paramName] = paramType

		if p.currentTokenKind() != lexer.CLOSE_PAREN {
			p.expect(lexer.COMMA)
		}
	}
	p.expect(lexer.CLOSE_PAREN)
	return params
}

func parse_if_statement(p *Parser) ast.Stmt {

	p.advance()

	condition := parse_expr(p, assignment) // using assignment as the lowest binding power

	consequentBlock := parse_block(p)

	var alternate ast.Stmt

	if p.currentTokenKind() == lexer.ELSE {
		p.advance()
		alternate = parse_block(p)
	} else if p.currentTokenKind() == lexer.ELSEIF {
		//p.advance()
		alternate = parse_if_statement(p)
	} else {
		p.expect(lexer.CLOSE_CURLY)
	}

	return ast.IfStmt{
		Kind:      ast.IF_STATEMENT,
		Condition: condition,
		Block:     consequentBlock,
		Alternate: alternate,
	}
}
