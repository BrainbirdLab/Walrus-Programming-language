package parser

import (
	"fmt"
	"rexlang/frontend/ast"
	"rexlang/frontend/lexer"
	"rexlang/utils"
)

func parse_stmt(p *Parser) ast.Stmt {

	stmt_fn, exists := stmtLookup[p.currentTokenKind()]

	if exists {
		return stmt_fn(p)
	}

	return parse_expression_stmt(p)
}

func parse_expression_stmt(p *Parser) ast.ExpressionStmt {

	start := p.currentToken().StartPos

	expression := parse_expr(p, DEFAULT_BP)

	end := p.expect(lexer.SEMI_COLON).EndPos

	return ast.ExpressionStmt{
		Kind:       ast.STATEMENT,
		Expression: expression,
		StartPos:  start,
		EndPos:  end,
	}
}

func parse_var_decl_stmt(p *Parser) ast.Stmt {

	start := p.currentToken().StartPos

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

		assignedValue = parse_expr(p, DEFAULT_BP)

		if assignedValue == nil {
			panic("Expected value after := operator")
		}
	} else {
		// then we expect type
		p.advance()
		explicitType = parse_type(p, DEFAULT_BP)
		if p.currentTokenKind() == lexer.ASSIGNMENT {
			// then we expect assignment
			p.advance()
			assignedValue = parse_expr(p, DEFAULT_BP)
		}
	}

	end := p.expect(lexer.SEMI_COLON).EndPos

	return ast.VariableDclStml{
		Kind:         ast.VARIABLE_DECLARATION_STATEMENT,
		IsConstant:   isConstant,
		Identifier:   varName,
		Value:        assignedValue,
		ExplicitType: explicitType,
		StartPos:    start,
		EndPos:    end,
	}
}

func parse_block(p *Parser) ast.BlockStmt {

	start := p.expect(lexer.OPEN_CURLY).StartPos

	body := make([]ast.Stmt, 0)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		body = append(body, parse_stmt(p))
	}

	end := p.expect(lexer.CLOSE_CURLY).EndPos

	return ast.BlockStmt{
		Kind: ast.BLOCK_STATEMENT,
		Body: body,
		StartPos: start,
		EndPos: end,
	}
}

func parse_function_decl_stmt(p *Parser) ast.Stmt {

	start := p.currentToken().StartPos

	p.expect(lexer.FUNCTION)

	functionName := p.expect(lexer.IDENTIFIER).Value
	//parse parameters
	params := parse_params(p)

	// if there is a ARROW token, then we have explicit return type. else we have implicit return type of void
	var explicitReturnType ast.Type
	if p.currentTokenKind() == lexer.ARROW {
		p.advance()
		explicitReturnType = parse_type(p, DEFAULT_BP)
	} else {
		explicitReturnType = ast.VoidType{}
	}

	// parse block
	functionBody := parse_block(p)

	end := functionBody.EndPos

	return ast.FunctionDeclStmt{
		Kind:         ast.FN_DECLARATION_STATEMENT,
		FunctionName: functionName,
		Parameters:   params,
		Block:        functionBody,
		ReturnType:   explicitReturnType,
		StartPos:    start,
		EndPos:    end,
	}
}

func parse_return_stmt(p *Parser) ast.Stmt {

	start := p.currentToken().StartPos

	p.expect(lexer.RETURN)

	var value ast.Expr

	if p.currentTokenKind() != lexer.SEMI_COLON {
		value = parse_expr(p, DEFAULT_BP)
	} else {
		value = ast.VoidExpr{}
	}

	end := p.expect(lexer.SEMI_COLON).EndPos

	return ast.ReturnStmt{
		Kind:       ast.RETURN_STATEMENT,
		Expression: value,
		StartPos:  start,
		EndPos:  end,
	}
}

func parse_struct_decl_stmt(p *Parser) ast.Stmt {

	start := p.currentToken().StartPos

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
		if p.currentTokenKind() == lexer.ACCESS {

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

				propertyType := parse_type(p, DEFAULT_BP)

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
					returnType = parse_type(p, DEFAULT_BP)
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

	end := p.expect(lexer.CLOSE_CURLY).EndPos

	return ast.StructDeclStatement{
		Kind:       ast.STRUCT_DECLARATION_STATEMENT,
		Properties: properties,
		Methods:    methods,
		StructName: structName,
		StartPos:  start,
		EndPos:  end,
	}
}

func parse_params(p *Parser) map[string]ast.Type {
	params := map[string]ast.Type{}
	//while )
	p.advance() // pass the open paren

	//parse the parameters
	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_PAREN {

		paramName := p.expect(lexer.IDENTIFIER).Value

		p.expect(lexer.COLON)

		paramType := parse_type(p, DEFAULT_BP)

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

	start := p.advance().StartPos

	condition := parse_expr(p, ASSIGNMENT) // using assignment as the lowest binding power

	consequentBlock := parse_block(p)

	var alternate ast.Stmt

	var end lexer.Position

	if p.currentTokenKind() == lexer.ELSE {
		p.advance()
		block := parse_block(p)
		alternate = block
		end = block.EndPos
	} else if p.currentTokenKind() == lexer.ELSEIF {
		//p.advance()
		stmt := parse_if_statement(p)
		alternate = stmt
		// type cast
		end = stmt.(ast.IfStmt).EndPos
	} else {
		end = p.expect(lexer.CLOSE_CURLY).EndPos
	}

	return ast.IfStmt{
		Kind:      ast.IF_STATEMENT,
		Condition: condition,
		Block:     consequentBlock,
		Alternate: alternate,
		StartPos: start,
		EndPos: end,
	}
}
