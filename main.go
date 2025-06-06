package main

import (
	"bufio"
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

	flag.Parse()

	var input, generatedFileComment string

	if len(strings.TrimSpace(*flagInputFile)) == 0 {
		os.Stderr.WriteString("input file was not specified - reading from stdin\n\n")
		for scanner := bufio.NewScanner(os.Stdin); scanner.Scan(); {
			input += scanner.Text()
		}

		generatedFileComment = "this file was generated"
	} else {
		rawBytes, err := os.ReadFile(*flagInputFile)
		if err != nil {
			fmt.Printf("failed to read input json file: %s\n", err)
			os.Exit(1)
		}

		input = string(rawBytes)
		generatedFileComment = fmt.Sprintf("this file was generated from %s", *flagInputFile)
	}

	generatedCode := Generate(generatedFileComment, input, *flagPackageName, *flagStructName, *flagInstanceVar)

	if len(strings.TrimSpace(*flagOutputFile)) == 0 {
		fmt.Println(generatedCode)
	} else {
		os.WriteFile(*flagOutputFile, []byte(generatedCode), 0644)
	}
}
