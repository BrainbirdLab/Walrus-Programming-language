package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
	"walrus/frontend/parser"
	"walrus/utils"
)

func main() {
	// time start

	timeStart := time.Now()

	targetDir := "./../code"

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

		ast := parser.Parse(filename, false)
	
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
	}


	// time end
	timeEnd := time.Now()

	fmt.Print(utils.Colorize(utils.GREEN, fmt.Sprintf("Compiled succesfully in: %v\n", timeEnd.Sub(timeStart))))
}
