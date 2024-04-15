package ast

import (
    "rexlang/helpers"
)

type NodeType int

const (
    // Statements
    PROGRAM int = iota
    STATEMENT         			
    BLOCK_STATEMENT    			
    VARIABLE_DECLARATION_STATEMENT
    CONTROL_FLOW_STATEMENT       
    WHILE_STATEMENT    			
    fOR_STATEMENT      			
    IF_STATEMENT       			
    ELSE_STATEMENT   			 			

    // Literals
    NUMERIC_LITERAL				
    STRING_LITERAL 				
    BOOLEAN_LITERAL			
    NULL_LITERAL				

    // Expressions
    ASSIGNMENT_EXPRESSION
    IDENTIFIER	
    BINARY_EXPRESSION
    LOGICAL_EXPRESSION    		

    // Unary Operations
    UNARY_EXPRESSION			
)

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