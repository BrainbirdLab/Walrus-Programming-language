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
		Kind:       ast.NodeType(ast.STATEMENT),
		Expression: expression,
	}
}

func parse_var_decl_stmt(p *Parser) ast.Stmt {

	isConstant := p.advance().Kind == lexer.CONST

	//varName := p.expectError(lexer.IDENTIFIER, "Expected identifier after " + (isConstant ? "const" : "let")  ).Value
	errMsg := fmt.Sprintf("Expected identifier after %s", utils.IF(isConstant, "const", "let"))

	varName := p.expectError(lexer.IDENTIFIER, errMsg).Value

	assignmentToken:= p.advance()

	fmt.Printf("Ass token (%s)", assignmentToken.Value)

	var value ast.Expr

	if assignmentToken.Kind != lexer.ASSIGNMENT {
		if isConstant {
			panic("Constant must be initialized with a value")
		} else if !isConstant && (assignmentToken.Kind != lexer.SEMI_COLON) {
			panic("Cannot determine the end of statement")
		}
	} else {
		// No assignment, just a declaration
		fmt.Printf("Parsing variable declaration with assignment for identifier: %s\n", varName)
		value = parse_expr(p, assignment)
		p.expect(lexer.SEMI_COLON)
		//p.expectAny(lexer.ENDLINE, lexer.EOF)
	}

	return ast.VariableDclStml{
		Kind: ast.NodeType(ast.VARIABLE_DECLARATION_STATEMENT),
		IsConstant: isConstant,
		Identifier: varName,
		Value:      value,
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
		Kind: ast.NodeType(ast.BLOCK_STATEMENT),
		Body: body,
	}
}

func parse_if_statement(p *Parser) ast.Stmt{
	
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
		Kind: ast.NodeType(ast.IF_STATEMENT),
		Condition: condition,
		Block: consequentBlock,
		Alternate: alternate,
	}
}