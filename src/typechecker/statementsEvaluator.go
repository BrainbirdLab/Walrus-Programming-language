package typechecker

import (
	"fmt"
	"walrus/frontend/ast"
	"walrus/frontend/lexer"
	"walrus/frontend/parser"
)

func EvaluateProgramBlock(block ast.ProgramStmt, env *Environment) RuntimeValue {
	for _, stmt := range block.Contents {
		rVal := Evaluate(stmt, env)
		if _, ok := rVal.(ReturnValue); ok {
			return rVal
		}
	}
	return MakeVOID()
}

func EvaluateVariableDeclarationStmt(stmt ast.VariableDclStml, env *Environment) RuntimeValue {

	var value RuntimeValue

	if stmt.Value != nil {
		value = Evaluate(stmt.Value, env)
	}

	if stmt.ExplicitType != nil {

		// size checking for the integer and float types
		var explicitSize uint8 = 0

		switch t := stmt.ExplicitType.(type) {
		case ast.IntegerType:
			explicitSize = t.BitSize

		case ast.FloatType:
			explicitSize = t.BitSize
		}

		switch v := value.(type) {
		case IntegerValue:
			//modify the Size field of the value. Update the original value
			v.Size = explicitSize // is it reference or copy? It is a copy. So, the original value is not updated
			t := v.Type.(ast.IntegerType)
			t.BitSize = explicitSize
			v.Type = t
			//update the Original value
			value = v
		case FloatValue:
			v.Size = explicitSize
			t := v.Type.(ast.FloatType)
			t.BitSize = explicitSize
			v.Type = t
			value = v
		}

		if value == nil {
			value = MakeDefaultRuntimeValue(stmt.ExplicitType)
		}

		//check user defined types with the value type
		start, end := stmt.Identifier.GetPos()
		checkTypes(env, stmt.ExplicitType, value, start, end)
	}

	val, err := env.DeclareVariable(stmt.Identifier.Identifier, value, stmt.IsConstant)

	if err != nil {
		parser.MakeError(env.parser, stmt.StartPos.Line, env.parser.FilePath, stmt.Identifier.StartPos, stmt.Identifier.EndPos, err.Error()).Display()
	}

	return val
}

func strFormatter(expected ast.Type, got RuntimeValue) string {
	var name string
	//if expected is userdefined type
	if userDefinedType, ok := expected.(ast.StructType); ok {
		name = userDefinedType.Name
	} else {
		name = string(expected.IType())
	}
	return fmt.Sprintf("cannot assign value of type '%s' to '%s'", GetRuntimeType(got), name)
}

func checkTypes(env *Environment, explicitType ast.Type, value RuntimeValue, startPos lexer.Position, endPos lexer.Position) {

	p := env.parser

	var msg string

	switch t := explicitType.(type) {
	case ast.IntegerType:
		if IsINT(value) {
			if t.BitSize != value.(IntegerValue).Size {
				msg = strFormatter(explicitType, value)
				msg += fmt.Sprintf(" of size %d to integer of size %d", value.(IntegerValue).Size, t.BitSize)
			}
		} else {
			msg = strFormatter(explicitType, value)
		}
	case ast.FloatType:
		if IsFLOAT(value) {
			if t.BitSize != value.(FloatValue).Size {
				msg = strFormatter(explicitType, value)
				msg += fmt.Sprintf(" of size %d to float of size %d", value.(FloatValue).Size, t.BitSize)
			}
		} else {
			msg = strFormatter(explicitType, value)
		}
	case ast.StructType:
		expected := t.Name
		got := string(GetRuntimeType(value))

		if !HasStruct(expected, env) {
			parser.MakeError(p, startPos.Line, p.FilePath, startPos, endPos, fmt.Sprintf("failed to validate types. struct '%s' is not defined", expected)).Display()
		} else if !HasStruct(string(got), env) {
			parser.MakeError(p, startPos.Line, p.FilePath, startPos, endPos, fmt.Sprintf("failed to validate types. struct '%s' is not defined", got)).Display()
		} else if expected != got {
			msg = strFormatter(explicitType, value)
		}
	default:
		if GetRuntimeType(value) != t.IType() {
			msg = strFormatter(explicitType, value)
		}
	}

	if msg != "" {
		parser.MakeError(p, startPos.Line, p.FilePath, startPos, endPos, msg).Display()
	}
}

