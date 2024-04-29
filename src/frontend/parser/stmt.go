package parser

import (
	"fmt"
	"rexlang/frontend/ast"
	"rexlang/frontend/lexer"
	"rexlang/utils"
)

func parse_stmt(p *Parser) ast.Stmt {

	fmt.Printf("Current token in parse_stmt: %s\n", p.currentToken().Value)

	stmt_fn, exists := stmt_lu[p.currentTokenKind()]

	if exists {
		return stmt_fn(p)
	}

	fmt.Printf("Current token: %s does not have a statement handler\n", p.currentToken().Value)
	return parse_expression_stmt(p)
}

func parse_expression_stmt(p *Parser) ast.ExpressionStmt {

	expression := parse_expr(p, default_bp)

	p.expect(lexer.SEMI_COLON)

	fmt.Printf("Parsed expression: %v\n", expression)

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

	fmt.Printf("Variable name: %s\n", varName)

	//p.expectError(lexer.COLON, "Expected type or value after variable name")
	if p.currentTokenKind() != lexer.COLON {
		// then we expect wallrus
		p.expect(lexer.WALRUS)
		// then we expect value
		if (p.currentTokenKind() == lexer.SEMI_COLON) {
			panic("Expected value after := operator")
		}
		assignedValue = parse_expr(p, default_bp)
		fmt.Printf("Assigned value: %v\n", assignedValue)
		if assignedValue == nil {
			panic("Expected value after := operator")
		}
	} else {
		// then we expect type
		p.advance()
		explicitType = parse_type(p, default_bp)
		fmt.Printf("Explicit type: %v\n", explicitType)
		if p.currentTokenKind() == lexer.ASSIGNMENT {
			// then we expect assignment
			p.advance()
			assignedValue = parse_expr(p, default_bp)
		}
	}

	fmt.Printf("Explicit type: %v\n", explicitType)
	
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

func parse_struct_decl_stmt(p *Parser) ast.Stmt {

	p.expect(lexer.STRUCT)

	properties := map[string]ast.StructProperty{}
	methods := map[string]ast.StructMethod{}
	structName := p.expect(lexer.IDENTIFIER).Value

	p.expect(lexer.OPEN_CURLY)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {

		var IsStatic bool
		var IsPublic bool
		var propname string

		if p.currentTokenKind() == lexer.STATIC {
			IsStatic = true
			//p.expect(lexer.STATIC) //! or use p.advance()
			p.advance() //advance to the next token
		}

		//property
		if p.currentTokenKind() == lexer.ACCESS_MODIFIER {

			if p.currentToken().Value == "pub" {
				IsPublic = true
				p.advance()
			} else {
				IsPublic = false
				p.advance()
			}

			propname = p.expect(lexer.IDENTIFIER).Value

			if p.currentTokenKind() == lexer.COLON {
				//then its a property

				p.advance()

				//fmt.Printf("Parsed property type: %v\n", parse_type(p, default_bp))
				propertyType := parse_type(p, default_bp)

				p.expect(lexer.SEMI_COLON)

				//check if already exists
				if _, exists := properties[propname]; exists {
					panic(fmt.Sprintf("Property %s already declared", propname))
				}

				properties[propname] = ast.StructProperty{
					IsStatic: IsStatic,
					IsPublic: IsPublic,
					Type: propertyType,
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
					IsStatic: IsStatic,
					IsPublic: IsPublic,
					Parameters: params,
					ReturnType: returnType,
				}
			}

			continue
		}
	}

	p.expect(lexer.CLOSE_CURLY)

	return ast.StructDeclStatement{
		Properties: properties,
		Methods: methods,
		StructName: structName,
		
	}
}

func parse_params(p *Parser) map[string]ast.Type {
	params := map[string]ast.Type{}
	//while )
	p.advance()
	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_PAREN {

		paramName := p.expect(lexer.IDENTIFIER).Value
		fmt.Printf("Param name: %s\n", paramName)

		p.expect(lexer.COLON)

		paramType := parse_type(p, default_bp)
		fmt.Printf("Param type: %v\n", paramType)
		
		
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

	//pass the if token
	fmt.Printf("Current token in if: %v\n", p.currentToken())

	p.advance()

	fmt.Printf("Remaining tokens: %v\n", p.currentToken())

	condition := parse_expr(p, assignment) // using assignment as the lowest binding power

	fmt.Println("Condition: ", condition)

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
