package main

import (

	"os"
	"rexlang/lexer"
	"rexlang/parser"

	"github.com/sanity-io/litter"
	
)


func main() {
	
	bytes, _ := os.ReadFile("./../examples/03.rx")

	source := string(bytes)

	//fmt.Printf("Source code: %s\n", source)

	tokens := lexer.Tokenize(source, true)

	ast := parser.Parse(tokens)

	litter.Dump(ast)
	
}