package ast

import "strconv"

type DATA_TYPE string

const (
	// Primitive Types
	T_VOID      DATA_TYPE = "void"
	T_INTEGER8  DATA_TYPE = "i8"
	T_INTEGER16 DATA_TYPE = "i16"
	T_INTEGER32 DATA_TYPE = "i32"
	T_INTEGER64 DATA_TYPE = "i64"

	T_UNSIGNED8  DATA_TYPE = "u8"
	T_UNSIGNED16 DATA_TYPE = "u16"
	T_UNSIGNED32 DATA_TYPE = "u32"
	T_UNSIGNED64 DATA_TYPE = "u64"

	T_FLOAT32   DATA_TYPE = "f32"
	T_FLOAT64   DATA_TYPE = "f64"
	T_BOOLEAN   DATA_TYPE = "boolean"
	T_STRING    DATA_TYPE = "str"
	T_CHARACTER DATA_TYPE = "chr"
	T_NULL      DATA_TYPE = "null"

	// Derived Types
	T_ARRAY DATA_TYPE = "array"

	T_STRUCT   DATA_TYPE = "struct"
	T_TRAIT    DATA_TYPE = "trait"
	T_ENUM     DATA_TYPE = "enum"
	T_FUNCTION DATA_TYPE = "function"

	//User Defined Types
	T_USER_DEFINED DATA_TYPE = "user_defined"
)

type IntegerType struct {
	Kind     DATA_TYPE
	BitSize  uint8
	IsSigned bool
}

func (i IntegerType) IType() DATA_TYPE {
	if i.IsSigned {
		return DATA_TYPE(("i" + strconv.Itoa(int(i.BitSize))))
	} else {
		return DATA_TYPE(("u" + strconv.Itoa(int(i.BitSize))))
	}
}

type FloatType struct {
	Kind    DATA_TYPE
	BitSize uint8
}

func (f FloatType) IType() DATA_TYPE {
	return DATA_TYPE(("f" + strconv.Itoa(int(f.BitSize))))
}

type BoolType struct {
	Kind DATA_TYPE
}

func (b BoolType) IType() DATA_TYPE {
	return b.Kind
}

type StringType struct {
	Kind DATA_TYPE
}

func (s StringType) IType() DATA_TYPE {
	return s.Kind
}

type CharType struct {
	Kind DATA_TYPE
}

func (c CharType) IType() DATA_TYPE {
	return c.Kind
}

type NullType struct {
	Kind DATA_TYPE
}

func (n NullType) IType() DATA_TYPE {
	return n.Kind
}

type VoidType struct {
	Kind DATA_TYPE
}

func (v VoidType) IType() DATA_TYPE {
	return v.Kind
}

type ArrayType struct {
	Kind        DATA_TYPE
	ElementType Type
	Size        int
}

func (a ArrayType) IType() DATA_TYPE {
	return a.Kind
}

type StructType struct {
	Kind DATA_TYPE
	Name string
}

func (s StructType) IType() DATA_TYPE {
	return s.Kind
}

type TraitType struct {
	Kind DATA_TYPE
	Name string
	For  string
}

func (t TraitType) IType() DATA_TYPE {
	return t.Kind
}

type EnumType struct {
	Kind   DATA_TYPE
	Fields []Type
}

func (e EnumType) IType() DATA_TYPE {
	return e.Kind
}

type FunctionType struct {
	Kind       DATA_TYPE
	Name       string
	ReturnType Type
	Parameters []FunctionParameter
}

func (f FunctionType) IType() DATA_TYPE {
	return f.Kind
}
