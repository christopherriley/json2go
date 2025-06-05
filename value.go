package main

import (
	"fmt"
	"strings"
)

type GoValue struct {
	any
	fieldName     string
	fieldTypeName string
	t             fieldType
	array         bool
	indent        string
	depth         int
}

func (v GoValue) String() string {
	var ret string

	if v.array {
		value := "{"

		switch v.t {
		case fieldString:
			for _, elem := range v.any.([]any) {
				value += fmt.Sprintf("\"%s\",", elem.(string))
			}

			value = value[:len(value)-1] + "}"

		case fieldInt:
			for _, elem := range v.any.([]any) {
				value += fmt.Sprintf("%d,", int(elem.(float64)))
			}

			value = value[:len(value)-1] + "}"

		case fieldFloat:
			for _, elem := range v.any.([]any) {
				value += fmt.Sprintf("%f,", elem.(float64))
			}

			value = value[:len(value)-1] + "}"

		case fieldBool:
			for _, elem := range v.any.([]any) {
				value += fmt.Sprintf("%t,", elem.(bool))
			}

			value = value[:len(value)-1] + "}"

		case fieldStruct:
			for _, elem := range v.any.([]any) {
				s := BuildGoStruct(elem.(map[string]any), "", v.depth+1, v.indent)
				value += "\n" + strings.Repeat(v.indent, v.depth+1)
				value += fmt.Sprintf("%s,", GoInstance{*s}.String())
			}

			value += "\n" + strings.Repeat(v.indent, v.depth) + "}"

		default:
			panic(fmt.Sprintf("field '%s' has unknown type %d", v.fieldName, v.t))
		}

		ret = fmt.Sprintf("[]%s%s", v.fieldTypeName, value)
	} else {
		switch v.t {
		case fieldString:
			ret = fmt.Sprintf("\"%s\"", v.any.([]any)[0].(string))
		case fieldBool:
			ret = fmt.Sprintf("%t", v.any.([]any)[0].(bool))
		case fieldInt:
			ret = fmt.Sprintf("%d", int(v.any.([]any)[0].(float64)))
		case fieldFloat:
			ret = fmt.Sprintf("%f", v.any.([]any)[0].(float64))
		default:
			panic(fmt.Sprintf("field '%s' has unknown type %d", v.fieldName, v.t))
		}
	}

	ret += ","

	return ret
}
