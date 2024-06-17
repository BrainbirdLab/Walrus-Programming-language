package typechecker

import (
	"fmt"
	"walrus/frontend/ast"
	"walrus/frontend/parser"
	"walrus/frontend/lexer"
	"walrus/helpers"
)

func EvaluateIdenitifierExpr(expr ast.IdentifierExpr, env *Environment) RuntimeValue {
	if !env.HasVariable(expr.Identifier) {

		msg := fmt.Sprintf("Variable %v is not declared in this scope\n", expr.Identifier)

		parser.MakeError(env.parser, expr.StartPos.Line, env.parser.FilePath, expr.StartPos, expr.EndPos, msg).Display()
	}
	return env.GetRuntimeValue(expr.Identifier)
}

func EvaluateUnaryExpression(unary ast.UnaryExpr, env *Environment) RuntimeValue {
	expr := Evaluate(unary.Argument, env)

	switch unary.Operator.Value {
	case "-", "+":
		if !helpers.TypesMatchT[IntegerValue](expr) {
			panic(fmt.Sprintf("Invalid unary operation for type %v", expr))
		}

		var value int

		if unary.Operator.Value == "-" {
			value = -expr.(IntegerValue).Value
		} else {
			value = expr.(IntegerValue).Value
		}

		return IntegerValue{
			Type:  "int",
			Value: value,
		}

	case "!":
		if !helpers.TypesMatchT[BooleanValue](expr) {
			panic(fmt.Sprintf("Invalid unary operation for type %v", expr))
		}

		return BooleanValue{
			Type:  "bool",
			Value: !expr.(BooleanValue).Value,
		}

	case "++", "--":
		if !helpers.TypesMatchT[IntegerValue](expr) {
			panic(fmt.Sprintf("Invalid unary operation for type %v", expr))
		}

		var value int

		if unary.Operator.Value == "++" {
			value = expr.(IntegerValue).Value + 1
		} else {
			value = expr.(IntegerValue).Value - 1
		}

		runtimeVal := IntegerValue{
			Type:  "int",
			Value: value,
		}

		if helpers.TypesMatchT[ast.IdentifierExpr](unary.Argument) {
			if env.HasVariable(unary.Argument.(ast.IdentifierExpr).Identifier) {
				env.AssignVariable(unary.Argument.(ast.IdentifierExpr).Identifier, runtimeVal)
			}
		}

		return runtimeVal

	default:
		return MAKE_NULL()
	}
}

func EvaluateBinaryExpr(binop ast.BinaryExpr, env *Environment) RuntimeValue {

	left := Evaluate(binop.Left, env)
	right := Evaluate(binop.Right, env)

	switch binop.Operator.Value {
	case "+", "-", "*", "/":
		if helpers.TypesMatchT[IntegerValue](left) && helpers.TypesMatchT[IntegerValue](right) {
			// Numeric expr
			return evaluateNumericExpr(left.(IntegerValue), right.(IntegerValue), binop.Operator.Value)
		} else if helpers.ContainsIn([]string{"==", "!="}, binop.Operator.Value) {
			// eval string expr
			return evaluateStringExpr(left.(StringValue), right.(StringValue), binop.Operator.Value)
		} else if binop.Operator.Value == "+" && helpers.TypesMatchT[StringValue](left) && helpers.TypesMatchT[StringValue](right) {
			// eval string concat
			return evaluateStringConcat(left.(StringValue), right.(StringValue))
		} else {
			panic(fmt.Sprintf("Unsupported binary operation between %v and %v", left, right))
		}

	case "==", "!=", ">", "<", ">=", "<=":
		if helpers.TypesMatchT[IntegerValue](left) && helpers.TypesMatchT[IntegerValue](right) {
			// Logical expr
			return evaluateLogicalExpr(left.(IntegerValue), right.(IntegerValue), binop.Operator.Value)
		} else if binop.Operator.Value == "==" || binop.Operator.Value == "!=" {
			// bool expr
			return evaluateBoolExpr(left.(BooleanValue), right.(BooleanValue), binop.Operator.Value)
		} else {
			panic(fmt.Sprintf("Unsupported binary operation between %v and %v", left, right))
		}

	case "+=", "-=", "*=", "/=", "%=":
		if !helpers.TypesMatchT[ast.IdentifierExpr](binop.Left) || !helpers.TypesMatchT[IntegerValue](right) {
			panic(fmt.Sprintf("Unsupported binary operation between %v and %v", left, right))
		}

		if env.HasVariable((binop.Left).(ast.IdentifierExpr).Identifier) {
			return env.AssignVariable((binop.Left).(ast.IdentifierExpr).Identifier, evaluateNumericExpr(left.(IntegerValue), right.(IntegerValue), binop.Operator.Value[:1]))
		} else {
			return evaluateNumericExpr(left.(IntegerValue), right.(IntegerValue), binop.Operator.Value[:1])
		}

	case "&&", "||":
		return evaluateLogicalExpr(left.(IntegerValue), right.(IntegerValue), binop.Operator.Value)

	default:
		panic(fmt.Sprintf("Unsupported binary operation between %v and %v", left, right))
	}
}

