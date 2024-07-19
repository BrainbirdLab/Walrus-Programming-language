package typechecker

import (
	"fmt"
	"walrus/frontend/ast"
	"walrus/frontend/parser"
)

type Environment struct {
	parent    *Environment
	variables map[string]RuntimeValue
	constants map[string]bool
	//user defined types declared with struct keyword
	structs map[string]RuntimeValue
	parser    *parser.Parser
}

func NewEnvironment(parent *Environment, p *parser.Parser) *Environment {
	return &Environment{
		parent:    parent,
		variables: make(map[string]RuntimeValue),
		constants: make(map[string]bool),
		structs:   make(map[string]RuntimeValue),
		parser:    p,
	}
}

func (e *Environment) DeclareVariable(name string, value RuntimeValue, isConstant bool) (RuntimeValue, error) {

	if e.variables[name] != nil {
		return nil, fmt.Errorf("variable %s already declared in this scope", name)
	}

	switch v := value.(type){
	case StructInstance:
		// check all fields are initialized
		structDeclaration := e.structs[v.StructName]

		for _, field := range structDeclaration.(StructValue).Fields {
			if v.Fields[field.Name] == nil {
				return nil, fmt.Errorf("field '%s' of struct '%s' is not initialized", field.Name, v.StructName)
			} else {
				// check type compatibility
				if GetRuntimeType(v.Fields[field.Name]) != field.Type.IType() {
					return nil, fmt.Errorf("field '%s' of struct '%s' is of type %s, but got %s", field.Name, v.StructName, field.Type.IType(), GetRuntimeType(v.Fields[field.Name]))
				}
			}
		}
	}

	e.variables[name] = value

	if isConstant {
		e.constants[name] = true
	}

	return value, nil
}

func (e *Environment) AssignVariable(name string, value RuntimeValue) (RuntimeValue, error) {

	env, err := e.ResolveVariable(name)

	if err != nil {
		return nil, err
	}

	if env.constants[name] {
		return nil, fmt.Errorf("cannot assign value to constant %s", name)
	}

	// check type compatibility
	variable := env.variables[name]

	if !((IsBothINT(variable, value) || IsBothFLOAT(variable, value)) || (GetRuntimeType(variable) == GetRuntimeType(value))) {
		return nil, fmt.Errorf("cannot assign value %v of type %s to %s", value, GetRuntimeType(value), GetRuntimeType(variable))
	} else {
		switch t := variable.(type) {
		case IntegerValue:
			if t.Size < value.(IntegerValue).Size {
				return nil, fmt.Errorf("potential data loss. %d bit value cannot be assigned to %s of size %d. You can try type casting", value.(IntegerValue).Size, t.Type, t.Size)
			}
		case FloatValue:
			if t.Size < value.(FloatValue).Size {
				return nil, fmt.Errorf("potential data loss. %d bit value cannot be assigned to %s of size %d. You can try type casting", value.(FloatValue).Size, t.Type, t.Size)
			}
		}
	}

	env.variables[name] = value

	return value, nil
}

func (e *Environment) DeclareFunction(name string, returnType ast.Type, parameters []ast.FunctionParameter, body ast.BlockStmt) error {

	if e.variables[name] != nil {
		return fmt.Errorf("identifier (function) %s already declared in this scope", name)
	}

	e.variables[name] = FunctionValue{
		Name:       name,
		Parameters: parameters,
		Body:       body,
		Type: 		ast.T_FN,
		ReturnType: returnType.IType(),
		DeclarationEnv: e,
	}

	e.constants[name] = true

	return nil
}

func (e *Environment) DeclareNativeFn(name string, fn RuntimeValue) error {

	if e.variables[name] != nil {
		return fmt.Errorf("identifier (function) %s already declared in this scope", name)
	}

	e.variables[name] = fn

	e.constants[name] = true

	return nil
}

func (e *Environment) ResolveVariable(name string) (*Environment, error) {

	if _, ok := e.variables[name]; ok {
		return e, nil
	}

	if e.parent == nil {
		return nil, fmt.Errorf("variable %s was not declared in this scope", name)
	}

	return e.parent.ResolveVariable(name)
}

func (e *Environment) GetStructType(name string) (RuntimeValue, error) {
	
	if _, ok := e.structs[name]; ok {
		return e.structs[name], nil
	}

	if e.parent == nil {
		return nil, fmt.Errorf("struct %s was not declared in this scope", name)
	}

	return e.parent.GetStructType(name)
}

func (e *Environment) GetRuntimeValue(name string) (RuntimeValue, error) {
	env, err := e.ResolveVariable(name)

	if err != nil {
		return nil, err
	}

	return env.variables[name], nil
}

func (e *Environment) HasVariable(name string) bool {

	if _, ok := e.variables[name]; ok {
		return true
	}
	if e.parent == nil {
		return false
	}
	return e.parent.HasVariable(name)
}
