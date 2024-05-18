package typechecker

import "fmt"

type Environment struct {
	parent    *Environment
	variables map[string]RuntimeValue
	constants map[string]bool
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

	if env.constants[name] == true {
		panic(fmt.Sprintf("Cannot assign value to constant %s", name))
	}

	env.variables[name] = value

	return value
}

func (e *Environment) ResolveVariable(name string) *Environment {
	if _, ok := e.variables[name]; ok {
		return e
	}

	if e.parent != nil {
		panic(fmt.Sprintf("Variable %s not declared\n", name))
	}

	return e.parent.ResolveVariable(name)
}

func (e *Environment) GetRuntimeValue(name string) RuntimeValue {
	env := e.ResolveVariable(name)
	return env.variables[name]
}

func (e *Environment) HasVariable(name string) bool {
	env := e.ResolveVariable(name)
	_, ok := env.variables[name]
	return ok
}
