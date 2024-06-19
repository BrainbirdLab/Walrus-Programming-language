package ast

type DATA_TYPE string

const (
	// Primitive Types
	VOID      DATA_TYPE = "void"
	INTEGER   DATA_TYPE = "integer"
	FLOATING  DATA_TYPE = "float"
	BOOLEAN   DATA_TYPE = "boolean"
	STRING    DATA_TYPE = "string"
	CHARACTER DATA_TYPE = "character"
	NULL      DATA_TYPE = "null"

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
	return i.Kind
}

type Float struct {
	Kind    DATA_TYPE
	BitSize uint8
}

func (f Float) IType() DATA_TYPE {
	return f.Kind
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
	Kind         DATA_TYPE
}

func (f FunctionType) IType() DATA_TYPE {
	return f.Kind
}

func GetType(t Type) DATA_TYPE {
	switch t.(type) {
	case Integer:
		return INTEGER
	case Float:
		return FLOATING
	case Boolean:
		return BOOLEAN
	case String:
		return STRING
	case Char:
		return CHARACTER
	case Null:
		return NULL
	case Void:
		return VOID
	case Array:
		return ARRAY
	case Struct:
		return STRUCT
	case Trait:
		return TRAIT
	case Enum:
		return ENUM
	case UserDefined:
		return USER_DEFINED
	}
	return "UNKNOWN"
}
