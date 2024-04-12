package ast

import "rexlang/lexer"

type NumericLiteral struct {
	Value float64
}

func (node NumericLiteral) expr() {

}

type StringLiteral struct {
	Value string
}

func (node StringLiteral) expr() {

}

type Symbol struct {
	Value string
}

func (node Symbol) expr() {

}

//Complex expressions

// 10 + 3 * 5
// Left: 10
// Operator: +
// Right: 3 * 5
type BinaryExpr struct {
	Left     Expr
	Operator lexer.Token
	Right    Expr
}

func (node BinaryExpr) expr() {
	// TODO:
}

type UnaryExpr struct {
	Operator lexer.Token
	Operand  Expr
}

func (node UnaryExpr) expr() {}