func EvaluateBlockStmt(block ast.BlockStmt, env *Environment) RuntimeValue {
	for _, stmt := range block.Items {
		switch stmt := stmt.(type) {
		case ast.ReturnStmt:
			// Evaluate the return expression and return its value immediately
			return Evaluate(stmt, env)
		default:
			rVal := Evaluate(stmt, env)
			//check if the runtime value is a return value
			if _, ok := rVal.(ReturnValue); ok {
				return rVal
			}
		}
	}

	return ReturnValue{
		Value: MakeVOID(),
	}
}

func EvaluateControlFlowStmt(astNode ast.IfStmt, env *Environment) RuntimeValue {

	condition := Evaluate(astNode.Condition, env)

	if IsTruthy(condition) {
		return Evaluate(astNode.Block, env)
	}

	if astNode.Alternate != nil {
		switch t := astNode.Alternate.(type) {
		case ast.IfStmt:
			return EvaluateControlFlowStmt(t, env)
		case ast.BlockStmt:
			return Evaluate(t, env)
		}
	}

	return MakeNULL()
}

func EvaluateFunctionDeclarationStmt(stmt ast.FunctionDeclStmt, env *Environment) RuntimeValue {

	err := env.DeclareFunction(stmt.Name.Identifier, stmt.ReturnType, stmt.Parameters, stmt.Block)

	if err != nil {
		parser.MakeError(env.parser, stmt.StartPos.Line, env.parser.FilePath, stmt.Name.StartPos, stmt.Name.EndPos, err.Error()).Display()
	}

	// eliminate the return statement from the body

	var returnStmt *ast.ReturnStmt

	funcEnv := NewEnvironment(env, env.parser)

	//for each parameter, declare a variable in the function environment
	for _, param := range stmt.Parameters {
		funcEnv.DeclareVariable(param.Identifier.Identifier, MakeDefaultRuntimeValue(param.Type), false)
	}

	for _, body := range stmt.Block.Items {
		switch t := body.(type) {
		case ast.VariableDclStml:
			val := Evaluate(t.Value, funcEnv)
			funcEnv.DeclareVariable(t.Identifier.Identifier, val, t.IsConstant)
		case ast.ReturnStmt:
			returnStmt = &t
		}
	}

	// a void function should not have a return statement with a value. It should be empty like return;
	if stmt.ReturnType.IType() == ast.T_VOID {
		if returnStmt != nil {
			if returnStmt.Kind != ast.NODE_TYPE(ast.T_VOID) {
				//return nil, fmt.Errorf("void function %s must not have a return statement with a value", name)
				parser.MakeError(funcEnv.parser, returnStmt.StartPos.Line, funcEnv.parser.FilePath, returnStmt.StartPos, returnStmt.EndPos, "void function must not have a return statement with a value").Display()
			}
		}
	} else {
		if returnStmt == nil {
			//return nil, fmt.Errorf("function %s must have a return statement", name)
			parser.MakeError(funcEnv.parser, stmt.StartPos.Line, funcEnv.parser.FilePath, stmt.Name.StartPos, stmt.Name.EndPos, "function must have a return value at the end").Display()
		} else {
			// check the return type of the function
			start, end := returnStmt.GetPos()

			// evaluate the return statement
			returnVal := Evaluate(returnStmt.Expression, funcEnv)

			returnType := GetRuntimeType(returnVal)

			if returnType != stmt.ReturnType.IType() {

				funcName := fmt.Sprintf("%s(", stmt.Name.Identifier)
				params := stmt.Parameters
				formattedParams := ""
				for i, param := range params {
					if i == 0 {
						formattedParams += fmt.Sprintf("%s: %s", param.Identifier.Identifier, param.Type.IType())
					} else {
						formattedParams += fmt.Sprintf(", %s: %s", param.Identifier.Identifier, param.Type.IType())
					}
				}
				funcName += formattedParams + ")"

				// if struct, check if the struct is defined
				if _, ok := stmt.ReturnType.(ast.StructType); ok {
					if HasStruct(stmt.ReturnType.(ast.StructType).Name, funcEnv) {
						return MakeVOID()
					} else {
						parser.MakeError(funcEnv.parser, start.Line, funcEnv.parser.FilePath, start, end, fmt.Sprintf("cannot return value of type '%s' from function %s with return type '%s'", returnType, stmt.Name.Identifier, stmt.ReturnType.(ast.StructType).Name)).Display()
					}
				} else {
					parser.MakeError(funcEnv.parser, start.Line, funcEnv.parser.FilePath, start, end, fmt.Sprintf("cannot return value of type '%s' from function %s with return type '%s'", returnType, funcName, stmt.ReturnType.IType())).Display()
				}
			}
		}
	}

	return MakeVOID()
}

