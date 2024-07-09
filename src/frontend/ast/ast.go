package ast

import (
	"walrus/frontend/lexer"
)

type NODE_TYPE string

const (
	// Statements
	PROGRAM NODE_TYPE = "program"

	MODULE_STATEMENT NODE_TYPE = "module statement"
	IMPORT_STATEMENT NODE_TYPE = "import statement"

	BLOCK_STATEMENT                NODE_TYPE = "block statement"
	VARIABLE_DECLARATION_STATEMENT NODE_TYPE = "variable declaration statement"
	CONTROL_FLOW_STATEMENT         NODE_TYPE = "control flow statement"
	WHILE_STATEMENT                NODE_TYPE = "while loop statement"
	SWITCH_STATEMENT               NODE_TYPE = "switch statement"
	SWITCH_CASE_STATEMENT          NODE_TYPE = "switch case statement"
	DEFAULT_CASE_STATEMENT         NODE_TYPE = "default case statement"
	IF_STATEMENT                   NODE_TYPE = "if statement"
	ELSE_STATEMENT                 NODE_TYPE = "else statement"
	FOR_LOOP_STATEMENT             NODE_TYPE = "for loop statement"
	FOREACH_LOOP_STATEMENT         NODE_TYPE = "foreach loop statement"
	FN_DECLARATION_STATEMENT       NODE_TYPE = "fn declaration statement"
	FN_PROTOTYPE_STATEMENT         NODE_TYPE = "fn prototype statement"
	RETURN_STATEMENT               NODE_TYPE = "return statement"
	BREAK_STATEMENT                NODE_TYPE = "break statement"
	CONTINUE_STATEMENT             NODE_TYPE = "continue statement"
	TRAIT_STATEMENT                NODE_TYPE = "trait statement"
	STRUCT_STATEMENT               NODE_TYPE = "struct statement"
	IMPLEMENTS_STATEMENT           NODE_TYPE = "implements statement"

	// Literals
	NUMERIC_LITERAL   NODE_TYPE = "NUMERIC_LITERAL"
	INTEGER_LITERAL   NODE_TYPE = "integer literal"
	FLOAT_LITERAL     NODE_TYPE = "float literal"
	STRING_LITERAL    NODE_TYPE = "string literal"
	CHARACTER_LITERAL NODE_TYPE = "character literal"
	BOOLEAN_LITERAL   NODE_TYPE = "boolean literal"
	NULL_LITERAL      NODE_TYPE = "null literal"
	VOID_LITERAL      NODE_TYPE = "void literal"
	ARRAY_LITERALS    NODE_TYPE = "array literals"
	STRUCT_LITERAL    NODE_TYPE = "struct literal"

	STRUCT_PROPERTY NODE_TYPE = "struct property"

	// Expressions
	ASSIGNMENT_EXPRESSION NODE_TYPE = "assignment expression"
	IDENTIFIER            NODE_TYPE = "identifier"
	BINARY_EXPRESSION     NODE_TYPE = "binary expression"
	LOGICAL_EXPRESSION    NODE_TYPE = "logical expression"

	// Functions
	FUNCTION_PARAMETER NODE_TYPE = "function parameter"

	FUNCTION_CALL_EXPRESSION NODE_TYPE = "function call expression"

	// Unary Operations
	UNARY_EXPRESSION NODE_TYPE = "unary expression"
)

type Node interface {
	INodeType() NODE_TYPE
	GetPos() (lexer.Position, lexer.Position)
}

type Statement interface {
	Node
	iStatement()
}

type Expression interface {
	Node
	iExpression()
}

type Type interface {
	IType() DATA_TYPE
}

type BaseStmt struct {
	Kind     NODE_TYPE
	StartPos lexer.Position
	EndPos   lexer.Position
}
