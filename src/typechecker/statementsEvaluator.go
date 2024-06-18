package typechecker

import (
	"walrus/frontend/ast"
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

	val, err := env.DeclareVariable(stmt.Identifier.Identifier, value, stmt.IsConstant)

	if err != nil {
		parser.MakeError(env.parser, stmt.StartPos.Line, env.parser.FilePath, stmt.Identifier.StartPos, stmt.Identifier.EndPos, err.Error()).Display()
	}

	return val
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