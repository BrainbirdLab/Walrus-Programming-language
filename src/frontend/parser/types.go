package parser

import (
	"fmt"
	"rexlang/frontend/ast"
	"rexlang/frontend/lexer"
	"rexlang/utils"
)

// Null denoted. Expect nothing to the left of the token
type type_nud_handler func(p *Parser) ast.Type

// Left denoted. Expect something to the left of the token
type type_led_handler func(p *Parser, left ast.Type, bp binding_power) ast.Type

type type_nud_lookup map[lexer.TokenKind]type_nud_handler
type type_led_lookup map[lexer.TokenKind]type_led_handler
type type_bp_lookup map[lexer.TokenKind]binding_power

var type_bp_lu = type_bp_lookup{}
var type_nud_lu = type_nud_lookup{}
var type_led_lu = type_led_lookup{}

func type_led(kind lexer.TokenKind, bp binding_power, led_fn type_led_handler) {
	type_bp_lu[kind] = bp
	type_led_lu[kind] = led_fn
}

func type_nud(kind lexer.TokenKind, nud_fn type_nud_handler) {
	type_nud_lu[kind] = nud_fn
}

func createTokenTypesLookups() {
	type_nud(lexer.IDENTIFIER, parse_data_type)
	type_nud(lexer.OPEN_BRACKET, parse_array_type)
}

func parse_data_type(p *Parser) ast.Type {
	value := p.expect(lexer.IDENTIFIER).Value
	switch value {
		case "i8","i16","i32","i64","i128":
			return ast.IntegerType{
				Kind: ast.INTEGER,
				BitSize: utils.BitSizeFromString(value),
				IsSigned: true,
			}
		case "u8","u16","u32","u64","u128":
			return ast.IntegerType{
				Kind: ast.INTEGER,
				BitSize: utils.BitSizeFromString(value),
				IsSigned: false,
			}
		case "f32", "f64":
			return ast.FloatingType{
				Kind: ast.FLOATING,
				BitSize: utils.BitSizeFromString(value),
			}
		case "bool":
			return ast.BooleanType{
				Kind: ast.BOOLEAN,
			}
		case "char":
			return ast.CharecterType{
				Kind: ast.CHARACTER,
			}
		case "str":
			return ast.StringType{
				Kind: ast.STRING,
			}
		default:
			return ast.NullType{
				Kind: ast.NULL,
			}
	}
}


func parse_array_type(p *Parser) ast.Type {

	p.advance()
	p.expect(lexer.CLOSE_BRACKET)

	elemType := parse_type(p, default_bp)

	return ast.ArrayType{
		ElementType: elemType,
	}
}

func parse_type(p *Parser, bp binding_power) ast.Type {
	// Fist parse the NUD
	tokenKind := p.currentTokenKind()
	nud_fn, exists := type_nud_lu[tokenKind]

	if !exists {
		panic(fmt.Sprintf("TYPE NUD handler expected for token %s\n", lexer.TokenKindString(tokenKind)))
		//return nil
	}

	left := nud_fn(p)

	for type_bp_lu[p.currentTokenKind()] > bp {

		tokenKind := p.currentTokenKind()

		led_fn, exists := type_led_lu[tokenKind]

		if !exists {
			panic(fmt.Sprintf("TYPE LED handler expected for token %s\n", lexer.TokenKindString(tokenKind)))
		}

		left = led_fn(p, left, type_bp_lu[p.currentTokenKind()])
	}

	return left
}
