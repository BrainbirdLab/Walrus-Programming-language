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
		if GetRuntimeType(value) == ast.INTEGER {
			if t.BitSize != value.(IntegerValue).Size {
				msg = strFormatter(explicitType, value)
				msg += fmt.Sprintf(" of size %d to integer of size %d", value.(IntegerValue).Size, t.BitSize)
			}
		} else {
			msg = strFormatter(explicitType, value)
		}
	case ast.Float:
		if GetRuntimeType(value) == ast.FLOATING {
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

	runtimeVal, err := env.DeclareFunction(stmt.FunctionName.Identifier, stmt.ReturnType, stmt.Parameters, stmt.Block)

	if err != nil {
		parser.MakeError(env.parser, stmt.StartPos.Line, env.parser.FilePath, stmt.FunctionName.StartPos, stmt.FunctionName.EndPos, err.Error()).Display()
	}

	return runtimeVal
}

func EvaluateFunctionCallExpr(expr ast.FunctionCallExpr, env *Environment) RuntimeValue {

	// check if the function is defined
	funcName := expr.Function.Identifier

	
	if  env.variables[funcName] == nil {
		parser.MakeError(env.parser, expr.StartPos.Line, env.parser.FilePath, expr.Function.StartPos, expr.Function.EndPos, fmt.Sprintf("function '%s' is not defined", funcName)).Display()
	}
	
	function := env.variables[funcName].(FunctionValue)

	args := expr.Args

	params := function.Parameters
	body := function.Body
	returnType := function.ReturnType

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

	if returnType != nil {
		start, end := expr.GetPos()
		if GetRuntimeType(lastVal) != returnType.IType() {
			parser.MakeError(env.parser, start.Line, env.parser.FilePath, start, end, fmt.Sprintf("return type mismatch: expected '%s' but got '%s'", returnType.IType(), GetRuntimeType(lastVal))).Display()
		}
	}

	return lastVal
}


func EvaluateReturnStmt(stmt ast.ReturnStmt, env *Environment) RuntimeValue {
	return Evaluate(stmt.Expression, env)
}