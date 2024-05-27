package evaluator

import (
	"fmt"
	"walrus/frontend/ast"
	"walrus/frontend/typechecker"
	"walrus/helpers"
)

func Evaluate_idenitifierExpr(expr ast.IdentifierExpr, env *typechecker.Environment) typechecker.RuntimeValue {
	return env.GetRuntimeValue(expr.Identifier)
}

func Evalueate_binaryExpr(binop ast.BinaryExpr, env *typechecker.Environment) typechecker.RuntimeValue {

	left := Evaluate(binop.Left, env)
	right := Evaluate(binop.Right, env)

	switch binop.Operator.Value {
	case "+", "-", "*", "/":
		if helpers.TypesMatchT[typechecker.IntegerValue](left) && helpers.TypesMatchT[typechecker.IntegerValue](right) {
			// Numeric expr
		} else if helpers.ContainsIn([]string{"==", "!="}, binop.Operator.Value) {
			// eval boolean expr
		} else {
			panic(fmt.Sprintf("Unsupported binary operation between %v and %v", left, right))
		}
	}

	return typechecker.IntegerValue{}
}

func Evaluate(astNode ast.Node, env *typechecker.Environment) typechecker.RuntimeValue {
	switch node := astNode.(type) {
	case ast.NumericLiteral:
		return typechecker.IntegerValue{
			Type:  "int",
			Value: int(node.Value),
			Size:  64,
		}
	case ast.StringLiteral:
		return typechecker.StringValue{
			Type:  "string",
			Value: node.Value,
		}
	case ast.BooleanLiteral:
		return typechecker.BooleanValue{
			Type:  "bool",
			Value: node.Value,
		}
	case ast.NullLiteral:
		return typechecker.NullValue{
			Type: "null",
		}

	default:
		panic(fmt.Sprintf("This AST node is not implemengted yet. Node: %v", astNode))
	}
}
