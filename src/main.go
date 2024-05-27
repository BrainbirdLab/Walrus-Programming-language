package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
	"walrus/frontend/parser"
)

func main() {
	// time start

	timeStart := time.Now()

	ast := parser.Parse("./../code/05.wal", false)

	//store as file
	file, err := os.Create("ast.json")

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
