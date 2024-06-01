package ast

import (
	"walrus/frontend/lexer"
)

type ModuleStmt struct {
	BaseStmt
	ModuleName string
}

func (m ModuleStmt) node() {
	// empty method implements the Node interface
} // implements the Statement interface
func (m ModuleStmt) GetPos() (lexer.Position, lexer.Position) {
	return m.StartPos, m.EndPos
}
func (m ModuleStmt) _statement() {
	// empty method implements the Statement interface
}

type ImportStmt struct {
	BaseStmt
	ModuleName  string
	Identifiers []string
}

func (i ImportStmt) node() {
	// empty method implements the Node interface
} // implements the Statement interface
func (i ImportStmt) GetPos() (lexer.Position, lexer.Position) {
	return i.StartPos, i.EndPos
}
func (i ImportStmt) _statement() {
	// empty method implements the Statement interface
}

type ProgramStmt struct {
	BaseStmt
	FileName   string
	ModuleName string
	Imports    []ImportStmt
	Contents   []Node
}

func (p ProgramStmt) node() {
	// empty method implements the Node interface
} // implements the Statement interface
func (p ProgramStmt) GetPos() (lexer.Position, lexer.Position) {
	return p.StartPos, p.EndPos
}
func (p ProgramStmt) _statement() {
	// empty method implements the Statement interface
}

type BlockStmt struct {
	BaseStmt
	Body []Node
}

func (b BlockStmt) node() {
	// empty method implements the Node interface
}
func (b BlockStmt) GetPos() (lexer.Position, lexer.Position) {
	return b.StartPos, b.EndPos
}
func (b BlockStmt) _statement() {
	// empty method implements the Statement interface
}

type VariableDclStml struct {
	BaseStmt
	IsConstant   bool
	Identifier   string
	Value        Expression
	ExplicitType Type
}

func (v VariableDclStml) node() {
	// empty method implements the Node interface
}
func (v VariableDclStml) GetPos() (lexer.Position, lexer.Position) {
	return v.StartPos, v.EndPos
}
func (v VariableDclStml) _statement() {
	// empty method implements the Statement interface
}

type FunctionPrototype struct {
	BaseStmt
	FunctionName string
	Parameters   map[string]Type
	ReturnType   Type
}

type FunctionDeclStmt struct {
	BaseStmt
	FunctionPrototype
	Block BlockStmt
}

func (f FunctionDeclStmt) node() {
	// empty method implements the Node interface
}
func (f FunctionDeclStmt) GetPos() (lexer.Position, lexer.Position) {
	return f.StartPos, f.EndPos
}
func (f FunctionDeclStmt) _statement() {
	// empty method implements the Statement interface
}

type ReturnStmt struct {
	BaseStmt
	Expression Expression
}

func (r ReturnStmt) node() {
	// empty method implements the Node interface
}
func (r ReturnStmt) GetPos() (lexer.Position, lexer.Position) {
	return r.StartPos, r.EndPos
}
func (r ReturnStmt) _statement() {
	// empty method implements the Statement interface
}

type BreakStmt struct {
	BaseStmt
}

func (b BreakStmt) node() {
	// empty method implements the Node interface
}
func (b BreakStmt) GetPos() (lexer.Position, lexer.Position) {
	return b.StartPos, b.EndPos
}
func (b BreakStmt) _statement() {
	// empty method implements the Statement interface
}

type ContinueStmt struct {
	BaseStmt
}

func (c ContinueStmt) node() {
	// empty method implements the Node interface
}
func (c ContinueStmt) GetPos() (lexer.Position, lexer.Position) {
	return c.StartPos, c.EndPos
}
func (c ContinueStmt) _statement() {
	// empty method implements the Statement interface
}

type Property struct {
	BaseStmt
	IsStatic bool
	IsPublic bool
	ReadOnly bool
	Type     Type
}

type StructDeclStatement struct {
	BaseStmt
	StructName string
	Properties map[string]Property
	Methods    map[string]FunctionPrototype
	Embeds     []string
}

func (s StructDeclStatement) node() {
	// empty method implements the Node interface
}
func (s StructDeclStatement) GetPos() (lexer.Position, lexer.Position) {
	return s.StartPos, s.EndPos
}
func (s StructDeclStatement) _statement() {
	// empty method implements the Statement interface
}

type Method struct {
	BaseStmt
	FunctionPrototype
	IsStatic bool
	IsPublic bool
}
type TraitDeclStatement struct {
	BaseStmt
	TraitName string
	Methods   map[string]Method
}

func (t TraitDeclStatement) node() {
	// empty method implements the Node interface
}
func (t TraitDeclStatement) GetPos() (lexer.Position, lexer.Position) {
	return t.StartPos, t.EndPos
}
func (t TraitDeclStatement) _statement() {
	// empty method implements the Statement interface
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

func (s ImplementStatement) node() {
	// empty method implements the Node interface
}
func (s ImplementStatement) GetPos() (lexer.Position, lexer.Position) {
	return s.StartPos, s.EndPos
}
func (s ImplementStatement) _statement() {
	// empty method implements the Statement interface
}

type IfStmt struct {
	BaseStmt
	Condition Expression
	Block     BlockStmt
	Alternate interface{}
}

func (i IfStmt) node() {
	// empty method implements the Node interface
}
func (i IfStmt) GetPos() (lexer.Position, lexer.Position) {
	return i.StartPos, i.EndPos
}
func (i IfStmt) _statement() {
	// empty method implements the Statement interface
}

type ForStmt struct {
	BaseStmt
	Variable  string
	Init      Expression
	Condition Expression
	Post      Expression
	Block     BlockStmt
}

func (f ForStmt) node() {
	// empty method implements the Node interface
}
func (f ForStmt) GetPos() (lexer.Position, lexer.Position) {
	return f.StartPos, f.EndPos
}
func (f ForStmt) _statement() {
	// empty method implements the Statement interface
}

type ForeachStmt struct {
	BaseStmt
	Variable      string
	IndexVariable string
	Iterable      Expression
	WhereClause   Expression
	Block         BlockStmt
}

func (f ForeachStmt) node() {
	// empty method implements the Node interface
}
func (f ForeachStmt) GetPos() (lexer.Position, lexer.Position) {
	return f.StartPos, f.EndPos
}
func (f ForeachStmt) _statement() {
	// empty method implements the Statement interface
}

type WhileLoopStmt struct {
	BaseStmt
	Condition Expression
	Block     BlockStmt
}

func (w WhileLoopStmt) node() {
	// empty method implements the Node interface
}
func (w WhileLoopStmt) GetPos() (lexer.Position, lexer.Position) {
	return w.StartPos, w.EndPos
}
func (w WhileLoopStmt) _statement() {
	// empty method implements the Statement interface
}

type SwitchCase struct {
	BaseStmt
	Consequent BlockStmt
	Test       Expression
}

type SwitchStmt struct {
	BaseStmt
	Discriminant Expression
	Cases        []SwitchCase
}

func (s SwitchStmt) node() {
	// empty method implements the Node interface
}
func (s SwitchStmt) GetPos() (lexer.Position, lexer.Position) {
	return s.StartPos, s.EndPos
}
func (s SwitchStmt) _statement() {
	// empty method implements the Statement interface
}
