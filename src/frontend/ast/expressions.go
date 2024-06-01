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

func (b BinaryExpr) node() {
	// empty method implements the Node interface
}
func (b BinaryExpr) GetPos() (lexer.Position, lexer.Position) {
	return b.StartPos, b.EndPos
}
func (b BinaryExpr) _expression() {
	// empty method implements the Expression interface
}

type UnaryExpr struct {
	BaseStmt
	Operator lexer.Token
	Argument Expression
}

func (u UnaryExpr) node() {
	// empty method implements the Node interface
}
func (u UnaryExpr) GetPos() (lexer.Position, lexer.Position) {
	return u.StartPos, u.EndPos
}
func (u UnaryExpr) _expression() {
	// empty method implements the Expression interface
}

type IdentifierExpr struct {
	BaseStmt
	Identifier string
	Type       string
}

func (i IdentifierExpr) node() {
	// empty method implements the Node interface
}
func (i IdentifierExpr) GetPos() (lexer.Position, lexer.Position) {
	return i.StartPos, i.EndPos
}
func (i IdentifierExpr) _expression() {
	// empty method implements the Expression interface
}

type NumericLiteral struct {
	BaseStmt
	Value float64
	Type  string
}

func (n NumericLiteral) node() {
	// empty method implements the Node interface
}
func (n NumericLiteral) GetPos() (lexer.Position, lexer.Position) {
	return n.StartPos, n.EndPos
}
func (n NumericLiteral) _expression() {
	// empty method implements the Expression interface
}

type StringLiteral struct {
	BaseStmt
	Value string
	Type  string
}

func (s StringLiteral) node() {
	// empty method implements the Node interface
}
func (s StringLiteral) GetPos() (lexer.Position, lexer.Position) {
	return s.StartPos, s.EndPos
}
func (s StringLiteral) _expression() {
	// empty method implements the Expression interface
}

type BooleanLiteral struct {
	BaseStmt
	Value bool
	Type  string
}

func (b BooleanLiteral) node() {
	// empty method implements the Node interface
}
func (b BooleanLiteral) GetPos() (lexer.Position, lexer.Position) {
	return b.StartPos, b.EndPos
}
func (b BooleanLiteral) _expression() {
	// empty method implements the Expression interface
}

type NullLiteral struct {
	BaseStmt
	Value string
	Type  string
}

func (n NullLiteral) node() {
	// empty method implements the Node interface
}
func (n NullLiteral) GetPos() (lexer.Position, lexer.Position) {
	return n.StartPos, n.EndPos
}
func (n NullLiteral) _expression() {
	// empty method implements the Expression interface
}

type VoidExpr struct {
	Kind  NODE_TYPE
	Value string
	Type  string
}

func (v VoidExpr) node() {
	// empty method implements the Node interface
}
func (v VoidExpr) GetPos() (lexer.Position, lexer.Position) {
	return lexer.Position{}, lexer.Position{}
}
func (v VoidExpr) _expression() {
	// empty method implements the Expression interface
}

type AssignmentExpr struct {
	BaseStmt
	Assigne  IdentifierExpr
	Value    Expression
	Operator lexer.Token
}

func (a AssignmentExpr) node() {
	// empty method implements the Node interface
}
func (a AssignmentExpr) GetPos() (lexer.Position, lexer.Position) {
	return a.StartPos, a.EndPos
}
func (a AssignmentExpr) _expression() {
	// empty method implements the Expression interface
}

type FunctionCallExpr struct {
	BaseStmt
	Function Expression
	Args     []Expression
}

func (c FunctionCallExpr) node() {
	// empty method implements the Node interface
}
func (c FunctionCallExpr) GetPos() (lexer.Position, lexer.Position) {
	return c.StartPos, c.EndPos
}
func (c FunctionCallExpr) _expression() {
	// empty method implements the Expression interface
}

type StructInstantiationExpr struct {
	BaseStmt
	StructName string
	Properties map[string]Expression
	Methods    map[string]FunctionDeclStmt
}

func (s StructInstantiationExpr) node() {
	// empty method implements the Node interface
}
func (s StructInstantiationExpr) GetPos() (lexer.Position, lexer.Position) {
	return s.StartPos, s.EndPos
}
func (s StructInstantiationExpr) _expression() {
	// empty method implements the Expression interface
}

type StructPropertyExpr struct {
	BaseStmt
	Object   Expression
	Property IdentifierExpr
}

func (s StructPropertyExpr) node() {
	// empty method implements the Node interface
}
func (s StructPropertyExpr) GetPos() (lexer.Position, lexer.Position) {
	return s.StartPos, s.EndPos
}
func (s StructPropertyExpr) _expression() {
	// empty method implements the Expression interface
}

type ArrayLiterals struct {
	BaseStmt
	Size     uint64
	Elements []Expression
}

func (a ArrayLiterals) node() {
	// empty method implements the Node interface
}
func (a ArrayLiterals) GetPos() (lexer.Position, lexer.Position) {
	return a.StartPos, a.EndPos
}
func (a ArrayLiterals) _expression() {
	// empty method implements the Expression interface
}
