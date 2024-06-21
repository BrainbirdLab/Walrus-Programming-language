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

	expr := Evaluate(unary.Argument, env)

	switch unary.Operator.Value {
	case "-", "+":
		if !helpers.TypesMatchT[IntegerValue](expr) {
			panic(fmt.Sprintf("Invalid unary operation for type %v", expr))
		}

		var value int64

		if unary.Operator.Value == "-" {
			value = -expr.(IntegerValue).Value
		} else {
			value = expr.(IntegerValue).Value
		}

		return MAKE_INT(value, 32, true)

	case "!":
		if !helpers.TypesMatchT[BooleanValue](expr) {
			panic(fmt.Sprintf("Invalid unary operation for type %v", expr))
		}

		return BooleanValue{
			Type: ast.BoolType{
				Kind: ast.T_BOOLEAN,
			},
			Value: !expr.(BooleanValue).Value,
		}

	case "++", "--":
		if !helpers.TypesMatchT[IntegerValue](expr) {
			panic(fmt.Sprintf("Invalid unary operation for type %v", expr))
		}

		var value int64

		if unary.Operator.Value == "++" {
			value = expr.(IntegerValue).Value + 1
		} else {
			value = expr.(IntegerValue).Value - 1
		}

		runtimeVal := MAKE_INT(value, 32, true)

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

	leftType := GetRuntimeType(left)
	rightType := GetRuntimeType(right)

	errMsg := fmt.Sprintf("Unsupported binary operation between %v and %v", leftType, rightType)

	switch binop.Operator.Value {
	case "+", "-", "*", "/", "^":
		if IsINT(left) || IsFLOAT(left) && IsINT(right) || IsFLOAT(right) {
			// Numeric expr
			val, err := evaluateNumericExpr(left, right, binop.Operator)

			if err != nil {
				parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.StartPos, binop.EndPos, err.Error()).Display()
			}

			return val

		} else if helpers.ContainsIn([]string{"==", "!="}, binop.Operator.Value) {
			// eval string expr
			val, err := evaluateStringExpr(left.(StringValue), right.(StringValue), binop.Operator)
			if err != nil {
				parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.StartPos, binop.EndPos, err.Error()).Display()
			}

			return val
		} else if binop.Operator.Value == "+" && helpers.TypesMatchT[StringValue](left) && helpers.TypesMatchT[StringValue](right) {
			// eval string concat
			val, err := evaluateStringConcat(left.(StringValue), right.(StringValue))
			if err != nil {
				parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.StartPos, binop.EndPos, err.Error()).Display()
			}

			return val
		} else {

			parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.Operator.StartPos, binop.Operator.EndPos, errMsg).Display()

			return nil
		}

	case "==", "!=", ">", "<", ">=", "<=":
		if helpers.TypesMatchT[IntegerValue](left) && helpers.TypesMatchT[IntegerValue](right) {
			// Logical expr
			val, err := evaluateComparisonExpr(left.(IntegerValue), right.(IntegerValue), binop.Operator)

			if err != nil {
				parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.StartPos, binop.EndPos, err.Error()).Display()
			}

			return val
		} else if binop.Operator.Value == "==" || binop.Operator.Value == "!=" {
			// bool expr
			val, err := evaluateBoolExpr(left.(BooleanValue), right.(BooleanValue), binop.Operator)

			if err != nil {
				parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.StartPos, binop.EndPos, err.Error()).Display()
			}

			return val
		} else {

			parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.Operator.StartPos, binop.Operator.EndPos, errMsg).Display()

			return nil
		}

	case "+=", "-=", "*=", "/=", "%=":
		if !helpers.TypesMatchT[ast.IdentifierExpr](binop.Left) || !helpers.TypesMatchT[IntegerValue](right) {

			parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.Operator.StartPos, binop.Operator.EndPos, errMsg).Display()

			return nil
		}

		if env.HasVariable((binop.Left).(ast.IdentifierExpr).Identifier) {

			exprVal, err := evaluateNumericExpr(left, right, binop.Operator)

			if err != nil {
				parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.StartPos, binop.EndPos, err.Error()).Display()
			}

			runtimeVal, err := env.AssignVariable((binop.Left).(ast.IdentifierExpr).Identifier, exprVal)

			if err != nil {
				parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.StartPos, binop.EndPos, err.Error()).Display()
			}

			if err != nil {
				parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.StartPos, binop.EndPos, err.Error()).Display()
			}

			return runtimeVal
		} else {
			val, err := evaluateNumericExpr(left, right, binop.Operator)

			if err != nil {
				parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.StartPos, binop.EndPos, err.Error()).Display()
			}

			return val
		}

	case "&&", "||":
		val, err := evaluateComparisonExpr(left, right, binop.Operator)

		if err != nil {
			parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.StartPos, binop.EndPos, err.Error()).Display()
		}

		return val

	default:

		parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.Operator.StartPos, binop.Operator.EndPos, errMsg).Display()

		return nil
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
			variableNameString = assignNode.Assigne.(ast.StructPropertyExpr).Object.(ast.IdentifierExpr).Identifier
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
	case lexer.PLUS_EQUALS_TOKEN, lexer.MINUS_EQUALS_TOKEN, lexer.TIMES_EQUALS_TOKEN, lexer.DIVIDE_EQUALS_TOKEN, lexer.MODULO_EQUALS_TOKEN:
		if (!IsINT(currentValueOfIdentifier) || !IsINT(valueToSet)) && (!IsFLOAT(currentValueOfIdentifier) || !IsFLOAT(valueToSet)){

			err = fmt.Errorf("invalid operation between %v and %v", GetRuntimeType(currentValueOfIdentifier), GetRuntimeType(valueToSet))

			parser.MakeError(env.parser, assignNode.StartPos.Line, env.parser.FilePath, assignNode.Operator.StartPos, assignNode.Operator.EndPos, err.Error()).Display()
		}

		valueToSet, err = evaluateNumericExpr(currentValueOfIdentifier, valueToSet, assignNode.Operator)

		if err != nil {

			valStart, valEnd := assignNode.Value.GetPos()

			parser.MakeError(env.parser, assignNode.StartPos.Line, env.parser.FilePath, valStart, valEnd, err.Error()).Display()
		}
	}

	runtimeVal, err := env.AssignVariable(variableToAssign.Identifier, valueToSet)

	if err != nil {
		start, end := assignNode.Value.GetPos()
		parser.MakeError(env.parser, assignNode.StartPos.Line, env.parser.FilePath, start, end, err.Error()).Display()
	}

	return runtimeVal
}

