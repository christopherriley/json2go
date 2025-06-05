package main

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
