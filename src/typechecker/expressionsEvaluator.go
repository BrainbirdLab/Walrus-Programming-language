package typechecker

import (
	"fmt"
	"walrus/frontend/ast"
	"walrus/frontend/lexer"
	"walrus/frontend/parser"
	"walrus/helpers"
)

func EvaluateIdenitifierExpr(expr ast.IdentifierExpr, env *Environment) RuntimeValue {
	if !env.HasVariable(expr.Identifier) {

		msg := fmt.Sprintf("variable %v is not declared in this scope\n", expr.Identifier)

		parser.MakeError(env.parser, expr.StartPos.Line, env.parser.FilePath, expr.StartPos, expr.EndPos, msg).Display()
	}

	runtimeVal, err := env.GetRuntimeValue(expr.Identifier)

	if err != nil {
		parser.MakeError(env.parser, expr.StartPos.Line, env.parser.FilePath, expr.StartPos, expr.EndPos, err.Error()).Display()
	}

	return runtimeVal
}

func EvaluateUnaryExpression(unary ast.UnaryExpr, env *Environment) RuntimeValue {
	// Evaluate the unary argument expression
	expr := Evaluate(unary.Argument, env)

	// Error message for unsupported unary operations
	errMsg := fmt.Errorf("unsupported unary operation for type %v", expr)

	// Switch based on the unary operator value
	switch unary.Operator.Value {
	case "-", "+":
		return handleUnaryAdditive(expr, unary, errMsg)
	case "!":
		// Handle unary logical NOT operator
		return handleUnaryNegation(expr, errMsg)

	case "++", "--":
		// Handle pre-increment and pre-decrement operators
		if !helpers.TypesMatchT[IntegerValue](expr) {
			panic(errMsg)
		}

		value := expr.(IntegerValue).Value
		if unary.Operator.Value == "++" {
			value++
		} else {
			value--
		}

		runtimeVal := MakeINT(value, 32, true)

		// Assign the incremented or decremented value if the argument is an identifier
		if idExpr, ok := unary.Argument.(ast.IdentifierExpr); ok {
			if env.HasVariable(idExpr.Identifier) {
				env.AssignVariable(idExpr.Identifier, runtimeVal)
			}
		}

		return runtimeVal

	default:
		// Default case for unsupported unary operators
		return MakeNULL()
	}
}

func handleUnaryNegation(expr RuntimeValue, errMsg error) RuntimeValue {
	if !helpers.TypesMatchT[BooleanValue](expr) {
		panic(errMsg)
	}

	return BooleanValue{
		Type:  ast.BoolType{Kind: ast.T_BOOLEAN},
		Value: !expr.(BooleanValue).Value,
	}
}

func handleUnaryAdditive(expr RuntimeValue, unary ast.UnaryExpr, errMsg error) RuntimeValue {
	// Handle unary minus and plus operators
	if !helpers.TypesMatchT[IntegerValue](expr) {
		panic(errMsg)
	}

	value := expr.(IntegerValue).Value
	if unary.Operator.Value == "-" {
		value = -value
	}

	return MakeINT(value, 32, true)
}

func EvaluateBinaryExpr(binop ast.BinaryExpr, env *Environment) RuntimeValue {

	left := Evaluate(binop.Left, env)
	right := Evaluate(binop.Right, env)

	switch binop.Operator.Value {
	// Arithmetic operators
	case "+", "-", "*", "/", "%", "^":
		return handleBinaryArithmeticExpr(left, right, binop, env)
	// Relational operators
	case "==", "!=", ">", "<", ">=", "<=":
		result, err := evaluateComparisonExpr(left, right, binop.Operator)
		if err != nil {
			handleBinaryExprError(err, binop, env)
		}
		return result

	// Logical operators
	case "&&":
		if IsTruthy(left) {
			return right
		}
		return left

	case "||":
		if IsTruthy(left) {
			return left
		}
		return right

	default:
		handleBinaryExprError(fmt.Errorf("unsupported operator: %v", binop.Operator.Value), binop, env)
	}

	return MakeNULL()
}

