package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
	"walrus/frontend/parser"
	"walrus/utils"
)

func main() {
	// time start

	timeStart := time.Now()

	var files []string;

	files = append(files, "modulesAndImport", "variables", "arrays", "conditionals", "loops", "structsAndTraits");

	for _, file := range files {

		filename := fmt.Sprintf("./../code/%s.wal", file)

		ast := parser.Parse(filename, false)
	
		//store as file
		file, err := os.Create(file + ".json")
	
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
