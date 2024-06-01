package parser

import (
	"fmt"
	"os"
	"walrus/frontend/ast"
	"walrus/frontend/lexer"
	"walrus/utils"
)

func parse_node(p *Parser) ast.Node {

	// can be a statement or an expression
	stmt_fn, exists := stmtLookup[p.currentTokenKind()]

	if exists {
		return stmt_fn(p)
	}

	// if not a statement, then it must be an expression
	expr := parse_expr(p, DEFAULT_BP)

	p.expect(lexer.SEMI_COLON)

	return expr
}
func parse_module_stmt(p *Parser) ast.Statement {
	start := p.currentToken().StartPos

	p.advance() // skip MODULE token

	moduleName := p.expect(lexer.IDENTIFIER).Value

	end := p.expect(lexer.SEMI_COLON).EndPos

	return ast.ModuleStmt{
		BaseStmt: ast.BaseStmt{
			Kind:     ast.MODULE_STATEMENT,
			StartPos: start,
			EndPos:   end,
		},
		ModuleName: moduleName,
	}
}

func parse_import_stmt(p *Parser) ast.Statement {
	start := p.currentToken().StartPos
	//advaced to the next token
	p.advance()

	identifiers := []string{}

	//expect the module name "..." or {x,y,z}
	if p.currentTokenKind() == lexer.OPEN_CURLY {

		p.advance()

		//expect identifiers inside the curly braces
		for p.currentTokenKind() != lexer.CLOSE_CURLY {

			identifier := p.expect(lexer.IDENTIFIER).Value
			identifiers = append(identifiers, identifier)

			//expect a comma between the identifiers
			if p.currentTokenKind() != lexer.CLOSE_CURLY {
				p.expect(lexer.COMMA)
			}
		}
		p.expect(lexer.CLOSE_CURLY)

		//expect the "from" keyword
		p.expect(lexer.FROM)
	}

	moduleName := p.expect(lexer.STRING).Value

	end := p.expect(lexer.SEMI_COLON).EndPos

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

func parse_var_decl_stmt(p *Parser) ast.Statement {

	start := p.currentToken().StartPos

	var explicitType ast.Type
	var assignedValue ast.Expression

	isConstant := p.advance().Kind == lexer.CONST

	//varName := p.expectError(lexer.IDENTIFIER, "Expected identifier after " + (isConstant ? "const" : "let")  ).Value
	errMsg := fmt.Sprintf("Expected identifier after %s", utils.IF(isConstant, "const", "let"))

	varName := p.expectError(lexer.IDENTIFIER, errMsg).Value

	//p.expectError(lexer.COLON, "Expected type or value after variable name")
	if p.currentTokenKind() == lexer.WALRUS {
		p.advance()

		assignedValue = parse_expr(p, DEFAULT_BP)

		if assignedValue == nil {
			p.MakeError(p.currentToken().StartPos.Line, p.FilePath, p.currentToken(), "Expected value after := operator").Display()
		}
	} else if p.currentTokenKind() == lexer.COLON {
		// then we expect type
		p.advance()
		explicitType = parse_type(p, DEFAULT_BP)
		if p.currentTokenKind() == lexer.ASSIGNMENT {
			// then we expect assignment
			p.advance()
			assignedValue = parse_expr(p, DEFAULT_BP)
		}
	} else {
		if p.currentTokenKind() == lexer.ASSIGNMENT {
			p.MakeError(p.currentToken().StartPos.Line, p.FilePath, p.currentToken(), "Invalid token").AddHint("Use ':=' instead\n", TEXT_HINT).Display()
		}
		p.MakeError(p.currentToken().StartPos.Line, p.FilePath, p.currentToken(), "Expected value or type").AddHint("You can declare a variable by\n", TEXT_HINT).AddHint(" let x : i8 = 4;", CODE_HINT).AddHint("\nor,", TEXT_HINT).AddHint("\n let x := 4;", CODE_HINT).Display()
	}

	if isConstant && assignedValue == nil {
		p.MakeError(p.currentToken().StartPos.Line, p.FilePath, p.currentToken(), "Expected value").AddHint("Constants must have a value while declaration", TEXT_HINT).Display()
	}

	end := p.expect(lexer.SEMI_COLON).EndPos

	return ast.VariableDclStml{
		BaseStmt: ast.BaseStmt{
			Kind:     ast.VARIABLE_DECLARATION_STATEMENT,
			StartPos: start,
			EndPos:   end,
		},
		IsConstant:   isConstant,
		Identifier:   varName,
		Value:        assignedValue,
		ExplicitType: explicitType,
	}
}

func parse_block_stmt(p *Parser) ast.Statement {
	return parse_block(p)
}

func parse_block(p *Parser) ast.BlockStmt {

	start := p.expect(lexer.OPEN_CURLY).StartPos

	body := make([]ast.Node, 0)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		body = append(body, parse_node(p))
	}

	end := p.expect(lexer.CLOSE_CURLY).EndPos

	return ast.BlockStmt{
		BaseStmt: ast.BaseStmt{
			Kind:     ast.BLOCK_STATEMENT,
			StartPos: start,
			EndPos:   end,
		},
		Body: body,
	}
}

