package main

import (
	"encoding/json"
	"fmt"
)

const sampleJson = `
{
    "name": "chris",
    "age": [27, 3, 0],
	"pref": [{
		"color": "blue",
		"food": "pizza"
	},
	{
		"color": "red",
		"food": "ice cream"
	}],
	"awesome": true
}
`

func main() {
	var rawJson map[string]any

	if err := json.Unmarshal([]byte(sampleJson), &rawJson); err != nil {
		panic(err)
	}

	goRep := BuildGoStruct(rawJson, "TestStruct", 0, "    ")
	instance := GoInstance{*goRep}

	fmt.Println(goRep)
	fmt.Println("\n\n\n")
	fmt.Println(instance)
}
