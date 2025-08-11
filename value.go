package main

import (
	"fmt"
	"math"
	"strings"
)

type GoValue struct {
	typeInfo
	val any
}

func getString(v any) string {
	str, ok := v.(string)
	if !ok {
		panic(fmt.Sprintf("cannot convert '%+v' to string", v))
	}

	str = strings.ReplaceAll(str, "\n", "\\n")
	str = strings.ReplaceAll(str, "\"", `\"`)

	return str
}

func getFloat(v any) float64 {
	f, ok := v.(float64)
	if !ok {
		panic(fmt.Sprintf("cannot convert '%+v' to float64", v))
	}

	return f
}

func getBool(v any) bool {
	b, ok := v.(bool)
	if !ok {
		panic(fmt.Sprintf("cannot convert '%+v' to bool", v))
	}

	return b
}

func getInt(v any) int {
	f := getFloat(v)

	if f != math.Trunc(f) {
		panic(fmt.Sprintf("cannot convert '%f' to int", f))
	}

	return int(f)
}

func getMap(v any) map[string]any {
	m, ok := v.(map[string]any)
	if !ok {
		panic(fmt.Sprintf("cannot convert '%+v' to map", v))
	}

	return m
}

func (v GoValue) String() string {
	var ret string
	var typeInfoStr string

	if v.t == fieldNil {
		return ""
	}

	if v.array {
		value := "{"

		switch v.t {
		case fieldString:
			for _, elem := range v.val.([]any) {
				value += fmt.Sprintf(",\"%s\"", getString(elem))
			}

			value = strings.Replace(value, ",", "", 1) + "}"

		case fieldInt:
			for _, elem := range v.val.([]any) {
				value += fmt.Sprintf(",%d", getInt(elem))
			}

			value = strings.Replace(value, ",", "", 1) + "}"

		case fieldFloat:
			for _, elem := range v.val.([]any) {
				value += fmt.Sprintf(",%f", getFloat(elem))
			}

			value = strings.Replace(value, ",", "", 1) + "}"

		case fieldBool:
			for _, elem := range v.val.([]any) {
				value += fmt.Sprintf(",%t", getBool(elem))
			}

			value = strings.Replace(value, ",", "", 1) + "}"

		case fieldStruct:
			f := NewField("xx", v.val, v.depth+1, v.indent)
			typeInfoStr = f.subStruct.String()
			for _, elem := range v.val.([]any) {
				sub := BuildGoStruct(getMap(elem), "", v.depth+1, v.indent)
				value += "\n" + strings.Repeat(v.indent, v.depth+1)
				value += fmt.Sprintf("%s,", GoInstance{sub}.String())
			}

			value += "\n" + strings.Repeat(v.indent, v.depth) + "}"

		default:
			panic(fmt.Sprintf("field '%s' has unknown type %d", v.fieldName, v.t))
		}

		if v.t != fieldStruct {
			typeInfoStr = v.typeInfo.String()
		}

		ret = fmt.Sprintf("[]%s%s", typeInfoStr, value)
	} else {
		switch v.t {
		case fieldString:
			ret = fmt.Sprintf("\"%s\"", getString(v.val))
		case fieldBool:
			ret = fmt.Sprintf("%t", getBool(v.val))
		case fieldInt:
			ret = fmt.Sprintf("%d", getInt(v.val))
		case fieldFloat:
			ret = fmt.Sprintf("%f", getFloat(v.val))
		case fieldStruct:
			s := BuildGoStruct(getMap(v.val), "", v.depth, v.indent)
			ret = fmt.Sprintf("%s %s", s.String(), GoInstance{s}.String())
		case fieldNil:
		default:
			panic(fmt.Sprintf("field '%s' has unknown type %d", v.fieldName, v.t))
		}
	}

	ret += ","

	return ret
}
