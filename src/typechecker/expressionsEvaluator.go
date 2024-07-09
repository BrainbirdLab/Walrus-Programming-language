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

	errMsg := fmt.Sprintf("unsupported unary operation for type %v", expr)

	switch unary.Operator.Value {
	case "-", "+":
		if !helpers.TypesMatchT[IntegerValue](expr) {
			panic(errMsg)
		}

		var value int64

		if unary.Operator.Value == "-" {
			value = -expr.(IntegerValue).Value
		} else {
			value = expr.(IntegerValue).Value
		}

		return MakeINT(value, 32, true)

	case "!":
		if !helpers.TypesMatchT[BooleanValue](expr) {
			panic(errMsg)
		}

		return BooleanValue{
			Type: ast.BoolType{
				Kind: ast.T_BOOLEAN,
			},
			Value: !expr.(BooleanValue).Value,
		}

	case "++", "--":
		if !helpers.TypesMatchT[IntegerValue](expr) {
			panic(errMsg)
		}

		var value int64

		if unary.Operator.Value == "++" {
			value = expr.(IntegerValue).Value + 1
		} else {
			value = expr.(IntegerValue).Value - 1
		}

		runtimeVal := MakeINT(value, 32, true)

		if helpers.TypesMatchT[ast.IdentifierExpr](unary.Argument) {
			if env.HasVariable(unary.Argument.(ast.IdentifierExpr).Identifier) {
				env.AssignVariable(unary.Argument.(ast.IdentifierExpr).Identifier, runtimeVal)
			}
		}

		return runtimeVal

	default:
		return MakeNULL()
	}
}

func EvaluateBinaryExpr(binop ast.BinaryExpr, env *Environment) RuntimeValue {

	left := Evaluate(binop.Left, env)
	right := Evaluate(binop.Right, env)

	leftType := GetRuntimeType(left)
	rightType := GetRuntimeType(right)

	errMsg := fmt.Sprintf("unsupported binary operation between %v and %v", leftType, rightType)

	switch binop.Operator.Value {
	//Arithmetic operators
	case "+", "-", "*", "/", "%", "^":
		if IsNumber(left) {
			if !IsNumber(right) {
				parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.Operator.StartPos, binop.Operator.EndPos, errMsg).Display()
			}

			result, err := evaluateNumericArithmeticExpr(left, right, binop.Operator)

			if err != nil {
				parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.Operator.StartPos, binop.Operator.EndPos, err.Error()).Display()
			}

			return result
		} else if IsString(left) {
			//convert right to string
			strVal, err := CastToStringValue(right)

			if err != nil {
				parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.Operator.StartPos, binop.Operator.EndPos, err.Error()).Display()
			}

			result, err := evaluateStringExpr(left.(StringValue), strVal, binop.Operator)

			if err != nil {
				parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.Operator.StartPos, binop.Operator.EndPos, err.Error()).Display()
			}

			return result
		} else {
			parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.Operator.StartPos, binop.Operator.EndPos, errMsg).Display()
		}

		//Relationa operators
	case "==", "!=", ">", "<", ">=", "<=":
		result, err := evaluateComparisonExpr(left, right, binop.Operator)

		if err != nil {
			parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.Operator.StartPos, binop.Operator.EndPos, err.Error()).Display()
		}

		return result

		//Logical operators
	case "&&":
		if IsTruthy(left) {
			if IsTruthy(right) {
				return right
			} else {
				return left
			}
		} else {
			return MakeBOOL(false)
		}

	case "||":
		if IsTruthy(left) {
			return left
		} else {
			if IsTruthy(right) {
				return right
			} else {
				return MakeBOOL(false)
			}
		}
	default:
		parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.Operator.StartPos, binop.Operator.EndPos, errMsg).Display()
	}
	return MakeNULL()
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
			return nil, divisionByZero
		}
		return MakeINT(left.Value/right.Value, highestBit, true), nil
	case "%", "%=":
		if right.Value == 0 {
			return nil, divisionByZero
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
			return nil, fmt.Errorf("division by zero is forbidden")
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

	result := false

	if GetRuntimeType(left) == GetRuntimeType(right) && GetRuntimeType(left) == ast.T_STRING {
		// string comparison
		switch operator.Value {
		case "==":
			result = left.(StringValue).Value == right.(StringValue).Value
		case "!=":
			result = left.(StringValue).Value != right.(StringValue).Value
		default:
			return nil, fmt.Errorf("operator %v is not supported in strings", operator.Value)
		}

		return MakeBOOL(result), nil
	}

	leftValue, err := GetNumericValue(left)

	if err != nil {
		return nil, err
	}

	rightValue, err := GetNumericValue(right)

	if err != nil {
		return nil, err
	}

	// these can be compared to int or float
	switch operator.Value {
	case ">":
		result = leftValue > rightValue
	case "<":
		result = leftValue < rightValue
	case ">=":
		result = leftValue >= rightValue
	case "<=":
		result = leftValue <= rightValue
	case "==":
		result = leftValue == rightValue
	case "!=":
		result = leftValue != rightValue
	// logical operators
	case "&&":
		// return the last truthy value if all are true. else return false
		if GetRuntimeType(left) == ast.T_BOOLEAN {
			if GetRuntimeType(right) == ast.T_BOOLEAN {
				if leftValue == 1 {
					if rightValue == 1 {
						return MakeBOOL(true), nil
					} else {
						return MakeBOOL(false), nil
					}
				} else {
					return MakeBOOL(false), nil
				}
			} else {
				// right is not a boolean
				if IsTruthy(left) {
					if IsTruthy(right) {
						return right, nil
					} else {
						return left, nil
					}
				} else {
					return MakeBOOL(false), nil
				}
			}
		} else {
			if IsTruthy(left) {
				if IsTruthy(right) {
					return right, nil
				} else {
					return left, nil
				}
			} else {
				return MakeBOOL(false), nil
			}
		}

	case "||":
		// return the first truthy value if any is true. else return false
		if GetRuntimeType(left) == ast.T_BOOLEAN {
			if GetRuntimeType(right) == ast.T_BOOLEAN {
				if leftValue == 1 {
					return MakeBOOL(true), nil
				} else {
					if rightValue == 1 {
						return MakeBOOL(true), nil
					} else {
						return MakeBOOL(false), nil
					}
				}
			} else {
				// right is not a boolean
				if IsTruthy(left) {
					return left, nil
				} else {
					if IsTruthy(right) {
						return right, nil
					} else {
						return MakeBOOL(false), nil
					}
				}
			}
		} else {
			if IsTruthy(left) {
				return left, nil
			} else {
				if IsTruthy(right) {
					return right, nil
				} else {
					return MakeBOOL(false), nil
				}
			}
		}

	default:
		return nil, fmt.Errorf("cannot evaluate comparison expression. unsupported operator %v", operator.Value)
	}

	return MakeBOOL(result), nil
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
