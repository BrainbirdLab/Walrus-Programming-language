package ast

type ExpressionStmt struct {
	Kind       NodeType
	Expression Expr
}

func (e ExpressionStmt) stmt() {} // implements the Stmt interface

type ProgramStmt struct {
	Kind NodeType
}

func (p ProgramStmt) stmt() {} // implements the Stmt interface

type BlockStmt struct {
	Kind NodeType
	Body []Stmt
}

func (b BlockStmt) stmt() {}

type VariableDclStml struct {
	Kind         NodeType
	IsConstant   bool
	Identifier   string
	Value        Expr
	ExplicitType Type
}

func (v VariableDclStml) stmt() {}

type StructProperty struct {
	IsStatic bool
	IsPublic bool
	ReadOnly bool
	Type Type
}

type StructMethod struct {
	IsStatic bool
	IsPublic bool
	Parameters map[string]Type
	ReturnType Type
}

type StructDeclStatement struct {
	StructName string
	Properties map[string]StructProperty
	Methods map[string]StructMethod
}
func (s StructDeclStatement) stmt() {}


type IfStmt struct {
	Kind      NodeType
	Condition Expr
	Block     BlockStmt
	Alternate interface{}
}
func (i IfStmt) stmt() {}
