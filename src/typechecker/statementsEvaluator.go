package typechecker

import (
	"fmt"
	"reflect"
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
		case ast.IntegerType:
			explicitSize = t.BitSize
		case ast.FloatingType:
			explicitSize = t.BitSize
		}

		switch v := value.(type) {
		case IntegerValue:
			//modify the Size field of the value. Update the original value
			v.Size = explicitSize // is it reference or copy? It is a copy. So, the original value is not updated
			//update the Original value
			value = v
		case FloatValue:
			v.Size = explicitSize
			value = v
		}


		start, end := stmt.Value.GetPos()
		//check user defined types with the value type
		checkTypes(env.parser, stmt.ExplicitType, value, start, end)
	}

	val, err := env.DeclareVariable(stmt.Identifier.Identifier, value, stmt.IsConstant)

	if err != nil {
		parser.MakeError(env.parser, stmt.StartPos.Line, env.parser.FilePath, stmt.Identifier.StartPos, stmt.Identifier.EndPos, err.Error()).Display()
	}

	return val
}

func checkTypes(p *parser.Parser, explicitType ast.Type, value RuntimeValue, startPos lexer.Position, endPos lexer.Position) {

	var msg string

	switch t := explicitType.(type) {
	case ast.IntegerType:
		if helpers.TypesMatchT[IntegerValue](value) {
			if t.BitSize != value.(IntegerValue).Size {
				msg = fmt.Sprintf("cannot assign integer of size %d to integer of size %d", value.(IntegerValue).Size, t.BitSize)
			}
		} else {
			msg = fmt.Sprintf("cannot assign %s to %s", reflect.TypeOf(value).Name(), "integer")
		}
	case ast.FloatingType:
		if helpers.TypesMatchT[FloatValue](value) {
			if t.BitSize != value.(FloatValue).Size {
				msg = fmt.Sprintf("cannot assign float of size %d to float of size %d", value.(FloatValue).Size, t.BitSize)
			}
		} else {
			msg = fmt.Sprintf("cannot assign %s to %s", reflect.TypeOf(value).Name(), "float")
		}
	case ast.StringType:
		if !helpers.TypesMatchT[StringValue](value) {
			msg = fmt.Sprintf("cannot assign %s to %s", reflect.TypeOf(value).Name(), "string")
		}
	case ast.BooleanType:
		if !helpers.TypesMatchT[BooleanValue](value) {
			msg = fmt.Sprintf("cannot assign %s to %s", reflect.TypeOf(value).Name(), "boolean")
		}
	}

	if msg != "" {
		parser.MakeError(p, startPos.Line, p.FilePath, startPos, endPos, msg).Display()
	}
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
	runtimeVal, err := env.DeclareFunction(stmt.FunctionName.Identifier, stmt.Parameters, stmt.Block)

	if err != nil {
		parser.MakeError(env.parser, stmt.StartPos.Line, env.parser.FilePath, stmt.FunctionName.StartPos, stmt.FunctionName.EndPos, err.Error()).Display()
	}

	return runtimeVal
}