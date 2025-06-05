package main

import "encoding/json"

func Generate(rawJson, pkgName, structName, varName string) string {
	var rawJsonMap map[string]any

	if err := json.Unmarshal([]byte(rawJson), &rawJsonMap); err != nil {
		panic(err)
	}

	goStruct := BuildGoStruct(rawJsonMap, structName, 0, "    ")
	instance := GoInstance{*goStruct}

	var ret string

	ret += "// generated file - do not modify" + "\n"
	ret += "\n"
	ret += "package " + pkgName + "\n"
	ret += "\n"
	ret += "type " + structName + " " + goStruct.String() + "\n"
	ret += "\n"
	ret += "var " + varName + " " + structName + " = " + instance.String() + "\n"

	return ret
}
