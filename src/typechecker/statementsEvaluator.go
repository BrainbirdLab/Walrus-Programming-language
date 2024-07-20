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
			t := v.Type
			v.Type = t
			//update the Original value
			value = v
		case FloatValue:
			v.Size = explicitSize
			t := v.Type
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

func checkTypes(env *Environment, explicitType ast.Type, value RuntimeValue, startPos lexer.Position, endPos lexer.Position) {

	switch t := explicitType.(type) {
	case ast.IntegerType:
		checkIntegerType(env, t, value, startPos, endPos)
	case ast.FloatType:
		checkFloatType(env, t, value, startPos, endPos)
	case ast.StructType:
		checkStructType(env, t, value, startPos, endPos)
	default:
		checkGeneralType(env, t, value, startPos, endPos)
	}
}

func checkIntegerType(env *Environment, explicitType ast.IntegerType, value RuntimeValue, startPos lexer.Position, endPos lexer.Position) {
	if IsINT(value) {
		if explicitType.BitSize != value.(IntegerValue).Size {
			displayTypeMismatchError(env.parser, explicitType, value, startPos, endPos, fmt.Sprintf("integer of size %d", explicitType.BitSize))
		}
	} else {
		displayTypeMismatchError(env.parser, explicitType, value, startPos, endPos, "")
	}
}

func checkFloatType(env *Environment, explicitType ast.FloatType, value RuntimeValue, startPos lexer.Position, endPos lexer.Position) {
	if IsFLOAT(value) {
		if explicitType.BitSize != value.(FloatValue).Size {
			displayTypeMismatchError(env.parser, explicitType, value, startPos, endPos, fmt.Sprintf("float of size %d", explicitType.BitSize))
		}
	} else {
		displayTypeMismatchError(env.parser, explicitType, value, startPos, endPos, "")
	}
}

func checkStructType(env *Environment, explicitType ast.StructType, value RuntimeValue, startPos lexer.Position, endPos lexer.Position) {
	expected := string(explicitType.Kind)
	got := string(GetRuntimeType(value))

	if !HasStruct(expected, env) {
		parser.MakeError(env.parser, startPos.Line, env.parser.FilePath, startPos, endPos, fmt.Sprintf("failed to validate types. struct '%s' is not defined", expected)).Display()
	} else if !HasStruct(got, env) {
		parser.MakeError(env.parser, startPos.Line, env.parser.FilePath, startPos, endPos, fmt.Sprintf("failed to validate types. struct '%s' is not defined", got)).Display()
	} else if expected != got {
		displayTypeMismatchError(env.parser, explicitType, value, startPos, endPos, "")
	}
}

func checkGeneralType(env *Environment, explicitType ast.Type, value RuntimeValue, startPos lexer.Position, endPos lexer.Position) {
	if GetRuntimeType(value) != explicitType.IType() {
		displayTypeMismatchError(env.parser, explicitType, value, startPos, endPos, "")
	}
}

func displayTypeMismatchError(p *parser.Parser, explicitType ast.Type, value RuntimeValue, startPos lexer.Position, endPos lexer.Position, additionalInfo string) {
	msg := strFormatter(explicitType, value)
	if additionalInfo != "" {
		msg += fmt.Sprintf(" to %s", additionalInfo)
	}
	parser.MakeError(p, startPos.Line, p.FilePath, startPos, endPos, msg).Display()
}

func strFormatter(expected ast.Type, got RuntimeValue) string {
	var name string
	switch t := expected.(type) {
	case ast.IntegerType:
		name = fmt.Sprintf("integer of size %d", t.BitSize)
	case ast.FloatType:
		name = fmt.Sprintf("float of size %d", t.BitSize)
	case ast.StructType:
		if userDefinedType, ok := expected.(ast.StructType); ok {
			name = string(userDefinedType.Kind)
		}
	default:
		name = string(expected.IType())
	}
	return fmt.Sprintf("cannot assign value of type '%s' to '%s'", GetRuntimeType(got), name)
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
	err := declareFunction(stmt, env)
	if err != nil {
		handleFunctionDeclarationError(stmt, env, err)
		return MakeVOID()
	}

	funcEnv := createFunctionEnvironment(stmt, env)
	processFunctionBody(stmt.Block, funcEnv)
	checkFunctionReturnType(stmt, funcEnv)

	return MakeVOID()
}

func declareFunction(stmt ast.FunctionDeclStmt, env *Environment) error {
	return env.DeclareFunction(stmt.Name.Identifier, stmt.ReturnType, stmt.Parameters, stmt.Block)
}

func handleFunctionDeclarationError(stmt ast.FunctionDeclStmt, env *Environment, err error) {
	parser.MakeError(env.parser, stmt.StartPos.Line, env.parser.FilePath, stmt.Name.StartPos, stmt.Name.EndPos, err.Error()).Display()
}

func createFunctionEnvironment(stmt ast.FunctionDeclStmt, env *Environment) *Environment {
	funcEnv := NewEnvironment(env, env.parser)

	for _, param := range stmt.Parameters {
		funcEnv.DeclareVariable(param.Identifier.Identifier, MakeDefaultRuntimeValue(param.Type), false)
	}

	return funcEnv
}

