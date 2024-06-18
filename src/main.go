package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
	"walrus/frontend/parser"
	"walrus/typechecker"
	"walrus/utils"
)

func main() {
	// time start

	timeStart := time.Now()

	targetDir := "./../code/test"

	dir, err := os.ReadDir(targetDir)

	if err != nil {
		panic(err)
	}

	for _, file := range dir {

		sf := strings.Split(file.Name(), ".")

		if file.IsDir() || sf[len(sf)-1] != "wal"{
			continue
		}

		filename := targetDir + "/" + file.Name()

		parserMachine := parser.NewParser(filename, false)

		ast := parserMachine.Parse()
	
		//store as file
		file, err := os.Create( targetDir + "/" + sf[0] + ".json")
	
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

		env := typechecker.NewEnvironment(nil, parserMachine)

		env.DeclareVariable("true", typechecker.MAKE_BOOL(true), true)
		env.DeclareVariable("false", typechecker.MAKE_BOOL(false), true)
		env.DeclareVariable("null", typechecker.MAKE_NULL(), true)

		result := typechecker.Evaluate(ast, 0, env)

		fmt.Printf("Result: %v\n", result)
	}


	// time end
	timeEnd := time.Now()

	fmt.Print(utils.Colorize(utils.GREEN, fmt.Sprintf("Compiled succesfully in: %v\n", timeEnd.Sub(timeStart))))
}
