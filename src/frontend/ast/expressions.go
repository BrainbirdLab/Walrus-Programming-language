package ast

import (
	"rexlang/frontend/lexer"
)

type BinaryExpr struct {
	Kind     NODE_TYPE
	Operator lexer.Token
	Left     Expr
	Right    Expr
	StartPos lexer.Position
	EndPos   lexer.Position
}
func (b BinaryExpr) expr() {}
func (b BinaryExpr) GetPos() (lexer.Position, lexer.Position) {
	return b.StartPos, b.EndPos
}


type UnaryExpr struct {
	Kind     NODE_TYPE
	Operator lexer.Token
	Argument Expr
	StartPos lexer.Position
	EndPos   lexer.Position
}
func (u UnaryExpr) expr() {}
func (u UnaryExpr) GetPos() (lexer.Position, lexer.Position) {
	return u.StartPos, u.EndPos
}

type SymbolExpr struct {
	Kind   NODE_TYPE
	Symbol string
	Type   string
	StartPos lexer.Position
	EndPos   lexer.Position
}
func (i SymbolExpr) expr() {}
func (i SymbolExpr) GetPos() (lexer.Position, lexer.Position) {
	return i.StartPos, i.EndPos
}

type NumericLiteral struct {
	Kind  NODE_TYPE
	Value float64
	Type  string
	StartPos lexer.Position
	EndPos   lexer.Position
}

func (n NumericLiteral) expr() {}
func (n NumericLiteral) GetPos() (lexer.Position, lexer.Position) {
	return n.StartPos, n.EndPos
}

type StringLiteral struct {
	Kind  NODE_TYPE
	Value string
	Type  string
	StartPos lexer.Position
	EndPos   lexer.Position
}

func (s StringLiteral) expr() {}
func (s StringLiteral) GetPos() (lexer.Position, lexer.Position) {
	return s.StartPos, s.EndPos
}

type BooleanLiteral struct {
	Kind  NODE_TYPE
	Value bool
	Type  string
	StartPos lexer.Position
	EndPos   lexer.Position
}

func (b BooleanLiteral) expr() {}
func (b BooleanLiteral) GetPos() (lexer.Position, lexer.Position) {
	return b.StartPos, b.EndPos
}

type NullLiteral struct {
	Kind  NODE_TYPE
	Value string
	Type  string
	StartPos lexer.Position
	EndPos   lexer.Position
}

func (n NullLiteral) expr() {}
func (n NullLiteral) GetPos() (lexer.Position, lexer.Position) {
	return n.StartPos, n.EndPos
}

type VoidExpr struct {
	Kind  NODE_TYPE
	Value string
	Type  string
}

func (v VoidExpr) expr() {}
func (v VoidExpr) GetPos() (lexer.Position, lexer.Position) {
	return lexer.Position{}, lexer.Position{}
}

type AssignmentExpr struct {
	Kind     NODE_TYPE
	Assigne  SymbolExpr
	Value    Expr
	Operator lexer.Token
	StartPos lexer.Position
	EndPos   lexer.Position
}
func (a AssignmentExpr) expr() {}
func (a AssignmentExpr) GetPos() (lexer.Position, lexer.Position) {
	return a.StartPos, a.EndPos
}

type FunctionCallExpr struct {
	Kind     NODE_TYPE
	Function Expr
	Args     []Expr
	StartPos lexer.Position
	EndPos   lexer.Position
}
func (c FunctionCallExpr) expr() {}
func (c FunctionCallExpr) GetPos() (lexer.Position, lexer.Position) {
	return c.StartPos, c.EndPos
}

type StructInstantiationExpr struct {
	StructName string
	Properties map[string]Expr
	Methods    map[string]FunctionDeclStmt
	StartPos lexer.Position
	EndPos   lexer.Position
}
func (s StructInstantiationExpr) expr() {}
func (s StructInstantiationExpr) GetPos() (lexer.Position, lexer.Position) {
	return s.StartPos, s.EndPos
}

type ArrayLiterals struct {
	Kind     NODE_TYPE
	Size     uint64
	Elements []Expr
	StartPos lexer.Position
	EndPos   lexer.Position
}

func (a ArrayLiterals) expr() {}
func (a ArrayLiterals) GetPos() (lexer.Position, lexer.Position) {
	return a.StartPos, a.EndPos
}