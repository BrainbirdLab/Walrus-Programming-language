package ast

import "strconv"

type DATA_TYPE string

const (
	// Primitive Types
	VOID      DATA_TYPE = "void"
	INTEGER8  DATA_TYPE = "i8"
	INTEGER16 DATA_TYPE = "i16"
	INTEGER32 DATA_TYPE = "i32"
	INTEGER64 DATA_TYPE = "i64"

	UNSIGNED8  DATA_TYPE = "u8"
	UNSIGNED16 DATA_TYPE = "u16"
	UNSIGNED32 DATA_TYPE = "u32"
	UNSIGNED64 DATA_TYPE = "u64"

	FLOAT32    DATA_TYPE = "f32"
	FLOAT64    DATA_TYPE = "f64"
	BOOLEAN    DATA_TYPE = "boolean"
	STRING     DATA_TYPE = "str"
	CHARACTER  DATA_TYPE = "chr"
	NULL       DATA_TYPE = "null"

	// Derived Types
	ARRAY DATA_TYPE = "array"

	STRUCT   DATA_TYPE = "struct"
	TRAIT    DATA_TYPE = "trait"
	ENUM     DATA_TYPE = "enum"
	FUNCTION DATA_TYPE = "function"

	//User Defined Types
	USER_DEFINED DATA_TYPE = "user_defined"
)

type Integer struct {
	Kind     DATA_TYPE
	BitSize  uint8
	IsSigned bool
}

func (i Integer) IType() DATA_TYPE {
	if i.IsSigned {
		return DATA_TYPE(("i" + strconv.Itoa(int(i.BitSize))))
	} else {
		return DATA_TYPE(("u" + strconv.Itoa(int(i.BitSize))))
	}
}

type Float struct {
	Kind    DATA_TYPE
	BitSize uint8
}

func (f Float) IType() DATA_TYPE {
	return DATA_TYPE(("f" + strconv.Itoa(int(f.BitSize))))
}

type Boolean struct {
	Kind DATA_TYPE
}

func (b Boolean) IType() DATA_TYPE {
	return b.Kind
}

type String struct {
	Kind DATA_TYPE
}

func (s String) IType() DATA_TYPE {
	return s.Kind
}

type Char struct {
	Kind DATA_TYPE
}

func (c Char) IType() DATA_TYPE {
	return c.Kind
}

type Null struct {
	Kind DATA_TYPE
}

func (n Null) IType() DATA_TYPE {
	return n.Kind
}

type Void struct {
	Kind DATA_TYPE
}

func (v Void) IType() DATA_TYPE {
	return v.Kind
}

type Array struct {
	Kind        DATA_TYPE
	ElementType Type
	Size        int
}

func (a Array) IType() DATA_TYPE {
	return a.Kind
}

type Struct struct {
	Kind DATA_TYPE
	Name string
}

func (s Struct) IType() DATA_TYPE {
	return s.Kind
}

type Trait struct {
	Kind DATA_TYPE
	Name string
	For  string
}

func (t Trait) IType() DATA_TYPE {
	return t.Kind
}

type Enum struct {
	Kind   DATA_TYPE
	Fields map[string]Type
}

func (e Enum) IType() DATA_TYPE {
	return e.Kind
}

type UserDefined struct {
	Kind DATA_TYPE
	Name string
}

func (u UserDefined) IType() DATA_TYPE {
	return u.Kind
}

type FunctionType struct {
	Kind DATA_TYPE
}

func (f FunctionType) IType() DATA_TYPE {
	return f.Kind
}
