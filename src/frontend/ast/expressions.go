package ast

import (
	"walrus/frontend/lexer"
)

type BinaryExpr struct {
	BaseStmt
	Operator lexer.Token
	Left     Node
	Right    Node
}

func (b BinaryExpr) INodeType() NODE_TYPE {
	return b.Kind
}

func (b BinaryExpr) GetPos() (lexer.Position, lexer.Position) {
	return b.StartPos, b.EndPos
}

type UnaryExpr struct {
	BaseStmt
	Operator lexer.Token
	Argument Node
}

func (u UnaryExpr) INodeType() NODE_TYPE {
	return u.Kind
}
func (u UnaryExpr) GetPos() (lexer.Position, lexer.Position) {
	return u.StartPos, u.EndPos
}

type IdentifierExpr struct {
	BaseStmt
	Identifier string
}

func (i IdentifierExpr) INodeType() NODE_TYPE {
	return i.Kind
}
func (i IdentifierExpr) GetPos() (lexer.Position, lexer.Position) {
	return i.StartPos, i.EndPos
}

type NumericLiteral struct {
	BaseStmt
	Value   string
	BitSize uint8
}

func (n NumericLiteral) INodeType() NODE_TYPE {
	return n.Kind
}
func (n NumericLiteral) GetPos() (lexer.Position, lexer.Position) {
	return n.StartPos, n.EndPos
}

type StringLiteral struct {
	BaseStmt
	Value string
}

func (s StringLiteral) INodeType() NODE_TYPE {
	return s.Kind
}
func (s StringLiteral) GetPos() (lexer.Position, lexer.Position) {
	return s.StartPos, s.EndPos
}

type CharacterLiteral struct {
	BaseStmt
	Value string
}

func (c CharacterLiteral) INodeType() NODE_TYPE {
	return c.Kind
}
func (c CharacterLiteral) GetPos() (lexer.Position, lexer.Position) {
	return c.StartPos, c.EndPos
}

type BooleanLiteral struct {
	BaseStmt
	Value bool
}

func (b BooleanLiteral) INodeType() NODE_TYPE {
	return b.Kind
}
func (b BooleanLiteral) GetPos() (lexer.Position, lexer.Position) {
	return b.StartPos, b.EndPos
}

type NullLiteral struct {
	BaseStmt
	Value string
}

func (n NullLiteral) INodeType() NODE_TYPE {
	return n.Kind
}
func (n NullLiteral) GetPos() (lexer.Position, lexer.Position) {
	return n.StartPos, n.EndPos
}


type VoidLiteral struct {
	Kind NODE_TYPE
}

func (v VoidLiteral) INodeType() NODE_TYPE {
	return v.Kind
}
func (v VoidLiteral) GetPos() (lexer.Position, lexer.Position) {
	return lexer.Position{}, lexer.Position{}
}

type AssignmentExpr struct {
	BaseStmt
	Assigne  Node
	Value    Node
	Operator lexer.Token
}

func (a AssignmentExpr) INodeType() NODE_TYPE {
	return a.Kind
}
func (a AssignmentExpr) GetPos() (lexer.Position, lexer.Position) {
	return a.StartPos, a.EndPos
}

type FunctionCallExpr struct {
	BaseStmt
	Caller IdentifierExpr
	Args   []Node
}

func (c FunctionCallExpr) INodeType() NODE_TYPE {
	return c.Kind
}
func (c FunctionCallExpr) GetPos() (lexer.Position, lexer.Position) {
	return c.StartPos, c.EndPos
}

type StructLiteral struct {
	BaseStmt
	StructName string
	Properties map[string]Node
}

func (s StructLiteral) INodeType() NODE_TYPE {
	return s.Kind
}
func (s StructLiteral) GetPos() (lexer.Position, lexer.Position) {
	return s.StartPos, s.EndPos
}


type StructPropertyExpr struct {
	BaseStmt
	Object   Node
	Property IdentifierExpr
}

func (s StructPropertyExpr) INodeType() NODE_TYPE {
	return s.Kind
}
func (s StructPropertyExpr) GetPos() (lexer.Position, lexer.Position) {
	return s.StartPos, s.EndPos
}

type ArrayLiterals struct {
	BaseStmt
	Size     uint64
	Elements []Node
}

func (a ArrayLiterals) INodeType() NODE_TYPE {
	return a.Kind
}
func (a ArrayLiterals) GetPos() (lexer.Position, lexer.Position) {
	return a.StartPos, a.EndPos
}

type ArrayIndexAccess struct {
	BaseStmt
	ArrayName string
	Index Node
}
func (a ArrayIndexAccess) INodeType() NODE_TYPE {
	return a.Kind
}
func (a ArrayIndexAccess) GetPos() (lexer.Position, lexer.Position) {
	return a.StartPos, a.EndPos
}