func processFunctionBody(body ast.BlockStmt, funcEnv *Environment) {
	var returnStmt *ast.ReturnStmt

	for _, stmt := range body.Items {
		switch stmt := stmt.(type) {
		case ast.VariableDclStml:
			val := Evaluate(stmt.Value, funcEnv)
			funcEnv.DeclareVariable(stmt.Identifier.Identifier, val, stmt.IsConstant)
		case ast.ReturnStmt:
			returnStmt = &stmt
		}
	}

	if returnStmt != nil && returnStmt.Kind == ast.NODE_TYPE(ast.T_VOID) {
		parser.MakeError(funcEnv.parser, returnStmt.StartPos.Line, funcEnv.parser.FilePath, returnStmt.StartPos, returnStmt.EndPos, "void function must not have a return statement with a value").Display()
	}
}

func checkFunctionReturnType(stmt ast.FunctionDeclStmt, funcEnv *Environment) {
	if stmt.ReturnType.IType() != ast.T_VOID {
		if len(stmt.Block.Items) == 0 {
			fmt.Println(stmt.StartPos, stmt.EndPos)
			parser.MakeError(funcEnv.parser, stmt.StartPos.Line, funcEnv.parser.FilePath, stmt.StartPos, stmt.EndPos, "no return statement found").AddHint("function is empty", parser.TEXT_HINT).Display()
		}
		lastStmt := stmt.Block.Items[len(stmt.Block.Items)-1]
		if returnStmt, ok := lastStmt.(ast.ReturnStmt); ok {
			returnVal := Evaluate(returnStmt.Expression, funcEnv)
			expectedType := fmt.Sprintf("%s", stmt.ReturnType)
			returnType := GetRuntimeType(returnVal)
			if GetRuntimeType(returnVal) != stmt.ReturnType.IType() {
				parser.MakeError(funcEnv.parser, returnStmt.StartPos.Line, funcEnv.parser.FilePath, returnStmt.StartPos, returnStmt.EndPos, fmt.Sprintf("cannot return value of type '%s' from function with return type '%s'", returnType, expectedType)).Display()
			}
		} else {
			parser.MakeError(funcEnv.parser, stmt.StartPos.Line, funcEnv.parser.FilePath, stmt.Name.StartPos, stmt.Name.EndPos, "function must have a return value at the end").Display()
		}
	}
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
		param := params[i].Identifier
		arg := args[i]

		expected := params[i].Type.IType()
		got := GetRuntimeType(arg)

		if expected != got {
			fmt.Printf("param: %v, arg: %v", params[i].Type.IType(), GetRuntimeType(arg))
			parser.MakeError(env.parser, expr.StartPos.Line, env.parser.FilePath, expr.StartPos, expr.EndPos, "function parameter and arguments type mismatched").AddHint(fmt.Sprintf("expected type '%s' but got '%s'", expected, got), parser.TEXT_HINT).Display()
		}

		scope.DeclareVariable(param.Identifier, arg, false)
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
		Type:    ast.DATA_TYPE(stmt.StructName),
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

func EvaluatePropertyExpr(expr ast.PropertyExpr, env *Environment) RuntimeValue {

	propname := expr.Property.Identifier

	switch obj := Evaluate(expr.Object, env).(type) {
	case StructInstance:
	
		if obj.Fields[propname] == nil {
			parser.MakeError(env.parser, expr.StartPos.Line, env.parser.FilePath, expr.Property.StartPos, expr.Property.EndPos, fmt.Sprintf("property '%s' does not exist in type '%s'", propname, obj.StructName)).Display()
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
	case ArrayValue:
		switch propname {
		case "length":
			//return the length of the array
			size := len(obj.Values)
			return MakeINT(int64(size), 32, true)
		default:
			parser.MakeError(env.parser, expr.StartPos.Line, env.parser.FilePath, expr.Property.StartPos, expr.Property.EndPos, fmt.Sprintf("property '%s' does not exist in type array", propname)).Display()
		}
	}
	return nil
}

func EvaluateArrayLiterals(node ast.Node, env *Environment) RuntimeValue {
	var values []RuntimeValue

	for _, value := range node.(ast.ArrayLiterals).Elements {
		values = append(values, Evaluate(value, env))
	}

	return ArrayValue{
		Values: values,
		Type:   ast.T_ARRAY,
	}
}

func EvaluateArrayAccess(node ast.Node, env *Environment) RuntimeValue {
	arr := node.(ast.ArrayIndexAccess)
	name := arr.ArrayName
	scope, err := env.ResolveVariable(name)
	errorPrinter := parser.MakeError(env.parser, arr.StartPos.Line, env.parser.FilePath, arr.StartPos, arr.EndPos, "array was not declared\n")
	if err != nil {
		errorPrinter.Display()
	}

	indexNumber := Evaluate(arr.Index, env)

	if _, ok := indexNumber.(IntegerValue); !ok {
		errorPrinter.Message = "invalid index value\n"
		errorPrinter.AddHint("index must be a valid integer\n", parser.TEXT_HINT).Display()
	}

	index := indexNumber.(IntegerValue).Value
	values := scope.variables[name].(ArrayValue).Values

	if index > int64(len(values)-1) {
		errorPrinter.Message = fmt.Sprintf("invalid index range %d\n", index)
		errorPrinter.AddHint(fmt.Sprintf("index must be within the range of 0 to %d\n", len(values)-1), parser.TEXT_HINT).Display()
	}

	return values[index]
}
