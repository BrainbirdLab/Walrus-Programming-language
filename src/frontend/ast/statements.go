package ast

import "walrus/frontend/lexer"

type ModuleStmt struct {
	BaseStmt
	ModuleName string
}

func (m ModuleStmt) node() {} // implements the Statement interface
func (m ModuleStmt) GetPos() (lexer.Position, lexer.Position) {
	return m.StartPos, m.EndPos
}
func (m ModuleStmt) _statement() {}

type ImportStmt struct {
	BaseStmt
	ModuleName  string
	Identifiers []string
}

func (i ImportStmt) node() {} // implements the Statement interface
func (i ImportStmt) GetPos() (lexer.Position, lexer.Position) {
	return i.StartPos, i.EndPos
}
func (i ImportStmt) _statement() {}

type ProgramStmt struct {
	BaseStmt
	FileName   string
	ModuleName string
	Imports    []ImportStmt
	Contents   []Node
}

func (p ProgramStmt) node() {} // implements the Statement interface
func (p ProgramStmt) GetPos() (lexer.Position, lexer.Position) {
	return p.StartPos, p.EndPos
}
func (p ProgramStmt) _statement() {}

type BlockStmt struct {
	BaseStmt
	Body []Node
}

func (b BlockStmt) node() {}
func (b BlockStmt) GetPos() (lexer.Position, lexer.Position) {
	return b.StartPos, b.EndPos
}
func (b BlockStmt) _statement() {}

type VariableDclStml struct {
	BaseStmt
	IsConstant   bool
	Identifier   string
	Value        Expression
	ExplicitType Type
}

func (v VariableDclStml) node() {}
func (v VariableDclStml) GetPos() (lexer.Position, lexer.Position) {
	return v.StartPos, v.EndPos
}
func (v VariableDclStml) _statement() {}

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

func (f FunctionDeclStmt) node() {}
func (f FunctionDeclStmt) GetPos() (lexer.Position, lexer.Position) {
	return f.StartPos, f.EndPos
}
func (f FunctionDeclStmt) _statement() {}

type ReturnStmt struct {
	BaseStmt
	Expression Expression
}

func (r ReturnStmt) node() {}
func (r ReturnStmt) GetPos() (lexer.Position, lexer.Position) {
	return r.StartPos, r.EndPos
}
func (r ReturnStmt) _statement() {}

type StructProperty struct {
	BaseStmt
	IsStatic bool
	IsPublic bool
	ReadOnly bool
	Type     Type
}

type StructDeclStatement struct {
	BaseStmt
	StructName string
	Properties map[string]StructProperty
	Methods    map[string]FunctionPrototype
	Embeds     []string
}

func (s StructDeclStatement) node() {}
func (s StructDeclStatement) GetPos() (lexer.Position, lexer.Position) {
	return s.StartPos, s.EndPos
}
func (s StructDeclStatement) _statement() {}

type TraitMethod struct {
	BaseStmt
	FunctionPrototype
	IsStatic bool
	IsPublic bool
}
type TraitDeclStatement struct {
	BaseStmt
	TraitName string
	Methods   map[string]TraitMethod
}

func (t TraitDeclStatement) node() {}
func (t TraitDeclStatement) GetPos() (lexer.Position, lexer.Position) {
	return t.StartPos, t.EndPos
}
func (t TraitDeclStatement) _statement() {}

type MethodImplementStmt struct {
	BaseStmt
	FunctionDeclStmt
	StructName string
	IsPublic   bool
	IsStatic   bool
}

type ImplementStatement struct {
	BaseStmt
	Impliments Type
	Methods    map[string]MethodImplementStmt
}

func (s ImplementStatement) node() {}
func (s ImplementStatement) GetPos() (lexer.Position, lexer.Position) {
	return s.StartPos, s.EndPos
}
func (s ImplementStatement) _statement() {}

type IfStmt struct {
	BaseStmt
	Condition Expression
	Block     BlockStmt
	Alternate interface{}
}

func (i IfStmt) node() {}
func (i IfStmt) GetPos() (lexer.Position, lexer.Position) {
	return i.StartPos, i.EndPos
}
func (i IfStmt) _statement() {}

type ForStmt struct {
	BaseStmt
	Variable  string
	Init      Expression
	Condition Expression
	Post      Expression
	Block     BlockStmt
}

func (f ForStmt) node() {}
func (f ForStmt) GetPos() (lexer.Position, lexer.Position) {
	return f.StartPos, f.EndPos
}
func (f ForStmt) _statement() {}

type ForeachStmt struct {
	BaseStmt
	Variable      string
	IndexVariable string
	Iterable      Expression
	WhereClause   Expression
	Block         BlockStmt
}

func (f ForeachStmt) node() {}
func (f ForeachStmt) GetPos() (lexer.Position, lexer.Position) {
	return f.StartPos, f.EndPos
}
func (f ForeachStmt) _statement() {}
