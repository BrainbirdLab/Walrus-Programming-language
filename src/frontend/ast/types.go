package ast

type DATA_TYPE string

const (
	// Primitive Types
	VOID DATA_TYPE = "VOID"

	INTEGER   = "INTEGER"
	FLOATING  = "FLOATING"
	BOOLEAN   = "BOOLEAN"
	STRING    = "STRING"
	CHARACTER = "CHARECTER"
	NULL      = "NULL"

	// Derived Types
	ARRAY = "ARRAY"

	STRUCT = "STRUCT"
	TRAIT  = "TRAIT"
	ENUM   = "ENUM"

	//User Defined Types
	USER_DEFINED = "USER_DEFINED"
)

type IntegerType struct {
	Kind     DATA_TYPE
	BitSize  uint8
	IsSigned bool
}

func (i IntegerType) iType() {
	// empty method implements the Type interface
}

type FloatingType struct {
	Kind    DATA_TYPE
	BitSize uint8
}

func (f FloatingType) iType() {
	// empty method implements the Type interface
}

type BooleanType struct {
	Kind DATA_TYPE
}

func (b BooleanType) iType() {
	// empty method implements the Type interface
}

type StringType struct {
	Kind DATA_TYPE
}

func (s StringType) iType() {
	// empty method implements the Type interface
}

type CharecterType struct {
	Kind DATA_TYPE
}

func (c CharecterType) iType() {
	// empty method implements the Type interface
}

type NullType struct {
	Kind DATA_TYPE
}

func (n NullType) iType() {
	// empty method implements the Type interface
}

type VoidType struct {
	Kind DATA_TYPE
}

func (v VoidType) iType() {
	// empty method implements the Type interface
}

type ArrayType struct {
	Kind        DATA_TYPE
	ElementType Type
	Size        int
}

func (a ArrayType) iType() {
	// empty method implements the Type interface
}

type StructType struct {
	Kind DATA_TYPE
	Name string
}

func (s StructType) iType() {
	// empty method implements the Type interface
}

type TraitType struct {
	Kind DATA_TYPE
	Name string
	For  string
}

func (t TraitType) iType() {
	// empty method implements the Type interface
}

type EnumType struct {
	Kind   DATA_TYPE
	Fields map[string]Type
}

func (e EnumType) iType() {
	// empty method implements the Type interface
}

type UserDefinedType struct {
	Kind DATA_TYPE
	Name string
}

func (u UserDefinedType) iType() {
	// empty method implements the Type interface
}
