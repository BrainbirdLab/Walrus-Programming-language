package parser

import (
	"fmt"
	"os"
	"walrus/frontend/ast"
	"walrus/frontend/lexer"
	"walrus/utils"
)

// Null denoted. Expect nothing to the left of the token
type typeNudNandlerType func(p *Parser) ast.Type

// Left denoted. Expect something to the left of the token
type typeLedHandlerType func(p *Parser, left ast.Type, bp BINDING_POWER) ast.Type

type typeNudLookupType map[lexer.TOKEN_KIND]typeNudNandlerType
type typeLedLookupType map[lexer.TOKEN_KIND]typeLedHandlerType
type typeBpLookupType map[lexer.TOKEN_KIND]BINDING_POWER

var typeBindindLookup = typeBpLookupType{}
var typeNudLookup = typeNudLookupType{}
var typeLedLookup = typeLedLookupType{}

func typeNUD(kind lexer.TOKEN_KIND, handleTypeNud typeNudNandlerType) {
	typeNudLookup[kind] = handleTypeNud
}

func createTokenTypesLookups() {
	typeNUD(lexer.IDENTIFIER_TOKEN, parseDataType)
	typeNUD(lexer.OPEN_BRACKET_TOKEN, parseArrayType)
}

func parseDataType(p *Parser) ast.Type {
	identifier := p.expect(lexer.IDENTIFIER_TOKEN)

	value := identifier.Value

	switch value {
	case "i8", "i16", "i32", "i64", "i128":
		return ast.IntegerType{
			Kind:     ast.DATA_TYPE(value),
			BitSize:  utils.BitSizeFromString(value),
			IsSigned: true,
		}
	case "u8", "u16", "u32", "u64", "u128":
		return ast.IntegerType{
			Kind:     ast.DATA_TYPE(value),
			BitSize:  utils.BitSizeFromString(value),
			IsSigned: false,
		}
	case "f32", "f64":
		return ast.FloatType{
			Kind:    ast.DATA_TYPE(value),
			BitSize: utils.BitSizeFromString(value),
		}
	case "bool":
		return ast.BoolType{
			Kind: ast.T_BOOLEAN,
		}
	case "chr":
		return ast.CharType{
			Kind: ast.T_CHARACTER,
		}
	case "str":
		return ast.StringType{
			Kind: ast.T_STRING,
		}
	default:
		return ast.StructType{
			Kind: ast.T_STRUCT,
			Name: value,
		}
		/*
			p.MakeError(identifier.StartPos.Line, p.FilePath, identifier, fmt.Sprintf("Unknown data type '%s'\n", value)).AddHint("You can use primitives types like i8, i16, i32, i64, i128, u8, u16, u32, u64, u128, f32, f64, bool, char, str, or arrays of them").Display()
			panic("Error while parsing")
		*/
	}
}

func parseArrayType(p *Parser) ast.Type {

	p.advance()
	p.expect(lexer.CLOSE_BRACKET_TOKEN)

	elemType := parseType(p, DEFAULT_BP)

	return ast.ArrayType{
		ElementType: elemType,
	}
}

func parseType(p *Parser, bp BINDING_POWER) ast.Type {
	// Fist parse the NUD
	tokenKind := p.currentTokenKind()
	nudFunction, exists := typeNudLookup[tokenKind]

	if !exists {
		//panic(fmt.Sprintf("TYPE NUD handler expected for token %s\n", tokenKind))
		err := MakeError(p, p.currentToken().StartPos.Line, p.FilePath, p.currentToken().StartPos, p.currentToken().EndPos, fmt.Sprintf("Unexpected token %s\n", tokenKind))

		err.AddHint("Follow ", TEXT_HINT)
		err.AddHint("let x := 10", CODE_HINT)
		err.AddHint(" syntax or", TEXT_HINT)
		err.AddHint("Use primitive types like ", TEXT_HINT)
		err.AddHint("i8, i16, i32, i64, i128, u8, u16, u32, u64, u128, f32, f64, bool, char, str", CODE_HINT)
		err.AddHint(" or arrays of them", TEXT_HINT)
		err.Display()

		os.Exit(1)
		//return nil
	}

	left := nudFunction(p)

	for typeBindindLookup[p.currentTokenKind()] > bp {

		tokenKind := p.currentTokenKind()

		ledFunction, exists := typeLedLookup[tokenKind]

		if !exists {
			panic(fmt.Sprintf("TYPE LED handler expected for token %s\n", tokenKind))
		}

		left = ledFunction(p, left, typeBindindLookup[p.currentTokenKind()])
	}

	return left
}
