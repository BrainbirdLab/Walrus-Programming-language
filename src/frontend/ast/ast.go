package ast

import (
	"walrus/frontend/lexer"
)

type NODE_TYPE string

const (
	// Statements
	PROGRAM NODE_TYPE 					= "PROGRAM"

	MODULE_STATEMENT 					= "MODULE_STATEMENT"
	IMPORT_STATEMENT 					= "IMPORT_STATEMENT"

	BLOCK_STATEMENT                		= "BLOCK_STATEMENT"
	VARIABLE_DECLARATION_STATEMENT 		= "VARIABLE_DECLARATION_STATEMENT"
	CONTROL_FLOW_STATEMENT         		= "CONTROL_FLOW_STATEMENT"
	WHILE_STATEMENT                		= "WHILE_STATEMENT"
	SWITCH_CASE_STATEMENT				= "SWITCH_CASE_STATEMENT"
	IF_STATEMENT                   		= "IF_STATEMENT"
	ELSE_STATEMENT                 		= "ELSE_STATEMENT"
	FOR_LOOP_STATEMENT             		= "FOR_LOOP_STATEMENT"
	FOREACH_LOOP_STATEMENT         		= "FOREACH_LOOP_STATEMENT"
	FN_DECLARATION_STATEMENT       		= "FN_DECLARATION_STATEMENT"
	FN_PROTOTYPE_STATEMENT         		= "FN_PROTOTYPE_STATEMENT"
	RETURN_STATEMENT               		= "RETURN_STATEMENT"
	TRAIT_STATEMENT    					= "TRAIT_STATEMENT"
	STRUCT_STATEMENT   					= "STRUCT_STATEMENT"
	IMPLEMENTS_STATEMENT           		= "IMPLEMENTS_STATEMENT"

	// Literals
	NUMERIC_LITERAL   					= "NUMERIC_LITERAL"
	STRING_LITERAL    					= "STRING_LITERAL"
	CHARACTER_LITERAL 					= "CHARACTER_LITERAL"
	BOOLEAN_LITERAL   					= "BOOLEAN_LITERAL"
	NULL_LITERAL      					= "NULL_LITERAL"
	VOID_LITERAL      					= "VOID_LITERAL"
	ARRAY_LITERALS    					= "ARRAY_LITERALS"
	STRUCT_LITERAL    					= "STRUCT_LITERAL"

	STRUCT_PROPERTY 					= "STRUCT_PROPERTY"

	// Expressions
	ASSIGNMENT_EXPRESSION 				= "ASSIGNMENT_EXPRESSION"
	IDENTIFIER            				= "IDENTIFIER"
	BINARY_EXPRESSION     				= "BINARY_EXPRESSION"
	LOGICAL_EXPRESSION    				= "LOGICAL_EXPRESSION"

	FUNCTION_CALL_EXPRESSION 			= "FUNCTION_CALL_EXPRESSION"

	// Unary Operations
	UNARY_EXPRESSION 					= "UNARY_EXPRESSION"
)

type Node interface {
	node()
	GetPos() (lexer.Position, lexer.Position)
}

type Statement interface {
	Node
	_statement()
}

type Expression interface {
	Node
	_expression()
}

type Type interface {
	_type()
}

type BaseStmt struct {
	Kind     NODE_TYPE
	StartPos lexer.Position
	EndPos   lexer.Position
}
