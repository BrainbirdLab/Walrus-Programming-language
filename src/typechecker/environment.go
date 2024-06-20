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

	if GetRuntimeType(variable) != GetRuntimeType(value) {
		return nil, fmt.Errorf("cannot assign value of type %s to %s", GetRuntimeType(value), GetRuntimeType(variable))
	} else {
		switch t := variable.(type) {
		case IntegerValue:
			if t.Size < value.(IntegerValue).Size {
				return nil, fmt.Errorf("potential data loss. %d bit value cannot be assigned to %s of size %d. You can try type casting", value.(IntegerValue).Size, t.Type.IType(), t.Size)
			}
		case FloatValue:
			if t.Size < value.(FloatValue).Size {
				return nil, fmt.Errorf("potential data loss. %d bit value cannot be assigned to %s of size %d. You can try type casting", value.(FloatValue).Size, t.Type.IType(), t.Size)
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
		Type: ast.FunctionType{
			Kind: ast.T_FUNCTION,
		},
		ReturnType: returnType,
	}

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