func evaluateNumericExpr(left RuntimeValue, right RuntimeValue, operator lexer.Token) (RuntimeValue, error) {

	if left != nil && right != nil {

		// evaluate both left, right as a, b where a and b can be int or float
		// if a is int return int, if a is float return float
		if IsINT(left) {
			if IsINT(right) {
				return evaluateIntInt(left.(IntegerValue), right.(IntegerValue), operator)
			} else {
				return evaluateIntFloat(left.(IntegerValue), right.(FloatValue), operator)
			}
		} else if IsFLOAT(left) {
			if IsINT(right) {
				return evaluateFloatInt(left.(FloatValue), right.(IntegerValue), operator)
			} else {
				return evaluateFloatFloat(left.(FloatValue), right.(FloatValue), operator)
			}
		}
	}
	return nil, fmt.Errorf("cannot evaluate numeric operation. unsupported operator %v", operator.Value)
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
		return MAKE_INT(left.Value+right.Value, highestBit, true), nil
	case "-", "-=":
		return MAKE_INT(left.Value-right.Value, highestBit, true), nil
	case "*", "*=":
		return MAKE_INT(left.Value*right.Value, highestBit, true), nil
	case "/", "/=":
		if right.Value == 0 {
			return nil, fmt.Errorf("division by zero is forbidden")
		}
		return MAKE_INT(left.Value/right.Value, highestBit, true), nil
	case "%", "%=":
		if right.Value == 0 {
			return nil, fmt.Errorf("division by zero is forbidden")
		}
		return MAKE_INT(left.Value%right.Value, highestBit, true), nil
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

		return MAKE_INT(result, highestBit, true), nil
	default:
		return nil, fmt.Errorf("cannot evaluate numeric operation. unsupported operator %v", operator.Value)
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
		return MAKE_INT(int64(float64(left.Value)+right.Value), highestBit, true), nil
	case "-", "-=":
		return MAKE_INT(int64(float64(left.Value)-right.Value), highestBit, true), nil
	case "*", "*=":
		return MAKE_INT(int64(float64(left.Value)*right.Value), highestBit, true), nil
	case "/", "/=":
		if right.Value == 0 {
			return nil, fmt.Errorf("division by zero is forbidden")
		}
		return MAKE_INT(int64(float64(left.Value)/right.Value), highestBit, true), nil
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

		return MAKE_INT(int64(result), highestBit, true), nil
	default:
		return nil, fmt.Errorf("cannot evaluate numeric operation. unsupported operator %v", operator.Value)
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
		return MAKE_FLOAT(left.Value+float64(right.Value), highestBit), nil
	case "-", "-=":
		return MAKE_FLOAT(left.Value-float64(right.Value), highestBit), nil
	case "*", "*=":
		return MAKE_FLOAT(left.Value*float64(right.Value), highestBit), nil
	case "/", "/=":
		if right.Value == 0 {
			return nil, fmt.Errorf("division by zero is forbidden")
		}
		return MAKE_FLOAT(left.Value/float64(right.Value), highestBit), nil
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

		return MAKE_FLOAT(result, highestBit), nil
	default:
		return nil, fmt.Errorf("cannot evaluate numeric operation. unsupported operator %v", operator.Value)
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
		return MAKE_FLOAT(left.Value+right.Value, highestBit), nil
	case "-", "-=":
		return MAKE_FLOAT(left.Value-right.Value, highestBit), nil
	case "*", "*=":
		return MAKE_FLOAT(left.Value*right.Value, highestBit), nil
	case "/", "/=":
		if right.Value == 0 {
			return nil, fmt.Errorf("division by zero is forbidden")
		}
		return MAKE_FLOAT(left.Value/right.Value, highestBit), nil
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

		return MAKE_FLOAT(result, highestBit), nil
	default:
		return nil, fmt.Errorf("cannot evaluate numeric operation. unsupported operator %v", operator.Value)
	}
}