func handleBinaryArithmeticExpr(left RuntimeValue, right RuntimeValue, binop ast.BinaryExpr, env *Environment) RuntimeValue {

	leftType := GetRuntimeType(left)
	rightType := GetRuntimeType(right)

	if IsNumber(left) && IsNumber(right) {
		result, err := evaluateNumericArithmeticExpr(left, right, binop.Operator)
		if err != nil {
			handleBinaryExprError(err, binop, env)
		}
		return result
	} else if IsString(left) {
		strVal, err := CastToStringValue(right)
		if err != nil {
			handleBinaryExprError(err, binop, env)
		}
		result, err := evaluateStringExpr(left.(StringValue), strVal, binop.Operator)
		if err != nil {
			handleBinaryExprError(err, binop, env)
		}
		return result
	} else {
		handleBinaryExprError(fmt.Errorf("operand types mismatch: %v and %v", leftType, rightType), binop, env)
	}

	return MakeNULL()
}

func handleBinaryExprError(err error, binop ast.BinaryExpr, env *Environment) {
	parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.Operator.StartPos, binop.Operator.EndPos, err.Error()).Display()
}

func evaluateNumericArithmeticExpr(left RuntimeValue, right RuntimeValue, operator lexer.Token) (RuntimeValue, error) {
	if IsINT(left) {
		if IsINT(right) {
			return evaluateIntInt(left.(IntegerValue), right.(IntegerValue), operator)
		} else {
			return evaluateIntFloat(left.(IntegerValue), right.(FloatValue), operator)
		}
	} else {
		if IsINT(right) {
			return evaluateFloatInt(left.(FloatValue), right.(IntegerValue), operator)
		} else {
			return evaluateFloatFloat(left.(FloatValue), right.(FloatValue), operator)
		}
	}
}

func EvaluateAssignmentExpr(assignNode ast.AssignmentExpr, env *Environment) RuntimeValue {

	var err error

	var variableToAssign ast.IdentifierExpr
	var variableNameString string

	switch assignNode.Assigne.(type) {
	case ast.IdentifierExpr:
		variableToAssign = assignNode.Assigne.(ast.IdentifierExpr)
		variableNameString = variableToAssign.Identifier
	case ast.StructPropertyExpr:
		variableToAssign = assignNode.Assigne.(ast.StructPropertyExpr).Property

		//check if object is an identifier
		object, ok := assignNode.Assigne.(ast.StructPropertyExpr).Object.(ast.IdentifierExpr)
		if !ok {
			err = fmt.Errorf("invalid left-hand side in assignment expression. expected identifier got %s", assignNode.Assigne.(ast.StructPropertyExpr).Object.INodeType())
			parser.MakeError(env.parser, assignNode.StartPos.Line, env.parser.FilePath, variableToAssign.StartPos, variableToAssign.EndPos, err.Error()).Display()
		}

		variableNameString = object.Identifier

	default:
		err = fmt.Errorf("invalid left-hand side in assignment expression")
		parser.MakeError(env.parser, assignNode.StartPos.Line, env.parser.FilePath, variableToAssign.StartPos, variableToAssign.EndPos, err.Error()).Display()
	}

	//if assigne is any of "false", "true", "null";
	if helpers.ContainsIn([]string{"false", "true", "null"}, variableToAssign.Identifier) {
		err = fmt.Errorf("cannot assign to built-in constant %v", variableToAssign.Identifier)
		parser.MakeError(env.parser, assignNode.StartPos.Line, env.parser.FilePath, variableToAssign.StartPos, variableToAssign.EndPos, err.Error()).Display()
	}

	currentValueOfIdentifier, err := env.GetRuntimeValue(variableNameString)

	if err != nil {
		valStart, valEnd := assignNode.Value.GetPos()
		parser.MakeError(env.parser, assignNode.StartPos.Line, env.parser.FilePath, valStart, valEnd, err.Error()).Display()
	}

	valueToSet := Evaluate(assignNode.Value, env)

	if helpers.TypesMatchT[StructInstance](currentValueOfIdentifier) {

		// assign struct instance properties
		// if object is declared in the current scope
		if !env.HasVariable(variableNameString) {
			err = fmt.Errorf("variable %v is not declared in this scope", variableNameString)
			parser.MakeError(env.parser, assignNode.StartPos.Line, env.parser.FilePath, variableToAssign.StartPos, variableToAssign.EndPos, err.Error()).Display()
		}

		structInstance := env.variables[variableNameString]
		i := structInstance.(StructInstance)

		i.Fields[assignNode.Assigne.(ast.StructPropertyExpr).Property.Identifier] = valueToSet

		env.structs[variableNameString] = i

		return valueToSet
	}

	switch assignNode.Operator.Kind {
	case lexer.PLUS_EQUALS_TOKEN, lexer.MINUS_EQUALS_TOKEN, lexer.TIMES_EQUALS_TOKEN, lexer.DIVIDE_EQUALS_TOKEN, lexer.MODULO_EQUALS_TOKEN, lexer.POWER_EQUALS_TOKEN:

		//remove the = from the operator
		opChar := assignNode.Operator.Value[:len(assignNode.Operator.Value)-1]

		operationValue := EvaluateBinaryExpr(ast.BinaryExpr{
			BaseStmt: ast.BaseStmt{
				Kind:     ast.BINARY_EXPRESSION,
				StartPos: assignNode.StartPos,
				EndPos:   assignNode.EndPos,
			},
			Left:  assignNode.Assigne,
			Right: assignNode.Value,
			Operator: lexer.Token{
				Kind:     lexer.TOKEN_KIND(opChar),
				Value:    opChar,
				StartPos: assignNode.Operator.StartPos,
				EndPos:   assignNode.Operator.EndPos,
			},
		}, env)

		valueToSet = operationValue
	}

	runtimeVal, err := env.AssignVariable(variableToAssign.Identifier, valueToSet)

	if err != nil {
		start, end := assignNode.Value.GetPos()
		parser.MakeError(env.parser, assignNode.StartPos.Line, env.parser.FilePath, start, end, err.Error()).Display()
	}

	return runtimeVal
}

