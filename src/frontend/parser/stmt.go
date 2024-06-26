package parser

import (
	"fmt"
	"os"
	"walrus/frontend/ast"
	"walrus/frontend/lexer"
	"walrus/utils"
)

func parseNode(p *Parser) ast.Node {

	// can be a statement or an expression
	stmt_fn, exists := stmtLookup[p.currentTokenKind()]

	if exists {
		return stmt_fn(p)
	}

	// if not a statement, then it must be an expression
	expr := parseExpr(p, DEFAULT_BP)

	p.expect(lexer.SEMI_COLON_TOKEN)

	return expr
}
func parseModuleStmt(p *Parser) ast.Statement {
	start := p.currentToken().StartPos

	p.advance() // skip MODULE token

	moduleName := p.expect(lexer.IDENTIFIER_TOKEN).Value

	end := p.expect(lexer.SEMI_COLON_TOKEN).EndPos

	return ast.ModuleStmt{
		BaseStmt: ast.BaseStmt{
			Kind:     ast.MODULE_STATEMENT,
			StartPos: start,
			EndPos:   end,
		},
		ModuleName: moduleName,
	}
}

func parseImportStmt(p *Parser) ast.Statement {
	start := p.currentToken().StartPos
	//advaced to the next token
	p.advance()

	identifiers := []string{}

	//expect the module name "..." or {x,y,z}
	if p.currentTokenKind() == lexer.OPEN_CURLY_TOKEN {

		p.advance()

		//expect identifiers inside the curly braces
		for p.currentTokenKind() != lexer.CLOSE_CURLY_TOKEN {

			identifier := p.expect(lexer.IDENTIFIER_TOKEN).Value
			identifiers = append(identifiers, identifier)

			//expect a comma between the identifiers
			if p.currentTokenKind() != lexer.CLOSE_CURLY_TOKEN {
				p.expect(lexer.COMMA_TOKEN)
			}
		}
		p.expect(lexer.CLOSE_CURLY_TOKEN)

		//expect the "from" keyword
		p.expect(lexer.FROM_TOKEN)
	}

	moduleName := p.expect(lexer.STRING_TOKEN).Value

	end := p.expect(lexer.SEMI_COLON_TOKEN).EndPos

	return ast.ImportStmt{
		BaseStmt: ast.BaseStmt{
			Kind:     ast.IMPORT_STATEMENT,
			StartPos: start,
			EndPos:   end,
		},
		Identifiers: identifiers,
		ModuleName:  moduleName,
	}
}

func parseVarDeclStmt(p *Parser) ast.Statement {

	start := p.currentToken().StartPos

	var explicitType ast.Type
	var assignedValue ast.Expression

	isConstant := p.advance().Kind == lexer.CONST_TOKEN

	//varName := p.expectError(lexer.IDENTIFIER, "Expected identifier after " + (isConstant ? "const" : "let")  ).Value
	errMsg := fmt.Sprintf("Expected identifier after %s", utils.IF(isConstant, "const", "let"))

	variable := p.expectError(lexer.IDENTIFIER_TOKEN, errMsg)

	value := ast.IdentifierExpr{
		BaseStmt: ast.BaseStmt{
			Kind:     ast.IDENTIFIER,
			StartPos: variable.StartPos,
			EndPos:   variable.EndPos,
		},
		Identifier: variable.Value,
	}

	//p.expectError(lexer.COLON, "Expected type or value after variable name")
	if p.currentTokenKind() == lexer.WALRUS_TOKEN {
		p.advance()

		assignedValue = parseExpr(p, DEFAULT_BP)

		if assignedValue == nil {
			MakeError(p, p.currentToken().StartPos.Line, p.FilePath, p.currentToken().StartPos, p.currentToken().EndPos, "Expected value after := operator").Display()
		}
	} else if p.currentTokenKind() == lexer.COLON_TOKEN {
		// then we expect type
		p.advance()
		explicitType = parseType(p, DEFAULT_BP)
		if p.currentTokenKind() == lexer.ASSIGNMENT_TOKEN {
			// then we expect assignment
			p.advance()
			assignedValue = parseExpr(p, DEFAULT_BP)
		}
	} else {
		if p.currentTokenKind() == lexer.ASSIGNMENT_TOKEN {
			MakeError(p, p.currentToken().StartPos.Line, p.FilePath, p.currentToken().StartPos, p.currentToken().EndPos, "Invalid token").AddHint("Use ':=' instead\n", TEXT_HINT).Display()
		}
		MakeError(p, p.currentToken().StartPos.Line, p.FilePath, p.currentToken().StartPos, p.currentToken().EndPos, "Expected value or type").AddHint("You can declare a variable by\n", TEXT_HINT).AddHint(" let x : i8 = 4;", CODE_HINT).AddHint("\nor,", TEXT_HINT).AddHint("\n let x := 4;", CODE_HINT).Display()
	}

	if isConstant && assignedValue == nil {
		MakeError(p, p.currentToken().StartPos.Line, p.FilePath, p.currentToken().StartPos, p.currentToken().EndPos, "Expected value").AddHint("Constants must have a value while declaration", TEXT_HINT).Display()
	}

	end := p.expect(lexer.SEMI_COLON_TOKEN).EndPos

	return ast.VariableDclStml{
		BaseStmt: ast.BaseStmt{
			Kind:     ast.VARIABLE_DECLARATION_STATEMENT,
			StartPos: start,
			EndPos:   end,
		},
		IsConstant:   isConstant,
		Identifier:   value,
		Value:        assignedValue,
		ExplicitType: explicitType,
	}
}