func parse_function_decl_stmt(p *Parser) ast.Statement {

	start := p.currentToken().StartPos

	p.expect(lexer.FUNCTION)

	functionName := p.expect(lexer.IDENTIFIER).Value
	//parse parameters
	params := parse_params(p)

	// if there is a ARROW token, then we have explicit return type. else we have implicit return type of void
	var explicitReturnType ast.Type
	if p.currentTokenKind() == lexer.ARROW {
		p.advance()
		explicitReturnType = parse_type(p, DEFAULT_BP)
	} else {
		explicitReturnType = ast.VoidType{}
	}

	// parse block
	//type assertion from ast.Statement to ast.BlockStmt
	functionBody := parse_block(p)

	end := functionBody.EndPos

	return ast.FunctionDeclStmt{
		BaseStmt: ast.BaseStmt{
			Kind:     ast.FN_DECLARATION_STATEMENT,
			StartPos: start,
			EndPos:   end,
		},
		FunctionPrototype: ast.FunctionPrototype{
			FunctionName: functionName,
			Parameters:   params,
			ReturnType:   explicitReturnType,
		},
		Block: functionBody,
	}
}

func parse_return_stmt(p *Parser) ast.Statement {

	start := p.currentToken().StartPos

	p.expect(lexer.RETURN)

	var value ast.Expression

	if p.currentTokenKind() != lexer.SEMI_COLON {
		value = parse_expr(p, DEFAULT_BP)
	} else {
		value = ast.VoidExpr{}
	}

	end := p.expect(lexer.SEMI_COLON).EndPos

	return ast.ReturnStmt{
		BaseStmt: ast.BaseStmt{
			Kind:     ast.RETURN_STATEMENT,
			StartPos: start,
			EndPos:   end,
		},
		Expression: value,
	}
}

func parse_break_stmt(p *Parser) ast.Statement {
	return parse_breakout_stmt(p)
}

func parse_continue_stmt(p *Parser) ast.Statement {
	return parse_breakout_stmt(p)
}

func parse_breakout_stmt(p *Parser) ast.Statement {

	start := p.advance().StartPos
	end := p.expect(lexer.SEMI_COLON).EndPos

	return ast.BreakStmt{
		BaseStmt: ast.BaseStmt{
			Kind:     ast.BREAK_STATEMENT,
			StartPos: start,
			EndPos:   end,
		},
	}
}

