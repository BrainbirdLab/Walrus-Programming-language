package main

import (
	"fmt"
	"os"
	"time"
	"encoding/json"
	"rexlang/lexer"
	"rexlang/parser"
)


func main() {
	// time start

	timeStart := time.Now()

	bytes, err := os.ReadFile("./../examples/03.rex")

	if err != nil {
		panic(err)
	}

	source := string(bytes)

	fmt.Printf("Source code: %s\n", source)

	tokens := lexer.Tokenize(source, true)

	ast := parser.Parse(tokens)

	//store as file
	file, err := os.Create("ast.json");

	if err != nil {
		panic(err)
	}

	//parse as string
	astString, err := json.MarshalIndent(ast, "", "  ")

	if err != nil {
		panic(err)
	}

	_, err = file.Write(astString)

	if err != nil {
		panic(err)
	}

	file.Close()

	// time end
	timeEnd := time.Now()

	fmt.Printf("Time taken: %v\n", timeEnd.Sub(timeStart))
	
}