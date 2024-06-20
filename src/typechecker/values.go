package typechecker

import (
	"fmt"
	"walrus/frontend/ast"
)

type RuntimeValue interface {
	rVal()
}

type IntegerValue struct {
	Value int64
	Size  uint8
	Type  ast.Type
}

func (i IntegerValue) rVal() {
	// empty function implements RuntimeValue interface
}

type FloatValue struct {
	Value float64
	Size  uint8
	Type  ast.Type
}

func (f FloatValue) rVal() {
	// empty function implements RuntimeValue interface
}

type BooleanValue struct {
	Value bool
	Type  ast.Type
}

func (b BooleanValue) rVal() {
	// empty function implements RuntimeValue interface
}

type StringValue struct {
	Value string
	Type  ast.Type
}

func (s StringValue) rVal() {
	// empty function implements RuntimeValue interface
}

type CharacterValue struct {
	Value byte
	Type  ast.Type
}

func (c CharacterValue) rVal() {
	// empty function implements RuntimeValue interface
}

type NullValue struct {
	Type ast.Type
}

func (n NullValue) rVal() {
	// empty function implements RuntimeValue interface
}

type VoidValue struct {
	Type ast.Type
}

func (v VoidValue) rVal() {
	// empty function implements RuntimeValue interface
}

type FunctionValue struct {
	Name       string
	Parameters []ast.FunctionParameter
	Body       ast.BlockStmt
	Type       ast.Type
	ReturnType ast.Type
}

func (f FunctionValue) rVal() {
	// empty function implements RuntimeValue interface
}

type FunctionCall struct {
	FunctionName string
	Arguments    []RuntimeValue
	Returns      RuntimeValue
	Type         ast.Type
}

func (f FunctionCall) rVal() {
	// empty function implements RuntimeValue interface
}

func MAKE_INT(value int64, size uint8, signed bool) IntegerValue {

	initial := "i"

	if !signed {
		initial = "u"
	}

	return IntegerValue{Value: value, Size: size, Type: ast.Integer{
			Kind:     ast.DATA_TYPE((initial + fmt.Sprintf("%d", size))),
			BitSize:  size,
			IsSigned: signed,
		},
	}
}

func MAKE_FLOAT(value float64, size uint8) FloatValue {
	
	return FloatValue{Value: value, Size: size, Type: ast.Float{
			Kind:    ast.DATA_TYPE(("f" + fmt.Sprintf("%d", size))),
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

func MakeDefaultRuntimeValue(t ast.Type) RuntimeValue {

	switch t := t.(type) {
	case ast.Integer:
		return MAKE_INT(0, t.BitSize, t.IsSigned)
	case ast.Float:
		return MAKE_FLOAT(0, t.BitSize)
	case ast.Boolean:
		return MAKE_BOOL(false)
	case ast.String:
		return MAKE_STRING("")
	case ast.Char:
		return MAKE_CHAR(0)
	case ast.Null:
		return MAKE_NULL()
	case ast.Void:
		return MAKE_VOID()
	default:
		panic(fmt.Sprintf("unsupported type %T", t))
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
