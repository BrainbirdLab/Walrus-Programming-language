package ast

type DATA_TYPE string

const (
	// Primitive Types
	VOID DATA_TYPE 	= "VOID"

	INTEGER   		= "INTEGER"
	FLOATING  		= "FLOATING"
	BOOLEAN   		= "BOOLEAN"
	STRING    		= "STRING"
	CHARACTER 		= "CHARECTER"
	NULL      		= "NULL"

	// Derived Types
	ARRAY 			= "ARRAY"

	STRUCT 			= "STRUCT"
	TRAIT			= "TRAIT"
	ENUM			= "ENUM"

	//User Defined Types
	USER_DEFINED	= "USER_DEFINED"
)

type IntegerType struct {
	Kind     DATA_TYPE
	BitSize  uint8
	IsSigned bool
}

func (i IntegerType) _type() {}

type FloatingType struct {
	Kind     DATA_TYPE
	BitSize  uint8
}

func (f FloatingType) _type() {}

type BooleanType struct {
	Kind     DATA_TYPE
}

func (b BooleanType) _type() {}

type StringType struct {
	Kind     DATA_TYPE
}

func (s StringType) _type() {}

type CharecterType struct {
	Kind     DATA_TYPE
}

func (c CharecterType) _type() {}

type NullType struct {
	Kind     DATA_TYPE
}

func (n NullType) _type() {}

type VoidType struct {
	Kind     DATA_TYPE
}

func (v VoidType) _type() {}

type ArrayType struct {
	Kind        DATA_TYPE
	ElementType Type
	Size        int
}

func (a ArrayType) _type() {}

type StructType struct {
	Kind     DATA_TYPE
	Name	 string
}
func (s StructType) _type() {}

type TraitType struct {
	Kind     	DATA_TYPE
	Name		string
	For 		string
}
func (t TraitType) _type() {}

type EnumType struct {
	Kind     DATA_TYPE
	Fields   map[string]Type
}
func (e EnumType) _type() {}

type UserDefinedType struct {
	Kind     DATA_TYPE
	Name     string
}
func (u UserDefinedType) _type() {}