func parseBlockStmt(p *Parser) ast.Statement {
	return parseBlock(p)
}

func parseBlock(p *Parser) ast.BlockStmt {

	start := p.expect(lexer.OPEN_CURLY_TOKEN).StartPos

	body := make([]ast.Node, 0)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY_TOKEN {
		body = append(body, parseNode(p))
	}

	end := p.expect(lexer.CLOSE_CURLY_TOKEN).EndPos

	return ast.BlockStmt{
		BaseStmt: ast.BaseStmt{
			Kind:     ast.BLOCK_STATEMENT,
			StartPos: start,
			EndPos:   end,
		},
		Items: body,
	}
}

func parseFunctionDeclStmt(p *Parser) ast.Statement {

	start := p.currentToken().StartPos

	p.expect(lexer.FUNCTION_TOKEN)

	function := p.expect(lexer.IDENTIFIER_TOKEN)

	functionName := ast.IdentifierExpr{
		BaseStmt: ast.BaseStmt{
			Kind:     ast.IDENTIFIER,
			StartPos: function.StartPos,
			EndPos:   function.EndPos,
		},
		Identifier: function.Value,
	}

	//parse parameters
	params := parseParams(p)

	// if there is a ARROW token, then we have explicit return type. else we have implicit return type of void
	var explicitReturnType ast.Type
	if p.currentTokenKind() == lexer.ARROW_TOKEN {
		p.advance()
		explicitReturnType = parseType(p, DEFAULT_BP)
	} else {
		explicitReturnType = ast.VoidType{
			Kind: ast.T_VOID,
		}
	}

	// parse block
	//type assertion from ast.Statement to ast.BlockStmt
	functionBody := parseBlock(p)

	end := functionBody.EndPos

	return ast.FunctionDeclStmt{
		BaseStmt: ast.BaseStmt{
			Kind:     ast.FN_DECLARATION_STATEMENT,
			StartPos: start,
			EndPos:   end,
		},
		FunctionPrototype: ast.FunctionPrototype{
			BaseStmt: ast.BaseStmt{
				Kind:     ast.FN_PROTOTYPE_STATEMENT,
				StartPos: function.StartPos,
				EndPos:   end,
			},
			Name:       functionName,
			Parameters: params,
			ReturnType: explicitReturnType,
		},
		Block: functionBody,
	}
}

func parseReturnStmt(p *Parser) ast.Statement {

	start := p.currentToken().StartPos

	p.expect(lexer.RETURN_TOKEN)

	var value ast.Expression

	if p.currentTokenKind() != lexer.SEMI_COLON_TOKEN {
		value = parseExpr(p, DEFAULT_BP)
	} else {
		value = ast.VoidLiteral{
			Kind: ast.VOID_LITERAL,
		}
	}

	end := p.expect(lexer.SEMI_COLON_TOKEN).EndPos

	return ast.ReturnStmt{
		BaseStmt: ast.BaseStmt{
			Kind:     ast.RETURN_STATEMENT,
			StartPos: start,
			EndPos:   end,
		},
		Expression: value,
	}
}

