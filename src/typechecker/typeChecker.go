package typechecker

import (
	"fmt"
	"strconv"
	"walrus/frontend/ast"
	"walrus/frontend/parser"
)

func GetRuntimeType(runtimeValue RuntimeValue) ast.DATA_TYPE {
	switch t := runtimeValue.(type) {
	case IntegerValue:
		return t.Type.IType()
	case FloatValue:
		return t.Type.IType()
	case BooleanValue:
		return t.Type.IType()
	case StringValue:
		return t.Type.IType()
	case CharacterValue:
		return t.Type.IType()
	case NullValue:
		return t.Type.IType()
	case VoidValue:
		return t.Type.IType()
	case FunctionValue:
		return t.Type.IType()
	case NativeFunctionValue:
		return t.Type.IType()
	case StructValue:
		return t.Type.IType()
	case StructInstance:
		return ast.DATA_TYPE(t.StructName)
	default:
		panic(fmt.Sprintf("This runtime value is not implemented yet: %T", runtimeValue))
	}
}

func GetNumericValue(runtimeValue RuntimeValue) (float64, error) {
	//cast to float64
	switch t := runtimeValue.(type) {
	case IntegerValue:
		return float64(t.Value), nil
	case FloatValue:
		return t.Value, nil
	case BooleanValue:
		if t.Value {
			return 1, nil
		} else {
			return 0, nil
		}
	case CharacterValue:
		return float64(t.Value), nil
	default:
		return 0, fmt.Errorf("cannot convert %T to a numeric value", runtimeValue)
	}
}

func IsINT(runtimeValue RuntimeValue) bool {
	switch GetRuntimeType(runtimeValue) {
	case ast.T_INTEGER8, ast.T_INTEGER16, ast.T_INTEGER32, ast.T_INTEGER64:
		return true
	default:
		return false
	}
}

func IsBothINT(runtimeValue1 RuntimeValue, runtimeValue2 RuntimeValue) bool {
	return IsINT(runtimeValue1) && IsINT(runtimeValue2)
}

func IsFLOAT(runtimeValue RuntimeValue) bool {
	switch GetRuntimeType(runtimeValue) {
	case ast.T_FLOAT32, ast.T_FLOAT64:
		return true
	default:
		return false
	}
}

func IsBothFLOAT(runtimeValue1 RuntimeValue, runtimeValue2 RuntimeValue) bool {
	return IsFLOAT(runtimeValue1) && IsFLOAT(runtimeValue2)
}

func IsNumber(runtimeValue RuntimeValue) bool {
	return IsINT(runtimeValue) || IsFLOAT(runtimeValue)
}

func IsString(runtimeValue RuntimeValue) bool {
	return GetRuntimeType(runtimeValue) == ast.T_STRING
}

func IsArithmetic(value RuntimeValue) bool {
	switch value.(type) {
	case IntegerValue, FloatValue, CharacterValue, BooleanValue:
		return true
	default:
		return false
	}
}

func IsFunction(value RuntimeValue) bool {
	return GetRuntimeType(value) == ast.T_FUNCTION || GetRuntimeType(value) == ast.T_NATIVE_FN
}

func CastToStringValue(value RuntimeValue) (StringValue, error) {

	switch t := value.(type) {
	case StringValue:
		return t, nil
	case IntegerValue:
		return MAKE_STRING(strconv.FormatInt(t.Value, 10)), nil
	case FloatValue:
		return MAKE_STRING(strconv.FormatFloat(t.Value, 'f', -1, 64)), nil
	case BooleanValue:
		return MAKE_STRING(strconv.FormatBool(t.Value)), nil
	case CharacterValue:
		return MAKE_STRING(string(t.Value)), nil
	default:
		return StringValue{}, fmt.Errorf("cannot cast %T to string", value)
	}
}

func Evaluate(astNode ast.Node, env *Environment) RuntimeValue {
	switch node := astNode.(type) {
	case ast.NumericLiteral:
		// Check if the number is an integer or a float
		if node.BaseStmt.Kind == ast.INTEGER_LITERAL {
			val, _ := strconv.ParseInt(node.Value, 10, int(node.BitSize))
			return MAKE_INT(val, node.BitSize, true)
		} else if node.BaseStmt.Kind == ast.FLOAT_LITERAL {
			val, _ := strconv.ParseFloat(node.Value, 64)
			return MAKE_FLOAT(val, node.BitSize)
		} else {
			parser.MakeError(env.parser, node.StartPos.Line, env.parser.FilePath, node.StartPos, node.EndPos, "invalid numeric literal").Display()
			return nil
		}
	case ast.StringLiteral:
		return MAKE_STRING(node.Value)
	case ast.CharacterLiteral:
		if len(node.Value) > 1 {
			parser.MakeError(env.parser, node.StartPos.Line, env.parser.FilePath, node.StartPos, node.EndPos, "character literals can only have one character").Display()
		}
		return MAKE_CHAR(node.Value[0])
	case ast.BooleanLiteral:
		return MAKE_BOOL(node.Value)
	case ast.NullLiteral:
		return MAKE_NULL()
	case ast.ProgramStmt:
		return EvaluateProgramBlock(node, env)
	case ast.VariableDclStml:
		return EvaluateVariableDeclarationStmt(node, env)
	case ast.AssignmentExpr:
		return EvaluateAssignmentExpr(node, env)
	case ast.UnaryExpr:
		return EvaluateUnaryExpression(node, env)
	case ast.BinaryExpr:
		return EvaluateBinaryExpr(node, env)
	case ast.IdentifierExpr:
		return EvaluateIdenitifierExpr(node, env)
	case ast.BlockStmt:
		return EvaluateBlockStmt(node, env)
	case ast.IfStmt:
		return EvaluateControlFlowStmt(node, env)
	case ast.FunctionDeclStmt:
		return EvaluateFunctionDeclarationStmt(node, env)
	case ast.FunctionCallExpr:
		return EvaluateFunctionCallExpr(node, env)
	case ast.ReturnStmt:
		return EvaluateReturnStmt(node, env)
	case ast.StructDeclStatement:
		return EvaluateStructDeclarationStmt(node, env)
	case ast.StructLiteral:
		return EvaluateStructLiteral(node, env)
	case ast.StructPropertyExpr:
		return EvaluateStructPropertyExpr(node, env)
	default:
		panic(fmt.Sprintf("This ast node is not implemented yet: %v", node))
	}
}

func HasStruct(name string, env *Environment) bool {
	// if not found in the current scope, check the parent scope
	if _, ok := env.structs[name]; ok {
		return true
	}

	if env.parent != nil {
		return HasStruct(name, env.parent)
	}

	return false
}
