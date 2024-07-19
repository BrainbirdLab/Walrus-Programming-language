package ast

import (
	"walrus/frontend/lexer"
)

type ModuleStmt struct {
	BaseStmt
	ModuleName string
}

func (m ModuleStmt) INodeType() NODE_TYPE {
	return m.BaseStmt.Kind
}
func (m ModuleStmt) GetPos() (lexer.Position, lexer.Position) {
	return m.StartPos, m.EndPos
}

type ImportStmt struct {
	BaseStmt
	ModuleName  string
	Identifiers []string
}

func (i ImportStmt) INodeType() NODE_TYPE {
	return i.Kind
}

// implements the Statement interface
func (i ImportStmt) GetPos() (lexer.Position, lexer.Position) {
	return i.StartPos, i.EndPos
}

type ProgramStmt struct {
	BaseStmt
	FileName   string
	ModuleName string
	Imports    []ImportStmt
	Contents   []Node
}

func (p ProgramStmt) INodeType() NODE_TYPE {
	return p.Kind
}

// implements the Statement interface
func (p ProgramStmt) GetPos() (lexer.Position, lexer.Position) {
	return p.StartPos, p.EndPos
}

type BlockStmt struct {
	BaseStmt
	Items []Node
}

func (b BlockStmt) INodeType() NODE_TYPE {
	return b.Kind
}
func (b BlockStmt) GetPos() (lexer.Position, lexer.Position) {
	return b.StartPos, b.EndPos
}

type VariableDclStml struct {
	BaseStmt
	IsConstant   bool
	Identifier   IdentifierExpr
	Value        Node
	ExplicitType Type
}

func (v VariableDclStml) INodeType() NODE_TYPE {
	return v.Kind
}
func (v VariableDclStml) GetPos() (lexer.Position, lexer.Position) {
	return v.StartPos, v.EndPos
}

type FunctionParameter struct {
	BaseStmt
	IsVariadic bool
	Identifier IdentifierExpr
	Type       Type
	DefaultVal Node
}

type FunctionPrototype struct {
	BaseStmt
	Name       IdentifierExpr
	Parameters []FunctionParameter
	ReturnType Type
}

type FunctionDeclStmt struct {
	BaseStmt
	FunctionPrototype
	Block BlockStmt
}

func (f FunctionDeclStmt) INodeType() NODE_TYPE {
	return f.Kind
}
func (f FunctionDeclStmt) GetPos() (lexer.Position, lexer.Position) {
	return f.StartPos, f.EndPos
}

type ReturnStmt struct {
	BaseStmt
	Expression Node
}

func (r ReturnStmt) INodeType() NODE_TYPE {
	return r.Kind
}
func (r ReturnStmt) GetPos() (lexer.Position, lexer.Position) {
	return r.StartPos, r.EndPos
}

type BreakStmt struct {
	BaseStmt
}

func (b BreakStmt) INodeType() NODE_TYPE {
	return b.Kind
}
func (b BreakStmt) GetPos() (lexer.Position, lexer.Position) {
	return b.StartPos, b.EndPos
}


type ContinueStmt struct {
	BaseStmt
}

func (c ContinueStmt) INodeType() NODE_TYPE {
	return c.Kind
}
func (c ContinueStmt) GetPos() (lexer.Position, lexer.Position) {
	return c.StartPos, c.EndPos
}

type Property struct {
	BaseStmt
	IsStatic bool
	IsPublic bool
	ReadOnly bool
	Name     string
	Type     Type
}

type StructDeclStatement struct {
	BaseStmt
	StructName string
	Properties map[string]Property
	Methods    map[string]FunctionType
	Embeds     []string
}

func (s StructDeclStatement) INodeType() NODE_TYPE {
	return s.Kind
}
func (s StructDeclStatement) GetPos() (lexer.Position, lexer.Position) {
	return s.StartPos, s.EndPos
}


type Method struct {
	BaseStmt
	FunctionType
	IsStatic bool
	IsPublic bool
}
type TraitDeclStatement struct {
	BaseStmt
	TraitName string
	Methods   map[string]Method
}

func (t TraitDeclStatement) INodeType() NODE_TYPE {
	return t.Kind
}
func (t TraitDeclStatement) GetPos() (lexer.Position, lexer.Position) {
	return t.StartPos, t.EndPos
}

type MethodImplementStmt struct {
	BaseStmt
	FunctionDeclStmt
	TypeToImplement string
	IsPublic        bool
	IsStatic        bool
}

type ImplementStatement struct {
	BaseStmt
	Impliments string
	Traits     []string
	Methods    map[string]MethodImplementStmt
}

func (s ImplementStatement) INodeType() NODE_TYPE {
	return s.Kind
}
func (s ImplementStatement) GetPos() (lexer.Position, lexer.Position) {
	return s.StartPos, s.EndPos
}

type IfStmt struct {
	BaseStmt
	Condition Node
	Block     BlockStmt
	Alternate interface{}
}

func (i IfStmt) INodeType() NODE_TYPE {
	return i.Kind
}
func (i IfStmt) GetPos() (lexer.Position, lexer.Position) {
	return i.StartPos, i.EndPos
}

type ForStmt struct {
	BaseStmt
	Variable  string
	Init      Node
	Condition Node
	Post      Node
	Block     BlockStmt
}

func (f ForStmt) INodeType() NODE_TYPE {
	return f.Kind
}
func (f ForStmt) GetPos() (lexer.Position, lexer.Position) {
	return f.StartPos, f.EndPos
}

type ForeachStmt struct {
	BaseStmt
	Variable      string
	IndexVariable string
	Iterable      Node
	WhereClause   Node
	Block         BlockStmt
}

func (f ForeachStmt) INodeType() NODE_TYPE {
	return f.Kind
}
func (f ForeachStmt) GetPos() (lexer.Position, lexer.Position) {
	return f.StartPos, f.EndPos
}

type WhileLoopStmt struct {
	BaseStmt
	Condition Node
	Block     BlockStmt
}

func (w WhileLoopStmt) INodeType() NODE_TYPE {
	return w.Kind
}
func (w WhileLoopStmt) GetPos() (lexer.Position, lexer.Position) {
	return w.StartPos, w.EndPos
}

type SwitchCase struct {
	BaseStmt
	Consequent BlockStmt
	Test       Node
}

type SwitchStmt struct {
	BaseStmt
	Discriminant Node
	Cases        []SwitchCase
}

func (s SwitchStmt) INodeType() NODE_TYPE {
	return s.Kind
}
func (s SwitchStmt) GetPos() (lexer.Position, lexer.Position) {
	return s.StartPos, s.EndPos
}