func parseBreakStmt(p *Parser) ast.Statement {
	return parseBreakoutStmt(p)
}

func parseContinueStmt(p *Parser) ast.Statement {
	return parseBreakoutStmt(p)
}

func parseBreakoutStmt(p *Parser) ast.Statement {

	start := p.advance().StartPos
	end := p.expect(lexer.SEMI_COLON_TOKEN).EndPos

	return ast.BreakStmt{
		BaseStmt: ast.BaseStmt{
			Kind:     ast.BREAK_STATEMENT,
			StartPos: start,
			EndPos:   end,
		},
	}
}

func parseStructDeclStmt(p *Parser) ast.Statement {

	start := p.currentToken().StartPos

	p.expect(lexer.STRUCT_TOKEN)

	properties := map[string]ast.Property{}
	structName := p.expect(lexer.IDENTIFIER_TOKEN).Value
	var embeds []string

	p.expect(lexer.OPEN_CURLY_TOKEN)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY_TOKEN {

		//property
		if p.currentTokenKind() == lexer.ACCESS_TOKEN {

			setProperties(p, properties)

			continue

		} else if p.currentTokenKind() == lexer.EMBED_TOKEN {
			p.advance()
			//parse the structname to be embeded into this struct
			embededStructName := p.expect(lexer.IDENTIFIER_TOKEN).Value

			embeds = append(embeds, embededStructName)

			p.expect(lexer.SEMI_COLON_TOKEN)
		} else {

			err := "Expected access modifier or embed keyword"

			MakeError(p, p.currentToken().StartPos.Line, p.FilePath, p.currentToken().StartPos, p.currentToken().EndPos, err).AddHint("Try adding access modifier - ", TEXT_HINT).AddHint("pub or priv", CODE_HINT).AddHint(" to the property.\n", TEXT_HINT).AddHint("Or,\nTo embed a struct, use the ", TEXT_HINT).AddHint("embed", CODE_HINT).AddHint(" keyword.", TEXT_HINT).Display()

			os.Exit(1)
		}
	}

	end := p.expect(lexer.CLOSE_CURLY_TOKEN).EndPos

	return ast.StructDeclStatement{
		BaseStmt: ast.BaseStmt{
			Kind:     ast.STRUCT_STATEMENT,
			StartPos: start,
			EndPos:   end,
		},
		Properties: properties,
		Embeds:     embeds,
		StructName: structName,
	}
}

func setProperties(p *Parser, propsMap map[string]ast.Property) {

	var isStatic bool
	var isPublic bool
	var readOnly bool

	if p.currentToken().Value == "pub" {
		isPublic = true
	} else {
		isPublic = false
	}

	p.advance() //pass the access modifier

	if p.currentTokenKind() == lexer.STATIC_TOKEN {
		isStatic = true
		p.advance()
	} else {
		isStatic = false
	}

	if p.currentTokenKind() == lexer.READONLY_TOKEN {
		readOnly = true
		p.advance()
	} else {
		readOnly = false
	}

	prop := p.expect(lexer.IDENTIFIER_TOKEN)

	if p.currentTokenKind() == lexer.COLON_TOKEN {
		//then its a property

		p.advance()

		propertyType := parseType(p, DEFAULT_BP)

		p.expect(lexer.SEMI_COLON_TOKEN)

		//check if already exists
		if _, exists := propsMap[prop.Value]; exists {
			//panic(fmt.Sprintf("Property %s already declared", propName))
			lineNo := prop.StartPos.Line
			filePath := p.FilePath

			errMsg := fmt.Sprintf("Property %s already declared", prop.Value)

			MakeError(p, lineNo, filePath, prop.StartPos, prop.EndPos, errMsg).AddHint("Try removing the duplicate", TEXT_HINT).Display()
		}

		propsMap[prop.Value] = ast.Property{
			BaseStmt: ast.BaseStmt{
				Kind:     ast.STRUCT_PROPERTY,
				StartPos: prop.StartPos,
				EndPos:   prop.EndPos,
			},
			IsStatic: isStatic,
			IsPublic: isPublic,
			ReadOnly: readOnly,
			Name:     prop.Value,
			Type:     propertyType,
		}
	}
}