func parse_struct_decl_stmt(p *Parser) ast.Statement {

	start := p.currentToken().StartPos

	p.expect(lexer.STRUCT)

	properties := map[string]ast.Property{}
	structName := p.expect(lexer.IDENTIFIER).Value
	var embeds []string

	p.expect(lexer.OPEN_CURLY)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {

		var IsStatic bool
		var IsPublic bool
		var ReadOnly bool
		var propname string

		//property
		if p.currentTokenKind() == lexer.ACCESS {

			if p.currentToken().Value == "pub" {
				IsPublic = true
			} else {
				IsPublic = false
			}

			p.advance() //pass the access modifier

			if p.currentTokenKind() == lexer.STATIC {
				IsStatic = true
				p.advance()
			} else {
				IsStatic = false
			}

			if p.currentTokenKind() == lexer.READONLY {
				ReadOnly = true
				p.advance()
			} else {
				ReadOnly = false
			}

			propname = p.expect(lexer.IDENTIFIER).Value

			if p.currentTokenKind() == lexer.COLON {
				//then its a property

				p.advance()

				propertyType := parse_type(p, DEFAULT_BP)

				p.expect(lexer.SEMI_COLON)

				//check if already exists
				if _, exists := properties[propname]; exists {
					panic(fmt.Sprintf("Property %s already declared", propname))
				}

				properties[propname] = ast.Property{
					IsStatic: IsStatic,
					IsPublic: IsPublic,
					ReadOnly: ReadOnly,
					Type:     propertyType,
					//Value: nil,
				}

			}

			continue
		} else if p.currentTokenKind() == lexer.EMBED {
			p.advance()
			//parse the structname to be embeded into this struct
			embededStructName := p.expect(lexer.IDENTIFIER).Value

			embeds = append(embeds, embededStructName)

			p.expect(lexer.SEMI_COLON)
		} else {

			err := "Expected access modifier or embed keyword"

			p.MakeError(p.currentToken().StartPos.Line, p.FilePath, p.currentToken(), err).AddHint("Try adding access modifier - ", TEXT_HINT).AddHint("pub or priv", CODE_HINT).AddHint(" to the property.\n", TEXT_HINT).AddHint("Or,\nTo embed a struct, use the ", TEXT_HINT).AddHint("embed", CODE_HINT).AddHint(" keyword.", TEXT_HINT).Display()

			os.Exit(1)
		}
	}

	end := p.expect(lexer.CLOSE_CURLY).EndPos

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

func parse_trait_decl_stmt(p *Parser) ast.Statement {

	start := p.currentToken().StartPos

	p.advance() //pass the trait token

	traitName := p.expect(lexer.IDENTIFIER).Value

	p.expect(lexer.OPEN_CURLY)

	methods := map[string]ast.Method{}

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {

		//parse access modifier
		isPublic := false
		isStatic := false

		if p.currentTokenKind() == lexer.ACCESS {
			if p.currentToken().Value == "pub" {
				isPublic = true
			}
			p.advance()
		}

		if p.currentTokenKind() == lexer.STATIC {
			isStatic = true
			p.advance()
		}

		// parse the method prototype: fn <name> (params) -> return_type; or fn <name> (params); <- void return type

		method := parse_function_prototype(p)

		traitMethod := ast.Method{
			BaseStmt: ast.BaseStmt{
				Kind:     ast.FN_PROTOTYPE_STATEMENT,
				StartPos: method.StartPos,
				EndPos:   method.EndPos,
			},
			IsPublic: isPublic,
			IsStatic: isStatic,
		}

		methods[method.FunctionName] = traitMethod
	}

	end := p.expect(lexer.CLOSE_CURLY).EndPos

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

func parse_function_prototype(p *Parser) ast.FunctionPrototype {

	start := p.expect(lexer.FUNCTION).StartPos

	Name := p.expect(lexer.IDENTIFIER).Value

	Parameters := parse_params(p)

	var ReturnType ast.Type

	if p.currentTokenKind() == lexer.ARROW {
		p.advance()
		ReturnType = parse_type(p, DEFAULT_BP)
	} else {
		ReturnType = ast.VoidType{}
	}

	end := p.expect(lexer.SEMI_COLON).EndPos

	return ast.FunctionPrototype{
		BaseStmt: ast.BaseStmt{
			Kind:     ast.FN_PROTOTYPE_STATEMENT,
			StartPos: start,
			EndPos:   end,
		},
		FunctionName: Name,
		Parameters:   Parameters,
		ReturnType:   ReturnType,
	}
}

func parse_implement_stmt(p *Parser) ast.Statement {

	//advance impl token
	start := p.advance().StartPos

	var traits []string

	var TypeToImplement string

	// syntax: impl A, B, C for T { ... } or impl A for T { ... } or impl T { ... }
	if p.currentTokenKind() == lexer.IDENTIFIER {
		traits = append(traits, p.expect(lexer.IDENTIFIER).Value)
	}

	//parse the trait names
	for p.currentTokenKind() == lexer.COMMA {
		p.advance()
		traits = append(traits, p.expect(lexer.IDENTIFIER).Value)
	}

	if p.currentTokenKind() != lexer.OPEN_CURLY {
		p.expect(lexer.FOR)
		TypeToImplement = p.expect(lexer.IDENTIFIER).Value
	} else {
		TypeToImplement = traits[0]
	}

	p.expect(lexer.OPEN_CURLY)

	methods := map[string]ast.MethodImplementStmt{}

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {

		start := p.currentToken().StartPos

		isPublic := false
		isStatic := false

		if p.currentTokenKind() == lexer.ACCESS {
			if p.currentToken().Value == "pub" {
				isPublic = true
			}
			p.advance()
		}

		if p.currentTokenKind() == lexer.STATIC {
			isStatic = true
			p.advance()
		}

		method := parse_function_decl_stmt(p).(ast.FunctionDeclStmt)

		methods[method.FunctionName] = ast.MethodImplementStmt{
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

	end := p.expect(lexer.CLOSE_CURLY).EndPos

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

func parse_params(p *Parser) map[string]ast.Type {
	params := map[string]ast.Type{}
	//while )
	p.advance() // pass the open paren

	//parse the parameters
	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_PAREN {

		paramName := p.expect(lexer.IDENTIFIER).Value

		p.expect(lexer.COLON)

		paramType := parse_type(p, DEFAULT_BP)

		//add to the map
		params[paramName] = paramType

		if p.currentTokenKind() != lexer.CLOSE_PAREN {
			p.expect(lexer.COMMA)
		}
	}
	p.expect(lexer.CLOSE_PAREN)
	return params
}

func parse_if_statement(p *Parser) ast.Statement {

	start := p.advance().StartPos

	condition := parse_expr(p, ASSIGNMENT) // using assignment as the lowest binding power

	consequentBlock := parse_block(p)

	var alternate ast.Statement

	if p.currentTokenKind() == lexer.ELSE {
		p.advance() //pass the else
		block := parse_block(p)
		alternate = block
	} else if p.currentTokenKind() == lexer.ELSEIF {
		//p.advance()
		stmt := parse_if_statement(p)
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

func parse_switch_case_stmt(p *Parser) ast.Statement {

    start := p.advance().StartPos // pass the switch token

    discriminant := parse_expr(p, ASSIGNMENT)

    p.expect(lexer.OPEN_CURLY)

    cases := []ast.SwitchCase{}
    var defaultCase *ast.SwitchCase = nil

    for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {

        caseStart := p.currentToken().StartPos
        
        if p.currentTokenKind() == lexer.CASE {
            tests := parse_case_stmt(p)
            block := parse_block(p)

            for _, test := range tests {
                cases = append(cases, ast.SwitchCase{
                    BaseStmt: ast.BaseStmt{
                        Kind:     ast.SWITCH_CASE_STATEMENT,
                        StartPos: caseStart,
                        EndPos:   block.EndPos,
                    },
                    Test:      test,
                    Consequent: block,
                })
            }
        } else if p.currentTokenKind() == lexer.DEFAULT {
            defaultCase = parse_default_case(p)
        } else {
            p.MakeError(p.currentToken().StartPos.Line, p.FilePath, p.currentToken(), "Unexpected token: '" + p.currentToken().Value + "'").AddHint("Switch can have only ", TEXT_HINT).AddHint("case or default", CODE_HINT).AddHint(" keyword", TEXT_HINT).Display()
			panic("-1")
        }
    }

    if defaultCase != nil {
        cases = append(cases, *defaultCase)
    }

    end := p.expect(lexer.CLOSE_CURLY).EndPos

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

func parse_case_stmt(p *Parser) []ast.Expression {
    p.expect(lexer.CASE)

    tests := []ast.Expression{}

    for p.hasTokens() && p.currentTokenKind() != lexer.OPEN_CURLY {
        test := parse_expr(p, ASSIGNMENT)

        tests = append(tests, test)

        if p.currentTokenKind() == lexer.COMMA {
            p.advance()
        }
    }

    return tests
}

func parse_default_case(p *Parser) *ast.SwitchCase {
    defaultStart := p.advance().StartPos // pass the default token
    block := parse_block(p)

    return &ast.SwitchCase{
        BaseStmt: ast.BaseStmt{
            Kind:     ast.DEFAULT_CASE_STATEMENT,
            StartPos: defaultStart,
            EndPos:   block.EndPos,
        },
        Test:      nil,
        Consequent: block,
    }
}



func parse_for_loop_stmt(p *Parser) ast.Statement {

	start := p.currentToken().StartPos

	loopKind := p.advance().Kind

	//parse the init
	identifier := p.expect(lexer.IDENTIFIER).Value

	if loopKind == lexer.FOR {
		p.expect(lexer.WALRUS)

		init := parse_expr(p, ASSIGNMENT)

		p.expect(lexer.SEMI_COLON)

		//parse the condition
		condition := parse_expr(p, ASSIGNMENT)

		p.expect(lexer.SEMI_COLON)

		//parse the post
		post := parse_expr(p, ASSIGNMENT)

		block := parse_block(p)

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

	} else if loopKind == lexer.FOREACH {

		var indexVar string

		/*
			for val, i in arr {
				...
			}
		*/

		if p.currentTokenKind() == lexer.COMMA {
			// then user wants index
			p.advance()

			indexVar = p.expect(lexer.IDENTIFIER).Value
		}

		p.expect(lexer.IN)

		//parse the array

		arr := parse_expr(p, ASSIGNMENT)

		var whereCause ast.Expression

		if p.currentTokenKind() != lexer.OPEN_CURLY {
			p.expect(lexer.WHERE)
			whereCause = parse_expr(p, ASSIGNMENT)
		}

		block := parse_block(p)

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
		p.MakeError(p.currentToken().StartPos.Line, p.FilePath, p.currentToken(), "Expected for or foreach keyword").Display()
		panic("-1")
	}
}

func parse_while_loop_stmt(p *Parser) ast.Statement {

	start := p.advance().StartPos // skip the for token

	cond := parse_expr(p, ASSIGNMENT)

	block := parse_block(p)

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
