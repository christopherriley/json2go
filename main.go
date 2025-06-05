package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	flagInputFile := flag.String("in", "", "input json file")
	flagOutputFile := flag.String("out", "", "go file to generate")
	flagPackageName := flag.String("package", "main", "the package for the generated code")
	flagStructName := flag.String("struct", "Anonymous", "the name for the generated struct type")
	flagInstanceVar := flag.String("var", "Instance", "the name for the generated instance variable")

	if len(strings.TrimSpace(*flagInputFile)) == 0 {
		fmt.Println("must provide input json file with -in")
		os.Exit(1)
	}

	raw, err := os.ReadFile(*flagInputFile)
	if err != nil {
		fmt.Printf("failed to read input json file: %s\n", err)
		os.Exit(1)
	}

	generatedCode := Generate(string(raw), *flagPackageName, *flagStructName, *flagInstanceVar)

	if len(strings.TrimSpace(*flagOutputFile)) == 0 {
		fmt.Println(generatedCode)
	} else {
		os.WriteFile(*flagOutputFile, []byte(generatedCode), 0644)
	}
}