func evaluateIntInt(left IntegerValue, right IntegerValue, operator lexer.Token) (RuntimeValue, error) {

	highestBit := uint8(0)

	if left.Size > right.Size {
		highestBit = left.Size
	} else {
		highestBit = right.Size
	}

	switch operator.Value {
	case "+", "+=":
		return MakeINT(left.Value+right.Value, highestBit, true), nil
	case "-", "-=":
		return MakeINT(left.Value-right.Value, highestBit, true), nil
	case "*", "*=":
		return MakeINT(left.Value*right.Value, highestBit, true), nil
	case "/", "/=":
		if right.Value == 0 {
			return nil, errorDivisionByZero
		}
		return MakeINT(left.Value/right.Value, highestBit, true), nil
	case "%", "%=":
		if right.Value == 0 {
			return nil, errorDivisionByZero
		}
		return MakeINT(left.Value%right.Value, highestBit, true), nil
	case "^":
		//power operation
		number := left.Value
		power := right.Value

		//use bit shifting to calculate power
		result := int64(1)
		for power > 0 {
			if power&1 == 1 {
				result *= number
			}
			number *= number
			power >>= 1
		}

		return MakeINT(result, highestBit, true), nil

	default:
		return nil, fmt.Errorf(invalidOperationMsg, operator.Value)
	}
}

func evaluateIntFloat(left IntegerValue, right FloatValue, operator lexer.Token) (RuntimeValue, error) {

	highestBit := uint8(0)

	if left.Size > right.Size {
		highestBit = left.Size
	} else {
		highestBit = right.Size
	}

	switch operator.Value {
	case "+", "+=":
		return MakeINT(int64(float64(left.Value)+right.Value), highestBit, true), nil
	case "-", "-=":
		return MakeINT(int64(float64(left.Value)-right.Value), highestBit, true), nil
	case "*", "*=":
		return MakeINT(int64(float64(left.Value)*right.Value), highestBit, true), nil
	case "/", "/=":
		if right.Value == 0 {
			return nil, errorDivisionByZero
		}
		return MakeINT(int64(float64(left.Value)/right.Value), highestBit, true), nil
	case "^":
		//power operation
		number := float64(left.Value)
		power := right.Value

		//use bit shifting to calculate power
		result := 1.0
		for power > 0 {
			if int64(power)&1 == 1 {
				result *= number
			}
			number *= number
			power /= 2
		}

		return MakeINT(int64(result), highestBit, true), nil
	default:
		return nil, fmt.Errorf(invalidOperationMsg, operator.Value)
	}
}

func evaluateFloatInt(left FloatValue, right IntegerValue, operator lexer.Token) (RuntimeValue, error) {

	highestBit := uint8(0)

	if left.Size > right.Size {
		highestBit = left.Size
	} else {
		highestBit = right.Size
	}

	switch operator.Value {
	case "+", "+=":
		return MakeFLOAT(left.Value+float64(right.Value), highestBit), nil
	case "-", "-=":
		return MakeFLOAT(left.Value-float64(right.Value), highestBit), nil
	case "*", "*=":
		return MakeFLOAT(left.Value*float64(right.Value), highestBit), nil
	case "/", "/=":
		if right.Value == 0 {
			return nil, fmt.Errorf("division by zero is forbidden")
		}
		return MakeFLOAT(left.Value/float64(right.Value), highestBit), nil
	case "^":
		//power operation
		number := left.Value
		power := float64(right.Value)

		//use bit shifting to calculate power
		result := 1.0
		for power > 0 {
			if int64(power)&1 == 1 {
				result *= number
			}
			number *= number
			power /= 2
		}

		return MakeFLOAT(result, highestBit), nil
	default:
		return nil, fmt.Errorf(invalidOperationMsg, operator.Value)
	}
}