func EvaluateAssignmentExpr(assignNode ast.AssignmentExpr, env *Environment) RuntimeValue {

	if assignNode.Assigne.Kind != ast.IDENTIFIER {
		panic(fmt.Sprintf("Invalid left-hand side in assignment expression %v", assignNode.Assigne))
	}

	//if assigne is any of "false", "true", "null";
	if helpers.ContainsIn([]string{"false", "true", "null"}, assignNode.Assigne.Identifier) {
		panic(fmt.Sprintf("Cannot assign to built-in constant %v", assignNode.Assigne.Identifier))
	}

	assigneValue := env.GetRuntimeValue(assignNode.Assigne.Identifier)

	value := Evaluate(assignNode.Value, env)

	switch assignNode.Operator.Kind {
	case lexer.PLUS_EQUALS, lexer.MINUS_EQUALS, lexer.TIMES_EQUALS, lexer.DIVIDE_EQUALS, lexer.MODULO_EQUALS:
		if !helpers.TypesMatchT[IntegerValue](assigneValue) || !helpers.TypesMatchT[IntegerValue](value) {
			panic(fmt.Sprintf("Invalid operation between %v and %v", assigneValue, value))
		}

		value = evaluateNumericExpr(assigneValue.(IntegerValue), value.(IntegerValue), assignNode.Operator.Value[:1])
	}

	return env.AssignVariable(assignNode.Assigne.Identifier, value)
}

func evaluateNumericExpr(left IntegerValue, right IntegerValue, operator string) RuntimeValue {
	result := 0

	switch operator {
	case "+":
		result = left.Value + right.Value
	case "-":
		result = left.Value - right.Value
	case "*":
		result = left.Value * right.Value
	case "/":
		if right.Value == 0 {
			panic("Division by zero is forbidden")
		}
		result = left.Value / right.Value
	case "%":
		if right.Value == 0 {
			panic("Division by zero is forbidden")
		}
		result = left.Value % right.Value
	default:
		panic(fmt.Sprintf("Unsupported operator %v", operator))
	}

	return IntegerValue{
		Type:  "int",
		Value: result,
	}
}

func evaluateLogicalExpr(left IntegerValue, right IntegerValue, operator string) RuntimeValue {

	result := false

	switch operator {
	case "==":
		result = left.Value == right.Value
	case "!=":
		result = left.Value != right.Value
	case ">":
		result = left.Value > right.Value
	case "<":
		result = left.Value < right.Value
	case ">=":
		result = left.Value >= right.Value
	case "<=":
		result = left.Value <= right.Value
	case "&&":
		result = left.Value != 0 && right.Value != 0
	case "||":
		result = left.Value != 0 || right.Value != 0
	default:
		panic(fmt.Sprintf("Unsupported operator %v", operator))
	}

	return BooleanValue{
		Type:  "bool",
		Value: result,
	}
}

func evaluateBoolExpr(left RuntimeValue, right RuntimeValue, operator string) RuntimeValue {
	result := false

	sameType := helpers.TypesMatch(left, right)

	switch operator {
	case "==":
		if sameType {
			result = left == right
		} else {
			result = false
		}
	case "!=":
		if sameType {
			result = left != right
		} else {
			result = true
		}
	case "&&":
		if sameType && helpers.TypesMatchT[BooleanValue](left) {
			result = left.(BooleanValue).Value && right.(BooleanValue).Value
		} else {
			result = IsTruthy(left) && IsTruthy(right)
		}
	case "||":
		if sameType && helpers.TypesMatchT[BooleanValue](left) {
			result = left.(BooleanValue).Value || right.(BooleanValue).Value
		} else {
			result = IsTruthy(left) || IsTruthy(right)
		}
	default:
		panic(fmt.Sprintf("Unsupported operator %v", operator))
	}

	return BooleanValue{
		Type:  "bool",
		Value: result,
	}
}

func evaluateStringExpr(left StringValue, right StringValue, operator string) RuntimeValue {
	result := false

	switch operator {
	case "==":
		result = left.Value == right.Value
	case "!=":
		result = left.Value != right.Value
	default:
		panic(fmt.Sprintf("Unsupported operator %v", operator))
	}

	return BooleanValue{
		Type:  "bool",
		Value: result,
	}
}

func evaluateStringConcat(left StringValue, right StringValue) RuntimeValue {
	return StringValue{
		Type:  "string",
		Value: left.Value + right.Value,
	}
}
