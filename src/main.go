package main

import (
	//"fmt"
	"os"
	"rexlang/lexer"
	"rexlang/parser"

	"github.com/sanity-io/litter"
)

func main() {
	bytes, _ := os.ReadFile("./../examples/02.rx")

	source := string(bytes)

	//fmt.Printf("Source code: %s\n", source)

	tokens := lexer.Tokenize(source)

	ast := parser.Parse(tokens)

	litter.Dump(ast)
}