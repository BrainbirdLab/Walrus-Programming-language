package tc

import (
	"fmt"
	"walrus/frontend/ast"
)

/*
{{program {3 1 4} {3 12 15}} ./../code/test/tc/tc.wal  [] [{{variable declaration statement {3 1 4} {3 12 15}} false {{identifier {3 5 8} {3 6 9}} a} {{integer literal {3 10 13} {3 11 14}} 5 32} <nil>}]}
*/

type TypeEnv struct {
	parent 		*TypeEnv
	variables 	map[string]ast.Type
	constants 	map[string]bool
	structs 	map[string]ast.Type
}

var typeEnv = new(TypeEnv)

func (t *TypeEnv) ResolveVar(name string) (*TypeEnv, error) {
	if _, ok := t.variables[name]; ok {
		return t, nil
	}

	//search for parent scope
	if t.parent == nil {
		return nil, fmt.Errorf("tc: %s is not declared in this scope", name)
	}

	return t.parent.ResolveVar(name)
}

func (t *TypeEnv) DeclareVar(name string, valueType ast.Type, isConst bool) error {
	if _, ok := t.variables[name]; ok {
		return fmt.Errorf("%s already declared in this scope", name)
	}

	switch vType := valueType.(type) {
	case ast.StructType:
		t.DeclareStruct(name, vType)
	}

	t.variables[name] = valueType

	if isConst {
		t.constants[name] = true
	}

	return nil
}

func (t *TypeEnv) DeclareStruct(name string, valueType ast.StructType) error {
	declaredStruct := t.structs[name]
	fmt.Printf("declared: %v\n", declaredStruct.(ast.StructType))
	return nil
}

func CheckType(astNode ast.Node, env *TypeEnv) (ast.Type, error) {
	fmt.Println("Checking types")
	fmt.Printf("%v\n", astNode)


	switch node := (astNode).(type) {
	case ast.ProgramStmt:
		return checkProgram(&node, typeEnv)
	case ast.VariableDclStml:
		return checkVarDecl(&node, typeEnv)
	default:
		return nil, fmt.Errorf("Error")
	}
}

func checkProgram(program *ast.ProgramStmt, env *TypeEnv) (ast.Type, error) {
	for _, item := range (*program).Contents {
		_, err := CheckType(item, env)
		if err != nil {
			return nil, err
		}
	}
	return ast.VoidType{
		Kind: ast.T_VOID,
	}, nil
}

func checkVarDecl(varDecl *ast.VariableDclStml, env *TypeEnv) (ast.Type, error) {

	iden := varDecl.Identifier
	valueToSet := varDecl.Value

	//env.DeclareVar(iden.Identifier, valueToSet)

	fmt.Printf("var: %v, value: %v\n", iden, valueToSet)
	return ast.VoidType{}, nil
}