func parseTraitDeclStmt(p *Parser) ast.Statement {

	start := p.currentToken().StartPos

	p.advance() //pass the trait token

	traitName := p.expect(lexer.IDENTIFIER_TOKEN).Value

	p.expect(lexer.OPEN_CURLY_TOKEN)

	methods := map[string]ast.Method{}

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY_TOKEN {

		//parse access modifier
		isPublic := false
		isStatic := false

		if p.currentTokenKind() == lexer.ACCESS_TOKEN {
			if p.currentToken().Value == "pub" {
				isPublic = true
			}
			p.advance()
		}

		if p.currentTokenKind() == lexer.STATIC_TOKEN {
			isStatic = true
			p.advance()
		}

		// parse the method prototype: fn <name> (params) -> return_type; or fn <name> (params); <- void return type

		method := parseFunctionPrototype(p)

		traitMethod := ast.Method{
			BaseStmt: ast.BaseStmt{
				Kind:     ast.FN_PROTOTYPE_STATEMENT,
				StartPos: method.StartPos,
				EndPos:   method.EndPos,
			},
			IsPublic: isPublic,
			IsStatic: isStatic,
		}

		methods[method.Name.Identifier] = traitMethod
	}

	end := p.expect(lexer.CLOSE_CURLY_TOKEN).EndPos

	return ast.TraitDeclStatement{
		BaseStmt: ast.BaseStmt{
			Kind:     ast.TRAIT_STATEMENT,
			StartPos: start,
			EndPos:   end,
		},
		TraitName: traitName,
		Methods:   methods,
	}
}

func parseFunctionPrototype(p *Parser) ast.FunctionPrototype {

	start := p.expect(lexer.FUNCTION_TOKEN).StartPos

	function := p.expect(lexer.IDENTIFIER_TOKEN)

	Parameters := parseParams(p)

	var ReturnType ast.Type

	if p.currentTokenKind() == lexer.ARROW_TOKEN {
		p.advance()
		ReturnType = parseType(p, DEFAULT_BP)
	} else {
		ReturnType = ast.VoidType{}
	}

	end := p.expect(lexer.SEMI_COLON_TOKEN).EndPos

	return ast.FunctionPrototype{
		BaseStmt: ast.BaseStmt{
			Kind:     ast.FN_PROTOTYPE_STATEMENT,
			StartPos: start,
			EndPos:   end,
		},
		Name: ast.IdentifierExpr{
			BaseStmt: ast.BaseStmt{
				Kind:     ast.IDENTIFIER,
				StartPos: function.StartPos,
				EndPos:   function.EndPos,
			},
			Identifier: function.Value,
		},
		Parameters: Parameters,
		ReturnType: ReturnType,
	}
}

func parseImplementStmt(p *Parser) ast.Statement {

	//advance impl token
	start := p.advance().StartPos

	var traits []string

	var TypeToImplement string

	// syntax: impl A, B, C for T { ... } or impl A for T { ... } or impl T { ... }
	if p.currentTokenKind() == lexer.IDENTIFIER_TOKEN {
		traits = append(traits, p.expect(lexer.IDENTIFIER_TOKEN).Value)
	}

	//parse the trait names
	for p.currentTokenKind() == lexer.COMMA_TOKEN {
		p.advance()
		traits = append(traits, p.expect(lexer.IDENTIFIER_TOKEN).Value)
	}

	if p.currentTokenKind() != lexer.OPEN_CURLY_TOKEN {
		p.expect(lexer.FOR_TOKEN)
		TypeToImplement = p.expect(lexer.IDENTIFIER_TOKEN).Value
	} else {
		TypeToImplement = traits[0]
	}

	p.expect(lexer.OPEN_CURLY_TOKEN)

	methods := map[string]ast.MethodImplementStmt{}

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY_TOKEN {

		start := p.currentToken().StartPos

		isPublic := false
		isStatic := false

		if p.currentTokenKind() == lexer.ACCESS_TOKEN {
			if p.currentToken().Value == "pub" {
				isPublic = true
			}
			p.advance()
		}

		if p.currentTokenKind() == lexer.STATIC_TOKEN {
			isStatic = true
			p.advance()
		}

		method := parseFunctionDeclStmt(p).(ast.FunctionDeclStmt)

		methods[method.Name.Identifier] = ast.MethodImplementStmt{
			BaseStmt: ast.BaseStmt{
				Kind:     ast.FN_DECLARATION_STATEMENT,
				StartPos: start,
				EndPos:   method.EndPos,
			},
			FunctionDeclStmt: method,
			TypeToImplement:  TypeToImplement,
			IsPublic:         isPublic,
			IsStatic:         isStatic,
		}
	}

	end := p.expect(lexer.CLOSE_CURLY_TOKEN).EndPos

	return ast.ImplementStatement{
		BaseStmt: ast.BaseStmt{
			Kind:     ast.IMPLEMENTS_STATEMENT,
			StartPos: start,
			EndPos:   end,
		},
		Impliments: TypeToImplement,
		Traits:     traits,
		Methods:    methods,
	}
}

