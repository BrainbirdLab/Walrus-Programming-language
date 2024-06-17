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
	parser    *parser.Parser
}

func NewEnvironment(parent *Environment, p *parser.Parser) *Environment {
	return &Environment{
		parent:    parent,
		variables: make(map[string]RuntimeValue),
		constants: make(map[string]bool),
		parser:    p,
	}
}

func (e *Environment) DeclareVariable(name string, value RuntimeValue, isConstant bool) RuntimeValue {
	if e.variables[name] != nil {
		panic(fmt.Sprintf("Variable %s already declared in this scope\n", name))
	}

	e.variables[name] = value

	if isConstant {
		e.constants[name] = true
	}

	return value
}

func (e *Environment) AssignVariable(name string, value RuntimeValue) RuntimeValue {
	env := e.ResolveVariable(name)

	if env.constants[name] {
		panic(fmt.Sprintf("Cannot assign value to constant %s", name))
	}

	env.variables[name] = value

	return value
}

func (e *Environment) DeclareFunction(name string, parameters map[string]ast.Type, body ast.BlockStmt) RuntimeValue {
	if e.variables[name] != nil {
		panic(fmt.Sprintf("Identifier (Function) %s already declared in this scope\n", name))
	}

	e.variables[name] = FunctionValue{
		Name:       name,
		Parameters: parameters,
		Body:       body,
	}

	return e.variables[name]
}

func (e *Environment) ResolveVariable(name string) *Environment {
	if _, ok := e.variables[name]; ok {
		return e
	}

	if e.parent != nil {
		panic(fmt.Sprintf("Variable %s was not declared in this scope\n", name))
	}

	return e.parent.ResolveVariable(name)
}

func (e *Environment) GetRuntimeValue(name string) RuntimeValue {
	env := e.ResolveVariable(name)
	return env.variables[name]
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
