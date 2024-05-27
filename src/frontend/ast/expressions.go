package ast

import (
	"walrus/frontend/lexer"
)

type BinaryExpr struct {
	BaseStmt
	Operator lexer.Token
	Left     Expression
	Right    Expression
}

func (b BinaryExpr) node() {}
func (b BinaryExpr) GetPos() (lexer.Position, lexer.Position) {
	return b.StartPos, b.EndPos
}
func (b BinaryExpr) _expression() {}

type UnaryExpr struct {
	BaseStmt
	Operator lexer.Token
	Argument Expression
}

func (u UnaryExpr) node() {}
func (u UnaryExpr) GetPos() (lexer.Position, lexer.Position) {
	return u.StartPos, u.EndPos
}
func (u UnaryExpr) _expression() {}

type IdentifierExpr struct {
	BaseStmt
	Identifier  string
	Type     	string
}

func (i IdentifierExpr) node() {}
func (i IdentifierExpr) GetPos() (lexer.Position, lexer.Position) {
	return i.StartPos, i.EndPos
}
func (i IdentifierExpr) _expression() {}

type NumericLiteral struct {
	BaseStmt
	Value    float64
	Type     string
}

func (n NumericLiteral) node() {}
func (n NumericLiteral) GetPos() (lexer.Position, lexer.Position) {
	return n.StartPos, n.EndPos
}
func (n NumericLiteral) _expression() {}

type StringLiteral struct {
	BaseStmt
	Value    string
	Type     string
}

func (s StringLiteral) node() {}
func (s StringLiteral) GetPos() (lexer.Position, lexer.Position) {
	return s.StartPos, s.EndPos
}
func (s StringLiteral) _expression() {}

type BooleanLiteral struct {
	BaseStmt
	Value    bool
	Type     string
}

func (b BooleanLiteral) node() {}
func (b BooleanLiteral) GetPos() (lexer.Position, lexer.Position) {
	return b.StartPos, b.EndPos
}
func (b BooleanLiteral) _expression() {}

type NullLiteral struct {
	BaseStmt
	Value    string
	Type     string
}

func (n NullLiteral) node() {}
func (n NullLiteral) GetPos() (lexer.Position, lexer.Position) {
	return n.StartPos, n.EndPos
}
func (n NullLiteral) _expression() {}

type VoidExpr struct {
	Kind  NODE_TYPE
	Value string
	Type  string
}
func (v VoidExpr) node() {}
func (v VoidExpr) GetPos() (lexer.Position, lexer.Position) {
	return lexer.Position{}, lexer.Position{}
}
func (v VoidExpr) _expression() {}

type AssignmentExpr struct {
	BaseStmt
	Assigne  IdentifierExpr
	Value    Expression
	Operator lexer.Token
}

func (a AssignmentExpr) node() {}
func (a AssignmentExpr) GetPos() (lexer.Position, lexer.Position) {
	return a.StartPos, a.EndPos
}
func (a AssignmentExpr) _expression() {}

type FunctionCallExpr struct {
	BaseStmt
	Function Expression
	Args     []Expression
}

func (c FunctionCallExpr) node() {}
func (c FunctionCallExpr) GetPos() (lexer.Position, lexer.Position) {
	return c.StartPos, c.EndPos
}
func (c FunctionCallExpr) _expression() {}

type StructInstantiationExpr struct {
	BaseStmt
	StructName string
	Properties map[string]Expression
	Methods    map[string]FunctionDeclStmt
}

func (s StructInstantiationExpr) node() {}
func (s StructInstantiationExpr) GetPos() (lexer.Position, lexer.Position) {
	return s.StartPos, s.EndPos
}
func (s StructInstantiationExpr) _expression() {}

type StructPropertyExpr struct {
	BaseStmt
	Object	 		Expression
	Property	 	IdentifierExpr
}
func (s StructPropertyExpr) node() {}
func (s StructPropertyExpr) GetPos() (lexer.Position, lexer.Position) {
	return s.StartPos, s.EndPos
}
func (s StructPropertyExpr) _expression() {}

type ArrayLiterals struct {
	BaseStmt
	Size     uint64
	Elements []Expression
}

func (a ArrayLiterals) node() {}
func (a ArrayLiterals) GetPos() (lexer.Position, lexer.Position) {
	return a.StartPos, a.EndPos
}
func (a ArrayLiterals) _expression() {}
