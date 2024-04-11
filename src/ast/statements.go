package ast

type BlockStmt struct {
	Body []Stmt
}

func (node BlockStmt) stmt() {}

type ExpressionStmt struct {
	Expression Expr
}

func (node ExpressionStmt) stmt() {}

type VarDeclStmt struct {
	Name string
	IsConstant bool
	Value Expr
	// Explicit type
}

func (node VarDeclStmt) stmt() {}