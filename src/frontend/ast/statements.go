package ast

import (
	"walrus/frontend/lexer"
)

type ModuleStmt struct {
	BaseStmt
	ModuleName string
}

func (m ModuleStmt) iNode() {
	// empty method implements the Node interface
} // implements the Statement interface
func (m ModuleStmt) GetPos() (lexer.Position, lexer.Position) {
	return m.StartPos, m.EndPos
}
func (m ModuleStmt) iStatement() {
	// empty method implements the Statement interface
}

type ImportStmt struct {
	BaseStmt
	ModuleName  string
	Identifiers []string
}

func (i ImportStmt) iNode() {
	// empty method implements the Node interface
} // implements the Statement interface
func (i ImportStmt) GetPos() (lexer.Position, lexer.Position) {
	return i.StartPos, i.EndPos
}
func (i ImportStmt) iStatement() {
	// empty method implements the Statement interface
}

type ProgramStmt struct {
	BaseStmt
	FileName   string
	ModuleName string
	Imports    []ImportStmt
	Contents   []Node
}

func (p ProgramStmt) iNode() {
	// empty method implements the Node interface
} // implements the Statement interface
func (p ProgramStmt) GetPos() (lexer.Position, lexer.Position) {
	return p.StartPos, p.EndPos
}
func (p ProgramStmt) iStatement() {
	// empty method implements the Statement interface
}

type BlockStmt struct {
	BaseStmt
	Body []Node
}

func (b BlockStmt) iNode() {
	// empty method implements the Node interface
}
func (b BlockStmt) GetPos() (lexer.Position, lexer.Position) {
	return b.StartPos, b.EndPos
}
func (b BlockStmt) iStatement() {
	// empty method implements the Statement interface
}

type VariableDclStml struct {
	BaseStmt
	IsConstant   bool
	Identifier   IdentifierExpr
	Value        Expression
	ExplicitType Type
}

func (v VariableDclStml) iNode() {
	// empty method implements the Node interface
}
func (v VariableDclStml) GetPos() (lexer.Position, lexer.Position) {
	return v.StartPos, v.EndPos
}
func (v VariableDclStml) iStatement() {
	// empty method implements the Statement interface
}

type FunctionParameter struct {
	BaseStmt
	IsVariadic bool
	Identifier IdentifierExpr
	Type       Type
	DefaultVal Expression
}

type FunctionPrototype struct {
	BaseStmt
	Name 	 	IdentifierExpr
	Parameters []FunctionParameter
	ReturnType  Type
}

type FunctionDeclStmt struct {
	BaseStmt
	FunctionPrototype
	Block BlockStmt
}

func (f FunctionDeclStmt) iNode() {
	// empty method implements the Node interface
}
func (f FunctionDeclStmt) GetPos() (lexer.Position, lexer.Position) {
	return f.StartPos, f.EndPos
}
func (f FunctionDeclStmt) iStatement() {
	// empty method implements the Statement interface
}

type ReturnStmt struct {
	BaseStmt
	Expression Expression
}

func (r ReturnStmt) iNode() {
	// empty method implements the Node interface
}
func (r ReturnStmt) GetPos() (lexer.Position, lexer.Position) {
	return r.StartPos, r.EndPos
}
func (r ReturnStmt) iStatement() {
	// empty method implements the Statement interface
}

type BreakStmt struct {
	BaseStmt
}

func (b BreakStmt) iNode() {
	// empty method implements the Node interface
}
func (b BreakStmt) GetPos() (lexer.Position, lexer.Position) {
	return b.StartPos, b.EndPos
}
func (b BreakStmt) iStatement() {
	// empty method implements the Statement interface
}

type ContinueStmt struct {
	BaseStmt
}

func (c ContinueStmt) iNode() {
	// empty method implements the Node interface
}
func (c ContinueStmt) GetPos() (lexer.Position, lexer.Position) {
	return c.StartPos, c.EndPos
}
func (c ContinueStmt) iStatement() {
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
	Methods    map[string]FunctionType
	Embeds     []string
}

func (s StructDeclStatement) iNode() {
	// empty method implements the Node interface
}
func (s StructDeclStatement) GetPos() (lexer.Position, lexer.Position) {
	return s.StartPos, s.EndPos
}
func (s StructDeclStatement) iStatement() {
	// empty method implements the Statement interface
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

func (t TraitDeclStatement) iNode() {
	// empty method implements the Node interface
}
func (t TraitDeclStatement) GetPos() (lexer.Position, lexer.Position) {
	return t.StartPos, t.EndPos
}
func (t TraitDeclStatement) iStatement() {
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

func (s ImplementStatement) iNode() {
	// empty method implements the Node interface
}
func (s ImplementStatement) GetPos() (lexer.Position, lexer.Position) {
	return s.StartPos, s.EndPos
}
func (s ImplementStatement) iStatement() {
	// empty method implements the Statement interface
}

type IfStmt struct {
	BaseStmt
	Condition Expression
	Block     BlockStmt
	Alternate interface{}
}

func (i IfStmt) iNode() {
	// empty method implements the Node interface
}
func (i IfStmt) GetPos() (lexer.Position, lexer.Position) {
	return i.StartPos, i.EndPos
}
func (i IfStmt) iStatement() {
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

func (f ForStmt) iNode() {
	// empty method implements the Node interface
}
func (f ForStmt) GetPos() (lexer.Position, lexer.Position) {
	return f.StartPos, f.EndPos
}
func (f ForStmt) iStatement() {
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

func (f ForeachStmt) iNode() {
	// empty method implements the Node interface
}
func (f ForeachStmt) GetPos() (lexer.Position, lexer.Position) {
	return f.StartPos, f.EndPos
}
func (f ForeachStmt) iStatement() {
	// empty method implements the Statement interface
}

type WhileLoopStmt struct {
	BaseStmt
	Condition Expression
	Block     BlockStmt
}

func (w WhileLoopStmt) iNode() {
	// empty method implements the Node interface
}
func (w WhileLoopStmt) GetPos() (lexer.Position, lexer.Position) {
	return w.StartPos, w.EndPos
}
func (w WhileLoopStmt) iStatement() {
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

func (s SwitchStmt) iNode() {
	// empty method implements the Node interface
}
func (s SwitchStmt) GetPos() (lexer.Position, lexer.Position) {
	return s.StartPos, s.EndPos
}
func (s SwitchStmt) iStatement() {
	// empty method implements the Statement interface
}
