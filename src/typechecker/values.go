package typechecker

import (
	"fmt"
	"walrus/frontend/ast"
)

type RuntimeValue interface{
	_val()
}

type IntegerValue struct {
	Value 	int64
	Size 	uint8
	Type 	ast.Type
}
func (i IntegerValue) _val() {}

type FloatValue struct {
	Value 	float64
	Size 	uint8
	Type 	ast.Type
}
func (f FloatValue) _val() {}

type BooleanValue struct {
	Value 	bool
	Type 	ast.Type
}
func (b BooleanValue) _val() {}

type StringValue struct {
	Value 	string
	Type 	ast.Type
}
func (s StringValue) _val() {}

type CharacterValue struct {
	Value 	byte
	Type 	ast.Type
}
func (c CharacterValue) _val() {}

type NullValue struct {
	Type 	ast.Type
}
func (n NullValue) _val() {}

type VoidValue struct {
	Type 	ast.Type
}
func (v VoidValue) _val() {}


type FunctionValue struct {
	Name 			string
	Parameters 		[]ast.FunctionParameter
	Body 			ast.BlockStmt
	Type			ast.Type
	ReturnType 		ast.Type
}
func (f FunctionValue) _val() {}

type FunctionCall struct {
	Value 	ast.Type
	Type 	ast.Type
}
func (f FunctionCall) _val() {}

func MAKE_INT(value int64, size uint8, signed bool) IntegerValue {
	return IntegerValue{Value: value, Size: size, Type: ast.Integer{
			Kind: ast.INTEGER,
			BitSize: size,
			IsSigned: signed,
		},
	}
}

func MAKE_FLOAT(value float64, size uint8) FloatValue {
	return FloatValue{Value: value, Size: size, Type: ast.Float{
			Kind: ast.FLOATING,
			BitSize: size,
		},
	}
}

func MAKE_BOOL(value bool) BooleanValue {
	return BooleanValue{Value: value, Type: ast.Boolean{
			Kind: ast.BOOLEAN,
		},
	}
}

func MAKE_STRING(value string) StringValue {
	return StringValue{Value: value, Type: ast.String{
			Kind: ast.STRING,
		},
	}
}

func MAKE_CHAR(value byte) CharacterValue {
	return CharacterValue{Value: value, Type: ast.Char{
			Kind: ast.CHARACTER,
		},
	}
}

func MAKE_NULL() NullValue {
	return NullValue{Type: ast.Null{
			Kind: ast.NULL,
		},
	}
}

func MAKE_VOID() VoidValue {
	return VoidValue{Type: ast.Void{
			Kind: ast.VOID,
		},
	}
}

func IsTruthy(value RuntimeValue) bool {
	if value == nil {
		return false
	}

	switch value := value.(type) {
	case IntegerValue:
		return value.Value != 0
	case FloatValue:
		return value.Value != 0
	case BooleanValue:
		return value.Value
	case StringValue:
		return value.Value != ""
	case CharacterValue:
		return value.Value != 0
	case NullValue, VoidValue:
		return false
	default:
		panic(fmt.Sprintf("unsupported type %T", value))
	}
}