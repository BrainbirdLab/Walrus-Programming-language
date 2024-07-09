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

type ReturnValue struct {
	Value RuntimeValue
}

func (r ReturnValue) rVal() {
	// empty function implements RuntimeValue interface
}

type FunctionValue struct {
	Name           string
	Parameters     []ast.FunctionParameter
	Body           ast.BlockStmt
	Type           ast.Type
	ReturnType     ast.Type
	DeclarationEnv *Environment
}

func (f FunctionValue) rVal() {
	// empty function implements RuntimeValue interface
}

type StructValue struct {
	Fields  map[string]ast.Property
	Methods map[string]ast.FunctionType
	Type    ast.Type
}

func (s StructValue) rVal() {
	// empty function implements RuntimeValue interface
}

type StructInstance struct {
	StructName string
	Fields     map[string]RuntimeValue
}

func (s StructInstance) rVal() {
	// empty function implements RuntimeValue interface
}

type FunctionCall = func(...RuntimeValue) RuntimeValue

type NativeFunctionValue struct {
	Caller FunctionCall
	Type   ast.Type
}

func (n NativeFunctionValue) rVal() {
	// empty function implements RuntimeValue interface
}

func MakeINT(value int64, size uint8, signed bool) IntegerValue {

	initial := "i"

	if !signed {
		initial = "u"
	}

	return IntegerValue{Value: value, Size: size, Type: ast.IntegerType{
		Kind:     ast.DATA_TYPE((initial + fmt.Sprintf("%d", size))),
		BitSize:  size,
		IsSigned: signed,
	},
	}
}

func MakeFLOAT(value float64, size uint8) FloatValue {

	return FloatValue{Value: value, Size: size, Type: ast.FloatType{
		Kind:    ast.DATA_TYPE(("f" + fmt.Sprintf("%d", size))),
		BitSize: size,
	},
	}
}

func MakeBOOL(value bool) BooleanValue {
	return BooleanValue{Value: value, Type: ast.BoolType{
		Kind: ast.T_BOOLEAN,
	},
	}
}

func MakeSTRING(value string) StringValue {
	return StringValue{Value: value, Type: ast.StringType{
		Kind: ast.T_STRING,
	},
	}
}

func MakeCHAR(value byte) CharacterValue {
	return CharacterValue{Value: value, Type: ast.CharType{
		Kind: ast.T_CHARACTER,
	},
	}
}

func MakeNULL() NullValue {
	return NullValue{Type: ast.NullType{
		Kind: ast.T_NULL,
	},
	}
}

func MakeVOID() VoidValue {
	return VoidValue{Type: ast.VoidType{
		Kind: ast.T_VOID,
	},
	}
}

func MakeNativeFUNCTION(call FunctionCall) NativeFunctionValue {
	return NativeFunctionValue{
		Caller: call,
		Type: ast.NativeFnType{
			Kind: ast.T_NATIVE_FN,
		},
	}
}

func MakeDefaultRuntimeValue(t ast.Type) RuntimeValue {

	switch t := t.(type) {
	case ast.IntegerType:
		return MakeINT(0, t.BitSize, t.IsSigned)
	case ast.FloatType:
		return MakeFLOAT(0, t.BitSize)
	case ast.BoolType:
		return MakeBOOL(false)
	case ast.StringType:
		return MakeSTRING("")
	case ast.CharType:
		return MakeCHAR(0)
	case ast.NullType:
		return MakeNULL()
	case ast.VoidType:
		return MakeVOID()
	case ast.StructType:
		return StructValue{
			Fields:  make(map[string]ast.Property),
			Methods: make(map[string]ast.FunctionType),
			Type:    t,
		}
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
