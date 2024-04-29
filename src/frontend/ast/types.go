package ast

type IntegerType struct{
	BitSize     uint8
	IsSigned    bool
}
func (i IntegerType) _type() {}

type FloatingType struct {
	BitSize uint8
}
func (f FloatingType) _type() {}

type BooleanType struct{}
func (b BooleanType) _type() {}

type StringType struct{}
func (s StringType) _type() {}

type CharecterType struct{}
func (c CharecterType) _type() {}

type NullType struct{}
func (n NullType) _type() {}

type ArrayType struct {
	ElementType Type
	Size        int
}
func (a ArrayType) _type() {}
