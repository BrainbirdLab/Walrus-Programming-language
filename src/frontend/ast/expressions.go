package ast

import (
	"rexlang/frontend/lexer"
)

type BinaryExpr struct {
	Kind     NODE_TYPE
	Operator lexer.Token
	Left     Expression
	Right    Expression
	StartPos lexer.Position
	EndPos   lexer.Position
}

func (b BinaryExpr) node() {}
func (b BinaryExpr) GetPos() (lexer.Position, lexer.Position) {
	return b.StartPos, b.EndPos
}
func (b BinaryExpr) _expression() {}

type UnaryExpr struct {
	Kind     NODE_TYPE
	Operator lexer.Token
	Argument Expression
	StartPos lexer.Position
	EndPos   lexer.Position
}

func (u UnaryExpr) node() {}
func (u UnaryExpr) GetPos() (lexer.Position, lexer.Position) {
	return u.StartPos, u.EndPos
}
func (u UnaryExpr) _expression() {}

type IdentifierExpr struct {
	Kind     	NODE_TYPE
	Identifier  string
	Type     	string
	StartPos 	lexer.Position
	EndPos   	lexer.Position
}

func (i IdentifierExpr) node() {}
func (i IdentifierExpr) GetPos() (lexer.Position, lexer.Position) {
	return i.StartPos, i.EndPos
}
func (i IdentifierExpr) _expression() {}

type NumericLiteral struct {
	Kind     NODE_TYPE
	Value    float64
	Type     string
	StartPos lexer.Position
	EndPos   lexer.Position
}

func (n NumericLiteral) node() {}
func (n NumericLiteral) GetPos() (lexer.Position, lexer.Position) {
	return n.StartPos, n.EndPos
}
func (n NumericLiteral) _expression() {}

type StringLiteral struct {
	Kind     NODE_TYPE
	Value    string
	Type     string
	StartPos lexer.Position
	EndPos   lexer.Position
}

func (s StringLiteral) node() {}
func (s StringLiteral) GetPos() (lexer.Position, lexer.Position) {
	return s.StartPos, s.EndPos
}
func (s StringLiteral) _expression() {}

type BooleanLiteral struct {
	Kind     NODE_TYPE
	Value    bool
	Type     string
	StartPos lexer.Position
	EndPos   lexer.Position
}

func (b BooleanLiteral) node() {}
func (b BooleanLiteral) GetPos() (lexer.Position, lexer.Position) {
	return b.StartPos, b.EndPos
}
func (b BooleanLiteral) _expression() {}

type NullLiteral struct {
	Kind     NODE_TYPE
	Value    string
	Type     string
	StartPos lexer.Position
	EndPos   lexer.Position
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
	Kind     NODE_TYPE
	Assigne  IdentifierExpr
	Value    Expression
	Operator lexer.Token
	StartPos lexer.Position
	EndPos   lexer.Position
}

func (a AssignmentExpr) node() {}
func (a AssignmentExpr) GetPos() (lexer.Position, lexer.Position) {
	return a.StartPos, a.EndPos
}
func (a AssignmentExpr) _expression() {}

type FunctionCallExpr struct {
	Kind     NODE_TYPE
	Function Expression
	Args     []Expression
	StartPos lexer.Position
	EndPos   lexer.Position
}

func (c FunctionCallExpr) node() {}
func (c FunctionCallExpr) GetPos() (lexer.Position, lexer.Position) {
	return c.StartPos, c.EndPos
}
func (c FunctionCallExpr) _expression() {}

type StructInstantiationExpr struct {
	Kind       NODE_TYPE
	StructName string
	Properties map[string]Expression
	Methods    map[string]FunctionDeclStmt
	StartPos   lexer.Position
	EndPos     lexer.Position
}

func (s StructInstantiationExpr) node() {}
func (s StructInstantiationExpr) GetPos() (lexer.Position, lexer.Position) {
	return s.StartPos, s.EndPos
}
func (s StructInstantiationExpr) _expression() {}

type StructPropertyExpr struct {
	Kind       		NODE_TYPE
	Object	 		Expression
	Property	 	IdentifierExpr
	StartPos   		lexer.Position
	EndPos     		lexer.Position
}
func (s StructPropertyExpr) node() {}
func (s StructPropertyExpr) GetPos() (lexer.Position, lexer.Position) {
	return s.StartPos, s.EndPos
}
func (s StructPropertyExpr) _expression() {}

type ArrayLiterals struct {
	Kind     NODE_TYPE
	Size     uint64
	Elements []Expression
	StartPos lexer.Position
	EndPos   lexer.Position
}

func (a ArrayLiterals) node() {}
func (a ArrayLiterals) GetPos() (lexer.Position, lexer.Position) {
	return a.StartPos, a.EndPos
}
func (a ArrayLiterals) _expression() {}
