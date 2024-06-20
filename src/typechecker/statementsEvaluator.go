package typechecker

import (
	"fmt"
	"walrus/frontend/ast"
	"walrus/frontend/lexer"
	"walrus/frontend/parser"
	"walrus/helpers"
)

func EvaluateProgramBlock(block ast.ProgramStmt, env *Environment) RuntimeValue {

	var lastEvaluated RuntimeValue = MAKE_NULL()

	for _, stmt := range block.Contents {
		lastEvaluated = Evaluate(stmt, env)
	}

	return lastEvaluated
}

func EvaluateVariableDeclarationStmt(stmt ast.VariableDclStml, env *Environment) RuntimeValue {

	var value RuntimeValue

	if stmt.Value != nil {
		value = Evaluate(stmt.Value, env)
	} else {
		value = MAKE_NULL()
	}

	if stmt.ExplicitType != nil {

		// size checking for the integer and float types
		var explicitSize uint8 = 0

		switch t := stmt.ExplicitType.(type) {
		case ast.Integer:
			explicitSize = t.BitSize
		case ast.Float:
			explicitSize = t.BitSize
		}

		switch v := value.(type) {
		case IntegerValue:
			//modify the Size field of the value. Update the original value
			v.Size = explicitSize // is it reference or copy? It is a copy. So, the original value is not updated
			t := v.Type.(ast.Integer)
			t.BitSize = explicitSize
			v.Type = t
			//update the Original value
			value = v
		case FloatValue:
			v.Size = explicitSize
			t := v.Type.(ast.Float)
			t.BitSize = explicitSize
			v.Type = t
			value = v
		}

		//check user defined types with the value type
		start, end := stmt.Value.GetPos()
		checkTypes(env.parser, stmt.ExplicitType, value, start, end)
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
	if udt, ok := expected.(ast.UserDefined); ok {
		name = udt.Name
	} else {
		name = string(expected.IType())
	}
	return fmt.Sprintf("cannot assign value of type '%s' to '%s'", GetRuntimeType(got), name)
}

func checkTypes(p *parser.Parser, explicitType ast.Type, value RuntimeValue, startPos lexer.Position, endPos lexer.Position) {

	var msg string

	switch t := explicitType.(type) {
	case ast.Integer:
		if IsINT(value) {
			if t.BitSize != value.(IntegerValue).Size {
				msg = strFormatter(explicitType, value)
				msg += fmt.Sprintf(" of size %d to integer of size %d", value.(IntegerValue).Size, t.BitSize)
			}
		} else {
			msg = strFormatter(explicitType, value)
		}
	case ast.Float:
		if IsFLOAT(value) {
			if t.BitSize != value.(FloatValue).Size {
				msg = strFormatter(explicitType, value)
				msg += fmt.Sprintf(" of size %d to float of size %d", value.(FloatValue).Size, t.BitSize)
			}
		} else {
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

	var lastEvaluated RuntimeValue = MAKE_NULL()

	for _, stmt := range block.Body {
		lastEvaluated = Evaluate(stmt, env)
	}

	return lastEvaluated

}

func EvaluateControlFlowStmt(astNode ast.IfStmt, env *Environment) RuntimeValue {

	condition := Evaluate(astNode.Condition, env)

	if IsTruthy(condition) {
		return Evaluate(astNode.Block, env)
	} else {
		for astNode.Alternate != nil && helpers.TypesMatchT[ast.IfStmt](astNode.Alternate) {
			alt := astNode.Alternate.(ast.IfStmt)
			condition = Evaluate(alt.Condition, env)
			if IsTruthy(condition) {
				return Evaluate(alt.Block, env)
			}
		}

		if astNode.Alternate != nil && helpers.TypesMatchT[ast.BlockStmt](astNode.Alternate) {
			return Evaluate(astNode.Alternate.(ast.BlockStmt), env)
		}
	}

	return MAKE_NULL()
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

	for _, body := range stmt.Block.Body {
		switch t := body.(type) {
		case ast.VariableDclStml:
			val := Evaluate(t.Value, funcEnv)
			funcEnv.DeclareVariable(t.Identifier.Identifier, val, t.IsConstant)
		case ast.ReturnStmt:
			returnStmt = &t
		}
	}

	// a void function should not have a return statement with a value. It should be empty like return;
	if stmt.ReturnType.IType() == ast.VOID {
		if returnStmt != nil {
			if returnStmt.Kind != ast.NODE_TYPE(ast.VOID) {
				//return nil, fmt.Errorf("void function %s must not have a return statement with a value", name)
				parser.MakeError(funcEnv.parser, returnStmt.StartPos.Line, funcEnv.parser.FilePath, returnStmt.StartPos, returnStmt.EndPos, "void function must not have a return statement with a value").Display()
			}
		}
	} else {
		if returnStmt == nil {
			//return nil, fmt.Errorf("function %s must have a return statement", name)
			parser.MakeError(funcEnv.parser, stmt.StartPos.Line, funcEnv.parser.FilePath, stmt.Name.StartPos, stmt.Name.EndPos, "function must have a return statement").Display()
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
				parser.MakeError(funcEnv.parser, start.Line, funcEnv.parser.FilePath, start, end, fmt.Sprintf("cannot return value of type '%s' from function %s with return type '%s'", returnType, funcName, stmt.ReturnType.IType())).Display()
			}
		}
	}

	return MAKE_VOID()
}

func EvaluateFunctionCallExpr(expr ast.FunctionCallExpr, env *Environment) RuntimeValue {

	// check if the function is defined
	funcName := expr.Function.Identifier

	if env.variables[funcName] == nil {
		parser.MakeError(env.parser, expr.StartPos.Line, env.parser.FilePath, expr.Function.StartPos, expr.Function.EndPos, fmt.Sprintf("function '%s' is not defined", funcName)).Display()
	}

	function := env.variables[funcName].(FunctionValue)

	args := expr.Args

	params := function.Parameters
	body := function.Body

	// check if the number of arguments match the number of parameters
	if len(args) != len(params) {
		parser.MakeError(env.parser, expr.StartPos.Line, env.parser.FilePath, expr.StartPos, expr.EndPos, fmt.Sprintf("function '%s' expects %d arguments but %d were provided", funcName, len(params), len(args))).Display()
	}

	// create a new environment for the function
	newEnv := NewEnvironment(env, env.parser)

	// check and set the arguments to the function parameters
	for i := 0; i < len(params); i++ {
		param := params[i]
		arg := args[i]

		if param.Type != nil {
			start, end := arg.GetPos()
			checkTypes(env.parser, param.Type, Evaluate(arg, env), start, end)
		}

		newEnv.DeclareVariable(param.Identifier.Identifier, Evaluate(arg, env), false)
	}

	lastVal := Evaluate(body, newEnv)

	return lastVal
}

func EvaluateReturnStmt(stmt ast.ReturnStmt, env *Environment) RuntimeValue {
	return Evaluate(stmt.Expression, env)
}
