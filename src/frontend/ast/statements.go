package ast

import "rexlang/frontend/lexer"

type ModuleStmt struct {
	Kind       NODE_TYPE
	ModuleName string
	StartPos   lexer.Position
	EndPos     lexer.Position
}
func (m ModuleStmt) stmt() {} // implements the Stmt interface
func (m ModuleStmt) GetPos() (lexer.Position, lexer.Position) {
	return m.StartPos, m.EndPos
}

type ImportStmt struct {
	Kind 		NODE_TYPE
	ModuleName 	string
	Symbols 	[]string
	StartPos 	lexer.Position
	EndPos   	lexer.Position
}
func (i ImportStmt) stmt() {} // implements the Stmt interface
func (i ImportStmt) GetPos() (lexer.Position, lexer.Position) {
	return i.StartPos, i.EndPos
}

type ExpressionStmt struct {
	Kind       NODE_TYPE
	Expression Expr
	StartPos   lexer.Position
	EndPos     lexer.Position
}
func (e ExpressionStmt) stmt() {} // implements the Stmt interface
func (e ExpressionStmt) GetPos() (lexer.Position, lexer.Position) {
	return e.StartPos, e.EndPos
}

type ProgramStmt 	struct {
	FileName   		string
	ModuleName 		string
	Imports  		[]ImportStmt
	Contents 		[]Stmt
	StartPos 		lexer.Position
	EndPos   		lexer.Position
}
func (p ProgramStmt) stmt() {} // implements the Stmt interface
func (p ProgramStmt) GetPos() (lexer.Position, lexer.Position) {
	return p.StartPos, p.EndPos
}

type BlockStmt struct {
	Kind 		NODE_TYPE
	Body 		[]Stmt
	StartPos 	lexer.Position
	EndPos   	lexer.Position
}
func (b BlockStmt) stmt() {}
func (b BlockStmt) GetPos() (lexer.Position, lexer.Position) {
	return b.StartPos, b.EndPos
}

type VariableDclStml struct {
	Kind         	NODE_TYPE
	IsConstant   	bool
	Identifier   	string
	Value        	Expr
	ExplicitType 	Type
	StartPos 		lexer.Position
	EndPos   		lexer.Position
}
func (v VariableDclStml) stmt() {}
func (v VariableDclStml) GetPos() (lexer.Position, lexer.Position) {
	return v.StartPos, v.EndPos
}

type FunctionDeclStmt struct {
	Kind         	NODE_TYPE
	FunctionName 	string
	Parameters   	map[string]Type
	ReturnType   	Type
	Block        	BlockStmt
	StartPos 		lexer.Position
	EndPos   		lexer.Position
}
func (f FunctionDeclStmt) stmt() {}
func (f FunctionDeclStmt) GetPos() (lexer.Position, lexer.Position) {
	return f.StartPos, f.EndPos
}

type ReturnStmt struct {
	Kind       	NODE_TYPE
	Expression 	Expr
	StartPos 	lexer.Position
	EndPos   	lexer.Position
}
func (r ReturnStmt) stmt() {}
func (r ReturnStmt) GetPos() (lexer.Position, lexer.Position) {
	return r.StartPos, r.EndPos
}

type StructProperty struct {
	IsStatic bool
	IsPublic bool
	ReadOnly bool
	Type     Type
	StartPos lexer.Position
	EndPos   lexer.Position
}

type StructMethod struct {
	IsStatic   	bool
	IsPublic   	bool
	Parameters 	map[string]Type
	ReturnType 	Type
	StartPos 	lexer.Position
	EndPos   	lexer.Position
}

type StructDeclStatement struct {
	Kind       	NODE_TYPE
	StructName 	string
	Properties 	map[string]StructProperty
	Methods    	map[string]StructMethod
	StartPos 	lexer.Position
	EndPos   	lexer.Position
}
func (s StructDeclStatement) stmt() {}
func (s StructDeclStatement) GetPos() (lexer.Position, lexer.Position) {
	return s.StartPos, s.EndPos
}

type IfStmt struct {
	Kind      	NODE_TYPE
	Condition 	Expr
	Block     	BlockStmt
	Alternate 	interface{}
	StartPos 	lexer.Position
	EndPos   	lexer.Position
}
func (i IfStmt) stmt() {}
func (i IfStmt) GetPos() (lexer.Position, lexer.Position) {
	return i.StartPos, i.EndPos
}
