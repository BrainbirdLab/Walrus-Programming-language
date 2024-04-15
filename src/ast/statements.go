package ast


type ExpressionStmt struct {
	Kind NodeType
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
	Kind       NodeType
	IsConstant bool
	Identifier string
	Value      Expr
}
func (v VariableDclStml) stmt() {}


type IfStmt struct {
	Kind      NodeType
	Condition Expr
	Block     BlockStmt
	Alternate interface{}
}
func (i IfStmt) stmt() {}