func parseParams(p *Parser) []ast.FunctionParameter {
	params := []ast.FunctionParameter{}
	//while )
	p.advance() // pass the open paren

	//parse the parameters
	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_PAREN_TOKEN {

		param := p.expect(lexer.IDENTIFIER_TOKEN)

		p.expect(lexer.COLON_TOKEN)

		paramType := parseType(p, DEFAULT_BP)

		//add to the map
		params = append(params, ast.FunctionParameter{
			BaseStmt: ast.BaseStmt{
				Kind:     ast.FUNCTION_PARAMETER,
				StartPos: param.StartPos,
				EndPos:   param.EndPos,
			},
			Identifier: ast.IdentifierExpr{
				BaseStmt: ast.BaseStmt{
					Kind:     ast.IDENTIFIER,
					StartPos: param.StartPos,
					EndPos:   param.EndPos,
				},
				Identifier: param.Value,
			},
			IsVariadic: false,
			Type:       paramType,
			DefaultVal: nil,
		})

		if p.currentTokenKind() != lexer.CLOSE_PAREN_TOKEN {
			p.expect(lexer.COMMA_TOKEN)
		}
	}
	p.expect(lexer.CLOSE_PAREN_TOKEN)
	return params
}

func parseIfStatement(p *Parser) ast.Statement {

	start := p.advance().StartPos

	condition := parseExpr(p, ASSIGNMENT) // using assignment as the lowest binding power

	consequentBlock := parseBlock(p)

	var alternate ast.Statement

	if p.currentTokenKind() == lexer.ELSE_TOKEN {
		p.advance() //pass the else
		block := parseBlock(p)
		alternate = block
	} else if p.currentTokenKind() == lexer.ELSEIF_TOKEN {
		//p.advance()
		stmt := parseIfStatement(p)
		alternate = stmt
	}

	return ast.IfStmt{
		BaseStmt: ast.BaseStmt{
			Kind:     ast.IF_STATEMENT,
			StartPos: start,
			EndPos:   consequentBlock.EndPos,
		},
		Condition: condition,
		Block:     consequentBlock,
		Alternate: alternate,
	}
}

func parseSwitchCaseStmt(p *Parser) ast.Statement {

	start := p.advance().StartPos // pass the switch token

	discriminant := parseExpr(p, ASSIGNMENT)

	p.expect(lexer.OPEN_CURLY_TOKEN)

	cases := []ast.SwitchCase{}
	var defaultCase *ast.SwitchCase = nil

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY_TOKEN {

		caseStart := p.currentToken().StartPos

		if p.currentTokenKind() == lexer.CASE_TOKEN {
			tests := parseCaseStmt(p)
			block := parseBlock(p)

			for _, test := range tests {
				cases = append(cases, ast.SwitchCase{
					BaseStmt: ast.BaseStmt{
						Kind:     ast.SWITCH_CASE_STATEMENT,
						StartPos: caseStart,
						EndPos:   block.EndPos,
					},
					Test:       test,
					Consequent: block,
				})
			}
		} else if p.currentTokenKind() == lexer.DEFAULT_TOKEN {
			defaultCase = parseDefaultCase(p)
		} else {
			MakeError(p, p.currentToken().StartPos.Line, p.FilePath, p.currentToken().StartPos, p.currentToken().EndPos, "Unexpected token: '"+p.currentToken().Value+"'").AddHint("Switch can have only ", TEXT_HINT).AddHint("case or default", CODE_HINT).AddHint(" keyword", TEXT_HINT).Display()
			panic("-1")
		}
	}

	if defaultCase != nil {
		cases = append(cases, *defaultCase)
	}

	end := p.expect(lexer.CLOSE_CURLY_TOKEN).EndPos

	return ast.SwitchStmt{
		BaseStmt: ast.BaseStmt{
			Kind:     ast.SWITCH_STATEMENT,
			StartPos: start,
			EndPos:   end,
		},
		Discriminant: discriminant,
		Cases:        cases,
	}
}

