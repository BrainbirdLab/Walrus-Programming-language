package ast

import (
	"rexlang/frontend/lexer"
)

type BinaryExpr struct {
	Kind     NodeType
	Operator lexer.Token
	Left     Expr
	Right    Expr
}
func (b BinaryExpr) expr() {}

type UnaryExpr struct {
	Kind     NodeType
	Operator lexer.Token
	Argument Expr
}
func (u UnaryExpr) expr() {}

type SymbolExpr struct {
	Kind   NodeType
	Symbol string
	Type   string
}
func (i SymbolExpr) expr() {}

type NumericLiteral struct {
	Kind  NodeType
	Value float64
	Type  string
}
func (n NumericLiteral) expr() {}

type StringLiteral struct {
	Kind  NodeType
	Value string
	Type  string
}
func (s StringLiteral) expr() {}

type BooleanLiteral struct {
	Kind  NodeType
	Value bool
	Type  string
}
func (b BooleanLiteral) expr() {}

type NullLiteral struct {
	Kind  NodeType
	Value string
	Type  string
}
func (n NullLiteral) expr() {}

type VoidExpr struct {
	Kind  NodeType
	Value string
	Type  string
}

func (v VoidExpr) expr() {}

type AssignmentExpr struct {
	Kind     NodeType
	Assigne  SymbolExpr
	Value    Expr
	Operator lexer.Token
}
func (a AssignmentExpr) expr() {}

type FunctionCallExpr struct {
  Kind     		NodeType
  Function 		Expr
  Args     		[]Expr
}
func (c FunctionCallExpr) expr() {}

type StructInstantiationExpr struct {
	StructName string
	Properties map[string]Expr
	Methods    map[string]FunctionDeclStmt
}
func (s StructInstantiationExpr) expr() {}

type ArrayLiterals struct {
	Kind          	NodeType
	Size     		uint64
	Elements 		[]Expr
}
func (a ArrayLiterals) expr() {}
