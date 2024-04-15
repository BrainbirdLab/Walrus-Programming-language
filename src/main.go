package main

import (

	"fmt"
	"os"
	"rexlang/lexer"
	"rexlang/parser"

	"github.com/sanity-io/litter"
	
)


func main() {
	
	bytes, err := os.ReadFile("./../examples/03.x")

	if err != nil {
		panic(err)
	}

	source := string(bytes)

	fmt.Printf("Source code: %s\n", source)

	tokens := lexer.Tokenize(source, true)

	ast := parser.Parse(tokens)

	litter.Dump(ast)
	
}