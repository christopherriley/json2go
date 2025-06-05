package main

// example app that demonstrates code generation from json

// how to try this example
// - go generate
// - go run .

// explanation of go generate directive:
// - input file is config.json, from current directory
// - config.go will be generated from the input json file
// - the struct name in the generated go file will be called 'Config'
// - the instance of the struct in the generated go file will be called 'config'

//go:generate go run github.com/christopherriley/json2go -in config.json -out config.go -struct Config -var config

import "fmt"

func main() {
	fmt.Println("config version: ", config.version)

	fmt.Println()

	fmt.Println("config windows output name: ", config.platform.windows.output)
	fmt.Println("config windows arch: ", config.platform.windows.arch)

	fmt.Println()

	fmt.Println("config linux output name: ", config.platform.linux.output)
	fmt.Println("config linux arch: ", config.platform.linux.arch)

	fmt.Println()

	fmt.Println("config darwin output name: ", config.platform.darwin.output)
	fmt.Println("config darwin arch: ", config.platform.darwin.arch)

	fmt.Println()
}
