package ast

import (
	"rexlang/frontend/lexer"
	"rexlang/helpers"
)

type NODE_TYPE string

const (
	// Statements
	PROGRAM NODE_TYPE 							= "PROGRAM"
	IMPORT_STATEMENT                         	= "IMPORT_STATEMENT"
	STATEMENT                                	= "STATEMENT"
	BLOCK_STATEMENT                          	= "BLOCK_STATEMENT"
	VARIABLE_DECLARATION_STATEMENT           	= "VARIABLE_DECLARATION_STATEMENT"
	CONTROL_FLOW_STATEMENT                   	= "CONTROL_FLOW_STATEMENT"
	WHILE_STATEMENT                          	= "WHILE_STATEMENT"
	FOR_STATEMENT                            	= "FOR_STATEMENT"
	IF_STATEMENT                             	= "IF_STATEMENT"
	ELSE_STATEMENT                           	= "ELSE_STATEMENT"
	FN_DECLARATION_STATEMENT                 	= "FN_DECLARATION_STATEMENT"
	RETURN_STATEMENT                         	= "RETURN_STATEMENT"
	STRUCT_DECLARATION_STATEMENT             	= "STRUCT_DECLARATION_STATEMENT"

	// Literals
	NUMERIC_LITERAL   						 	= "NUMERIC_LITERAL"
	STRING_LITERAL    						 	= "STRING_LITERAL"
	CHARACTER_LITERAL 						 	= "CHARACTER_LITERAL"
	BOOLEAN_LITERAL   						 	= "BOOLEAN_LITERAL"
	NULL_LITERAL      						 	= "NULL_LITERAL"
	VOID_LITERAL      						 	= "VOID_LITERAL"
	ARRAY_LITERALS    						 	= "ARRAY_LITERALS"

	// Expressions
	ASSIGNMENT_EXPRESSION 						= "ASSIGNMENT_EXPRESSION"
	IDENTIFIER            						= "IDENTIFIER"
	BINARY_EXPRESSION     						= "BINARY_EXPRESSION"
	LOGICAL_EXPRESSION    						= "LOGICAL_EXPRESSION"

	FUNCTION_CALL_EXPRESSION 					= "FUNCTION_CALL_EXPRESSION"

	// Unary Operations
	UNARY_EXPRESSION 							= "UNARY_EXPRESSION"
)

type Node interface {
	node()
}

type Stmt interface {
	stmt()
	GetPos() (lexer.Position, lexer.Position)
}

type Expr interface {
	expr()
	GetPos() (lexer.Position, lexer.Position)
}

type Type interface {
	_type()
}

func ExpectExpr[T Expr](expr Expr) T {
	return helpers.ExpectType[T](expr)
}

func ExpectStmt[T Stmt](expr Stmt) T {
	return helpers.ExpectType[T](expr)
}

func ExpectType[T Type](_type Type) T {
	return helpers.ExpectType[T](_type)
}