func evaluateComparisonExpr(left RuntimeValue, right RuntimeValue, operator lexer.Token) (RuntimeValue, error) {

	result := false

	castedLeft := left.(FloatValue)
	castedRight := right.(FloatValue)

	// these can be compared to int or float
	switch operator.Value {
	case ">":
		result = castedLeft.Value > castedRight.Value
	case "<":
		result = castedLeft.Value < castedRight.Value
	case ">=":
		result = castedLeft.Value >= castedRight.Value
	case "<=":
		result = castedLeft.Value <= castedRight.Value
	default:
		return nil, fmt.Errorf("cannot evaluate comparison expression. unsupported operator %v", operator.Value)
	}

	return MAKE_BOOL(result), nil
}

func evaluateBoolExpr(left RuntimeValue, right RuntimeValue, operator lexer.Token) (RuntimeValue, error) {
	result := false

	sameType := helpers.TypesMatch(left, right)

	switch operator.Value {
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
		return nil, fmt.Errorf("cannot evaluate boolean expression. unsupported operator %v", operator)
	}

	return MAKE_BOOL(result), nil
}

func evaluateStringExpr(left StringValue, right StringValue, operator lexer.Token) (RuntimeValue, error) {
	result := false

	switch operator.Value {
	case "==":
		result = left.Value == right.Value
	case "!=":
		result = left.Value != right.Value
	default:
		return nil, fmt.Errorf("cannot evaluate string operation. unsupported operator %v", operator)
	}

	return MAKE_BOOL(result), nil
}

func evaluateStringConcat(left StringValue, right StringValue) (RuntimeValue, error) {
	return MAKE_STRING(left.Value + right.Value), nil
}
