package main

import (
	"fmt"
	"os"
	"rexlang/lexer"
)

func main() {
	bytes, _ := os.ReadFile("./../examples/00.rx")

	source := string(bytes)

	fmt.Printf("Source code: %s\n", source)

	tokens := lexer.Tokenize(source)

	for _, token := range tokens {
		token.Debug()
	}
}