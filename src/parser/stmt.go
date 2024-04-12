package parser

import (
	"fmt"
	"rexlang/ast"
	"rexlang/lexer"
	"rexlang/utils"
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
	errMsg := fmt.Sprintf("Expected identifier after %s", utils.IF(isConstant, "const", "let"))
	
	varName := p.expectError(lexer.IDENTIFIER, errMsg).Value
	
	assignMentOrSemiColon := p.advance()

	var value ast.Expr

	if assignMentOrSemiColon.Kind != lexer.ASSIGNMENT {
		if isConstant {
			panic("Constant must be initialized with a value")
		} else if !isConstant && assignMentOrSemiColon.Kind != lexer.SEMI_COLON{
			panic("Expected value or ; after let")
		}
	} else {
		// No assignment, just a declaration
		value = parse_expr(p, assignment)
		p.expect(lexer.SEMI_COLON)
	}
		

	return ast.VarDeclStmt{
		IsConstant: isConstant,
		Name: varName,
		Value: value,
	}
}