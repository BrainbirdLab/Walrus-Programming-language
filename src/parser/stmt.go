package parser

import (
	"fmt"
	"rexlang/ast"
	"rexlang/lexer"
)

func parse_stmt(p *Parser) ast.Stmt {

	stmt_fn, exists := stmt_lu[p.currentTokenKind()]

	if exists {
		return stmt_fn(p)
	}

	expression := parse_expr(p, default_bp)

	p.expect(lexer.SEMI_COLON)

	return ast.ExpressionStmt{
		Expression: expression,
	}
}

func parse_var_decl_stmt(p *Parser) ast.Stmt {

	isConstant := p.advance().Kind == lexer.CONST

	//varName := p.expectError(lexer.IDENTIFIER, "Expected identifier after " + (isConstant ? "const" : "let")  ).Value
	errMsg := fmt.Sprintf("Expected identifier after %s", func() string {
		if isConstant {
			return "const"
		}
		return "let"
	}())

	
	varName := p.expectError(lexer.IDENTIFIER, errMsg).Value
	p.expect(lexer.ASSIGNMENT)
	value := parse_expr(p, assignment)
	p.expect(lexer.SEMI_COLON)

	return ast.VarDeclStmt{
		IsConstant: isConstant,
		Name: varName,
		Value: value,
	}
}