func evaluateFloatFloat(left FloatValue, right FloatValue, operator lexer.Token) (RuntimeValue, error) {

	highestBit := uint8(0)

	if left.Size > right.Size {
		highestBit = left.Size
	} else {
		highestBit = right.Size
	}

	switch operator.Value {
	case "+", "+=":
		return MakeFLOAT(left.Value+right.Value, highestBit), nil
	case "-", "-=":
		return MakeFLOAT(left.Value-right.Value, highestBit), nil
	case "*", "*=":
		return MakeFLOAT(left.Value*right.Value, highestBit), nil
	case "/", "/=":
		if right.Value == 0 {
			return nil, fmt.Errorf("division by zero is forbidden")
		}
		return MakeFLOAT(left.Value/right.Value, highestBit), nil
	case "^":
		//power operation
		number := left.Value
		power := right.Value

		//use bit shifting to calculate power
		result := 1.0
		for power > 0 {
			if int64(power)&1 == 1 {
				result *= number
			}
			number *= number
			power /= 2
		}

		return MakeFLOAT(result, highestBit), nil
	default:
		return nil, fmt.Errorf(invalidOperationMsg, operator.Value)
	}
}

func evaluateComparisonExpr(left RuntimeValue, right RuntimeValue, operator lexer.Token) (RuntimeValue, error) {
	// Handle string comparison
	if GetRuntimeType(left) == ast.T_STRING && GetRuntimeType(right) == ast.T_STRING {
		switch operator.Value {
		case "==":
			return MakeBOOL(left.(StringValue).Value == right.(StringValue).Value), nil
		case "!=":
			return MakeBOOL(left.(StringValue).Value != right.(StringValue).Value), nil
		default:
			return nil, fmt.Errorf("operator %v is not supported for string comparison", operator.Value)
		}
	}

	// Handle numeric comparison
	leftValue, err := GetNumericValue(left)
	if err != nil {
		return nil, err
	}

	rightValue, err := GetNumericValue(right)
	if err != nil {
		return nil, err
	}

	switch operator.Value {
	case ">":
		return MakeBOOL(leftValue > rightValue), nil
	case "<":
		return MakeBOOL(leftValue < rightValue), nil
	case ">=":
		return MakeBOOL(leftValue >= rightValue), nil
	case "<=":
		return MakeBOOL(leftValue <= rightValue), nil
	case "==":
		return MakeBOOL(leftValue == rightValue), nil
	case "!=":
		return MakeBOOL(leftValue != rightValue), nil
	case "&&":
		// Logical AND operator
		return evaluateLogicalAND(leftValue, rightValue), nil
	case "||":
		// Logical OR operator
		return evaluateLogicalOR(leftValue, rightValue), nil
	default:
		return nil, fmt.Errorf("unsupported operator %v for comparison", operator.Value)
	}
}

func evaluateLogicalAND(left float64, right float64) RuntimeValue {
	// Logical AND evaluates to the last truthy value if all are true, otherwise false
	if left == 1 && right == 1 {
		return MakeBOOL(true)
	}
	return MakeBOOL(false)
}

func evaluateLogicalOR(left float64, right float64) RuntimeValue {
	// Logical OR evaluates to the first truthy value, otherwise false
	if left == 1 || right == 1 {
		return MakeBOOL(true)
	}
	return MakeBOOL(false)
}

func evaluateStringExpr(left StringValue, right StringValue, operator lexer.Token) (RuntimeValue, error) {

	result := false

	switch operator.Value {
	case "+", "+=":
		return evaluateStringConcat(left, right)
	case "==":
		result = left.Value == right.Value
	case "!=":
		result = left.Value != right.Value
	default:
		return nil, fmt.Errorf("cannot evaluate string operation. unsupported operator %v", operator)
	}

	return MakeBOOL(result), nil
}

func evaluateStringConcat(left StringValue, right StringValue) (RuntimeValue, error) {
	return MakeSTRING(left.Value + right.Value), nil
}
