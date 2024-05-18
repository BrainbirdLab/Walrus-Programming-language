package typechecker

type RuntimeValue interface{
	_val()
}

type IntegerValue struct {
	Value 	int
	Size 	uint8
	Type 	string
}
func (i IntegerValue) _val() {}

type FloatValue struct {
	Value 	float64
	Size 	uint8
	Type 	string
}
func (f FloatValue) _val() {}

type BooleanValue struct {
	Value 	bool
	Type 	string
}
func (b BooleanValue) _val() {}

type StringValue struct {
	Value 	string
	Type 	string
}
func (s StringValue) _val() {}

type CharacterValue struct {
	Value 	byte
	Type 	string
}
func (c CharacterValue) _val() {}

type NullValue struct {
	Type 	string
}
func (n NullValue) _val() {}

type VoidValue struct {
	Type 	string
}
func (v VoidValue) _val() {}

func MAKE_INT(value int, size uint8) IntegerValue {
	return IntegerValue{Value: value, Size: size, Type: "INTEGER"}
}

func MAKE_FLOAT(value float64, size uint8) FloatValue {
	return FloatValue{Value: value, Size: size, Type: "FLOAT"}
}

func MAKE_BOOL(value bool) BooleanValue {
	return BooleanValue{Value: value, Type: "BOOL"}
}

func MAKE_STRING(value string) StringValue {
	return StringValue{Value: value, Type: "STRING"}
}

func MAKE_CHAR(value byte) CharacterValue {
	return CharacterValue{Value: value, Type: "CHARACTER"}
}

func MAKE_NULL() NullValue {
	return NullValue{Type: "NULL"}
}

func MAKE_VOID() VoidValue {
	return VoidValue{Type: "VOID"}
}