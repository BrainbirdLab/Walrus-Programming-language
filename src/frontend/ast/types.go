package ast

type Datatype int

const (
	// Primitive Types
	INTEGER Datatype = iota
	FLOATING
	BOOLEAN
	STRING
	CHARACTER
	NULL

	// Derived Types
	ARRAY
)

type IntegerType struct{
	Kind 	  	Datatype
	BitSize     uint8
	IsSigned    bool
}
func (i IntegerType) _type() {}

type FloatingType struct {
	Kind 	  	Datatype
	BitSize 	uint8
}
func (f FloatingType) _type() {}

type BooleanType struct{
	Kind 	  	Datatype
}
func (b BooleanType) _type() {}

type StringType struct{
	Kind 	  	Datatype
}
func (s StringType) _type() {}

type CharecterType struct{
	Kind 	  	Datatype
}
func (c CharecterType) _type() {}

type NullType struct{
	Kind 	  	Datatype
}
func (n NullType) _type() {}

type ArrayType struct {
	Kind        Datatype
	ElementType Type
	Size        int
}
func (a ArrayType) _type() {}
