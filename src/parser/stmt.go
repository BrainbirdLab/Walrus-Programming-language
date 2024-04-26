package parser

import (
	"fmt"
	"rexlang/ast"
	"rexlang/lexer"
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
	} else {
		// then we expect type
		p.advance()
		explicitType = parse_type(p, default_bp)
		fmt.Printf("Explicit type: %v\n", explicitType)
		p.expect(lexer.ASSIGNMENT)
	}

	fmt.Printf("Explicit type: %v\n", explicitType)

	fmt.Printf("Current token: %s\n", p.currentToken().Value)

	if p.currentTokenKind() != lexer.SEMI_COLON {
		// then we expect assignment
		//p.expect(lexer.ASSIGNMENT)
		assignedValue = parse_expr(p, default_bp)
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

func parse_struct_decl_stmt(p *Parser) ast.Stmt {

	p.expect(lexer.STRUCT)

	properties := map[string]ast.StructProperty{}
	methods := map[string]ast.StructMethod{}
	structName := p.expect(lexer.IDENTIFIER).Value

	p.expect(lexer.OPEN_CURLY)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {

		var IsStatic bool
		var propName string

		if p.currentTokenKind() == lexer.STATIC {
			IsStatic = true
			p.expect(lexer.STATIC) //! or use p.advance()
		}

		//property
		if p.currentTokenKind() == lexer.IDENTIFIER {
			propName = p.expect(lexer.IDENTIFIER).Value
			p.expectError(lexer.COLON, "Expected ':' after property name")
			structType := parse_type(p, default_bp)
			p.expect(lexer.COMMA)

			_, exists := properties[propName]
			if exists {
				panic(fmt.Sprintf("Property %s already exists", propName))
			}

			properties[propName] = ast.StructProperty{
				IsStatic: IsStatic,
				Type:     structType,
				PropertyName: propName,
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
