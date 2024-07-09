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

func nativePrint(args ...typechecker.RuntimeValue) typechecker.RuntimeValue {

	//if no arguments
	if len(args) == 0 {
		fmt.Println()
		return typechecker.MakeVOID()
	}

	for _, arg := range args {
		val, err := typechecker.CastToStringValue(arg)

		if err != nil {
			continue
		}

		//colorize
		fmt.Print(utils.Colorize(utils.YELLOW, val.Value))
	}
	fmt.Println()
	return typechecker.MakeVOID()
}

func nativeTime(args ...typechecker.RuntimeValue) typechecker.RuntimeValue {
	t := time.Now().Unix()
	return typechecker.MakeINT(t, 64, true)
}

func main() {
	// time start

	timeStart := time.Now()

	targetDir := "./../code/test/ret"

	dir, err := os.ReadDir(targetDir)

	if err != nil {
		panic(err)
	}

	fmt.Println("Compiling...")

	for _, file := range dir {

		sf := strings.Split(file.Name(), ".")

		if file.IsDir() || sf[len(sf)-1] != "wal" {
			continue
		}

		filename := targetDir + "/" + file.Name()

		parserMachine := parser.NewParser(filename, false)

		ast := parserMachine.Parse()

		fmt.Printf("Parsed: %v\n", filename)

		//store as file
		file, err := os.Create(targetDir + "/" + sf[0] + ".json")

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

		env.DeclareVariable("true", typechecker.MakeBOOL(true), true)
		env.DeclareVariable("false", typechecker.MakeBOOL(false), true)
		env.DeclareVariable("null", typechecker.MakeNULL(), true)

		env.DeclareNativeFn("print", typechecker.MakeNativeFUNCTION(nativePrint))
		env.DeclareNativeFn("time", typechecker.MakeNativeFUNCTION(nativeTime))

		fmt.Printf("Evaluating: %v\n", filename)

		typechecker.Evaluate(ast, env)
	}

	// time end
	timeEnd := time.Now()

	fmt.Print(utils.Colorize(utils.GREEN, fmt.Sprintf("Compiled succesfully in: %v\n", timeEnd.Sub(timeStart))))
	//wait for user input to close
	fmt.Scanln()
}