func EvaluateFunctionCallExpr(expr ast.FunctionCallExpr, env *Environment) RuntimeValue {

	var args []RuntimeValue

	for _, arg := range expr.Args {
		args = append(args, Evaluate(arg, env))
	}

	fn := Evaluate(expr.Caller, env)

	if !IsFunction(fn) {
		parser.MakeError(env.parser, expr.StartPos.Line, env.parser.FilePath, expr.StartPos, expr.EndPos, fmt.Sprintf("could not call. %s not a function", expr.Caller.Identifier)).Display()
	}

	if GetRuntimeType(fn) == ast.T_NATIVE_FN {
		return fn.(NativeFunctionValue).Caller(args...)
	}

	function := fn.(FunctionValue)
	scope := NewEnvironment(function.DeclarationEnv, env.parser)

	params := function.Parameters

	// check if the number of arguments match the number of parameters
	if len(args) != len(params) {
		parser.MakeError(env.parser, expr.StartPos.Line, env.parser.FilePath, expr.StartPos, expr.EndPos, fmt.Sprintf("function '%s' expects %d arguments but %d were provided", function.Name, len(params), len(args))).Display()
	}

	// check and set the arguments to the function parameters
	for i := 0; i < len(params); i++ {
		scope.DeclareVariable(params[i].Identifier.Identifier, args[i], false)
	}

	for _, stmt := range function.Body.Items {
		rVal := Evaluate(stmt, scope)
		if _, ok := rVal.(ReturnValue); ok {
			return rVal.(ReturnValue).Value
		}
	}

	return MakeVOID()
}

func EvaluateReturnStmt(stmt ast.ReturnStmt, env *Environment) RuntimeValue {
	expr := stmt.Expression
	val := Evaluate(expr, env)

	return ReturnValue{
		Value: val,
	}
}

func EvaluateStructDeclarationStmt(stmt ast.StructDeclStatement, env *Environment) RuntimeValue {

	env.structs[stmt.StructName] = StructValue{
		Fields:  stmt.Properties,
		Methods: stmt.Methods,
		Type: ast.StructType{
			Kind: ast.T_STRUCT,
			Name: stmt.StructName,
		},
	}

	return MakeVOID()
}

func EvaluateStructLiteral(stmt ast.StructLiteral, env *Environment) RuntimeValue {

	//check if the struct is defined
	if !HasStruct(stmt.StructName, env) {
		parser.MakeError(env.parser, stmt.StartPos.Line, env.parser.FilePath, stmt.StartPos, stmt.EndPos, fmt.Sprintf("cannot evaluate struct literal. struct '%s' is not defined", stmt.StructName)).Display()
	}

	properties := make(map[string]RuntimeValue)

	for name, value := range stmt.Properties {
		properties[name] = Evaluate(value, env)
	}

	return StructInstance{
		StructName: stmt.StructName,
		Fields:     properties,
	}
}

func EvaluateStructPropertyExpr(expr ast.StructPropertyExpr, env *Environment) RuntimeValue {

	obj := Evaluate(expr.Object, env).(StructInstance)

	propname := expr.Property.Identifier

	if obj.Fields[propname] == nil {
		parser.MakeError(env.parser, expr.StartPos.Line, env.parser.FilePath, expr.Property.StartPos, expr.Property.EndPos, fmt.Sprintf("property '%s' is not defined in struct '%s'", propname, obj.StructName)).Display()
	}

	structValue, err := env.GetStructType(obj.StructName)

	if err != nil {
		parser.MakeError(env.parser, expr.StartPos.Line, env.parser.FilePath, expr.StartPos, expr.EndPos, err.Error()).Display()
	}

	// check if the property is public
	if structValue.(StructValue).Fields[propname].IsPublic {
		return obj.Fields[propname]
	} else {
		parser.MakeError(env.parser, expr.StartPos.Line, env.parser.FilePath, expr.Property.StartPos, expr.Property.EndPos, fmt.Sprintf("property '%s' is private in struct '%s'", propname, obj.StructName)).Display()
		return nil
	}
}
