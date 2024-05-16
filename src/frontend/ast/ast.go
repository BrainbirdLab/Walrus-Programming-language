package ast

import (
    "rexlang/helpers"
)

type NodeType int

const (
    // Statements
    PROGRAM NodeType = iota
    STATEMENT         			
    BLOCK_STATEMENT    			
    VARIABLE_DECLARATION_STATEMENT
    CONTROL_FLOW_STATEMENT       
    WHILE_STATEMENT    			
    FOR_STATEMENT      			
    IF_STATEMENT       			
    ELSE_STATEMENT   			 			

    // Literals
    NUMERIC_LITERAL				
    STRING_LITERAL
    CHARACTER_LITERAL			
    BOOLEAN_LITERAL		
    NULL_LITERAL		
    ARRAY_LITERALS

    // Expressions
    ASSIGNMENT_EXPRESSION
    IDENTIFIER
    BINARY_EXPRESSION
    LOGICAL_EXPRESSION

    FUNCTION_EXPRESSION
    
    // Unary Operations
    UNARY_EXPRESSION
)

type Node interface {
    node()
}

type Stmt interface {
	stmt()
}

type Expr interface {
	expr()
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