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

func (b BinaryExpr) INodeType() NODE_TYPE {
	return b.Kind
}

func (b BinaryExpr) GetPos() (lexer.Position, lexer.Position) {
	return b.StartPos, b.EndPos
}
func (b BinaryExpr) iExpression() {
	// empty method implements the Expression interface
}

type UnaryExpr struct {
	BaseStmt
	Operator lexer.Token
	Argument Expression
}

func (u UnaryExpr) INodeType() NODE_TYPE {
	return u.Kind
}
func (u UnaryExpr) GetPos() (lexer.Position, lexer.Position) {
	return u.StartPos, u.EndPos
}
func (u UnaryExpr) iExpression() {
	// empty method implements the Expression interface
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
func (i IdentifierExpr) iExpression() {
	// empty method implements the Expression interface
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
func (n NumericLiteral) iExpression() {
	// empty method implements the Expression interface
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
func (s StringLiteral) iExpression() {
	// empty method implements the Expression interface
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
func (c CharacterLiteral) iExpression() {
	// empty method implements the Expression interface
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
func (b BooleanLiteral) iExpression() {
	// empty method implements the Expression interface
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
func (n NullLiteral) iExpression() {
	// empty method implements the Expression interface
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
func (v VoidLiteral) iExpression() {
	// empty method implements the Expression interface
}

type AssignmentExpr struct {
	BaseStmt
	Assigne  Expression
	Value    Expression
	Operator lexer.Token
}

func (a AssignmentExpr) INodeType() NODE_TYPE {
	return a.Kind
}
func (a AssignmentExpr) GetPos() (lexer.Position, lexer.Position) {
	return a.StartPos, a.EndPos
}
func (a AssignmentExpr) iExpression() {
	// empty method implements the Expression interface
}

type FunctionCallExpr struct {
	BaseStmt
	Caller IdentifierExpr
	Args   []Expression
}

func (c FunctionCallExpr) INodeType() NODE_TYPE {
	return c.Kind
}
func (c FunctionCallExpr) GetPos() (lexer.Position, lexer.Position) {
	return c.StartPos, c.EndPos
}
func (c FunctionCallExpr) iExpression() {
	// empty method implements the Expression interface
}

type StructLiteral struct {
	BaseStmt
	StructName string
	Properties map[string]Expression
}

func (s StructLiteral) INodeType() NODE_TYPE {
	return s.Kind
}
func (s StructLiteral) GetPos() (lexer.Position, lexer.Position) {
	return s.StartPos, s.EndPos
}
func (s StructLiteral) iExpression() {
	// empty method implements the Expression interface
}

type StructPropertyExpr struct {
	BaseStmt
	Object   Expression
	Property IdentifierExpr
}

func (s StructPropertyExpr) INodeType() NODE_TYPE {
	return s.Kind
}
func (s StructPropertyExpr) GetPos() (lexer.Position, lexer.Position) {
	return s.StartPos, s.EndPos
}
func (s StructPropertyExpr) iExpression() {
	// empty method implements the Expression interface
}

type ArrayLiterals struct {
	BaseStmt
	Size     uint64
	Elements []Expression
}

func (a ArrayLiterals) INodeType() NODE_TYPE {
	return a.Kind
}
func (a ArrayLiterals) GetPos() (lexer.Position, lexer.Position) {
	return a.StartPos, a.EndPos
}
func (a ArrayLiterals) iExpression() {
	// empty method implements the Expression interface
}