func parseCaseStmt(p *Parser) []ast.Expression {
	p.expect(lexer.CASE_TOKEN)

	tests := []ast.Expression{}

	for p.hasTokens() && p.currentTokenKind() != lexer.OPEN_CURLY_TOKEN {
		test := parseExpr(p, ASSIGNMENT)

		tests = append(tests, test)

		if p.currentTokenKind() == lexer.COMMA_TOKEN {
			p.advance()
		}
	}

	return tests
}

func parseDefaultCase(p *Parser) *ast.SwitchCase {
	defaultStart := p.advance().StartPos // pass the default token
	block := parseBlock(p)

	return &ast.SwitchCase{
		BaseStmt: ast.BaseStmt{
			Kind:     ast.DEFAULT_CASE_STATEMENT,
			StartPos: defaultStart,
			EndPos:   block.EndPos,
		},
		Test:       nil,
		Consequent: block,
	}
}

func parseForLoopStmt(p *Parser) ast.Statement {

	start := p.currentToken().StartPos

	loopKind := p.advance().Kind

	//parse the init
	identifier := p.expect(lexer.IDENTIFIER_TOKEN).Value

	if loopKind == lexer.FOR_TOKEN {
		p.expect(lexer.WALRUS_TOKEN)

		init := parseExpr(p, ASSIGNMENT)

		p.expect(lexer.SEMI_COLON_TOKEN)

		//parse the condition
		condition := parseExpr(p, ASSIGNMENT)

		p.expect(lexer.SEMI_COLON_TOKEN)

		//parse the post
		post := parseExpr(p, ASSIGNMENT)

		block := parseBlock(p)

		end := block.EndPos

		return ast.ForStmt{
			BaseStmt: ast.BaseStmt{
				Kind:     ast.FOR_LOOP_STATEMENT,
				StartPos: start,
				EndPos:   end,
			},
			Variable:  identifier,
			Init:      init,
			Condition: condition,
			Post:      post,
			Block:     block,
		}

	} else if loopKind == lexer.FOREACH_TOKEN {

		var indexVar string

		/*
			for val, i in arr {
				...
			}
		*/

		if p.currentTokenKind() == lexer.COMMA_TOKEN {
			// then user wants index
			p.advance()

			indexVar = p.expect(lexer.IDENTIFIER_TOKEN).Value
		}

		p.expect(lexer.IN_TOKEN)

		//parse the array

		arr := parseExpr(p, ASSIGNMENT)

		var whereCause ast.Expression

		if p.currentTokenKind() != lexer.OPEN_CURLY_TOKEN {
			p.expect(lexer.WHERE_TOKEN)
			whereCause = parseExpr(p, ASSIGNMENT)
		}

		block := parseBlock(p)

		end := block.EndPos

		return ast.ForeachStmt{
			BaseStmt: ast.BaseStmt{
				Kind:     ast.FOREACH_LOOP_STATEMENT,
				StartPos: start,
				EndPos:   end,
			},
			Variable:      identifier,
			IndexVariable: indexVar,
			Iterable:      arr,
			WhereClause:   whereCause,
			Block:         block,
		}

	} else {
		MakeError(p, p.currentToken().StartPos.Line, p.FilePath, p.currentToken().StartPos, p.currentToken().EndPos, "Expected for or foreach keyword").Display()
		panic("-1")
	}
}

func parseWhileLoopStmt(p *Parser) ast.Statement {

	start := p.advance().StartPos // skip the for token

	cond := parseExpr(p, ASSIGNMENT)

	block := parseBlock(p)

	_, end := block.GetPos()

	return ast.WhileLoopStmt{
		BaseStmt: ast.BaseStmt{
			Kind:     ast.WHILE_STATEMENT,
			StartPos: start,
			EndPos:   end,
		},
		Condition: cond,
		Block:     block,
	}
}
