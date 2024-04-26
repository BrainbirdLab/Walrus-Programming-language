package ast

type SymbolType struct {
	Name string
}
func (s SymbolType) _type() {}

type ArrayType struct {
	UnderlayingType Type
}
func (a ArrayType) _type() {}