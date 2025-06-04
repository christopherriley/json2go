package main

import (
	"fmt"
	"math"
	"strings"
)

type fieldType int

const (
	fieldString fieldType = iota
	fieldInt
	fieldFloat
	fieldBool
	fieldStruct
)

type GoField struct {
	name      string
	array     bool
	t         fieldType
	subStruct *GoStruct
	val       any
}

func newArrayField(v []any, name string, depth int, indent string) GoField {
	goField := newScalarField(v[0], name, depth, indent)
	goField.array = true
	goField.val = v

	if goField.t == fieldFloat {
		goField.t = fieldInt
	loop:
		for _, elem := range v {
			switch v := elem.(type) {
			case float64:
				if v != math.Trunc(v) {
					goField.t = fieldFloat
					break loop
				}
			default:
				panic("mixed arrays not allowed")
			}
		}
	}

	return goField
}

func newScalarField(v any, name string, depth int, indent string) GoField {
	f := GoField{
		name: name,
		val:  v,
	}

	switch v.(type) {
	case float64:
		f.t = fieldFloat
	case string:
		f.t = fieldString
	case int:
		f.t = fieldInt
	case bool:
		f.t = fieldBool
	default:
		sub, ok := v.(map[string]any)
		if !ok {
			panic("error parsing json")
		}

		f.t = fieldStruct
		f.subStruct = BuildGoStruct(sub, name, depth+1, indent)
	}

	return f
}

func NewField(k string, v any, depth int, indent string) GoField {
	var f GoField

	switch v := v.(type) {
	case []any:
		f = newArrayField(v, k, depth, indent)
	default:
		f = newScalarField(v, k, depth, indent)
	}

	return f
}

func (f GoField) String() string {
	s := f.name + " "

	if f.array {
		s += "[]"
	}

	switch f.t {
	case fieldString:
		s += "string"
	case fieldFloat:
		s += "float64"
	case fieldInt:
		s += "int"
	case fieldBool:
		s += "bool"
	case fieldStruct:
		s += f.subStruct.String()
	default:
		panic(fmt.Sprintf("field '%s' has unknown type %d", f.name, f.t))
	}

	return s
}

func (f GoField) Value() string {
	var ret string

	if f.array {
		var typeName string
		value := "{"

		switch f.t {
		case fieldString:
			typeName = "string"
			for _, elem := range f.val.([]any) {
				value += fmt.Sprintf("\"%s\",", elem.(string))
			}

			value = value[:len(value)-1] + "}"

		case fieldInt:
			typeName = "int"
			for _, elem := range f.val.([]any) {
				value += fmt.Sprintf("%d,", int(elem.(float64)))
			}

			value = value[:len(value)-1] + "}"

		case fieldFloat:
			typeName = "float64"
			for _, elem := range f.val.([]any) {
				value += fmt.Sprintf("%f,", elem.(float64))
			}

			value = value[:len(value)-1] + "}"

		case fieldBool:
			typeName = "bool"
			for _, elem := range f.val.([]any) {
				value += fmt.Sprintf("%t,", elem.(bool))
			}

			value = value[:len(value)-1] + "}"

		case fieldStruct:
			typeName = f.subStruct.String()
			for _, elem := range f.val.([]any) {
				s := BuildGoStruct(elem.(map[string]any), "", f.subStruct.depth+1, f.subStruct.indent)
				value += "\n" + strings.Repeat(f.subStruct.indent, f.subStruct.depth+1)
				value += fmt.Sprintf("%s,", GoInstance{*s}.String())
			}

			value += "\n" + strings.Repeat(f.subStruct.indent, f.subStruct.depth) + "}"

		default:
			panic(fmt.Sprintf("field '%s' has unknown type %d", f.name, f.t))
		}

		ret = fmt.Sprintf("[]%s%s", typeName, value)
	} else {
		switch f.t {
		case fieldString:
			ret = fmt.Sprintf("\"%s\"", f.val.(string))
		case fieldBool:
			ret = fmt.Sprintf("%t", f.val.(bool))
		case fieldInt:
			ret = fmt.Sprintf("%d", int(f.val.(float64)))
		case fieldFloat:
			ret = fmt.Sprintf("%f", f.val.(float64))
		default:
			panic(fmt.Sprintf("field '%s' has unknown type %d", f.name, f.t))
		}
	}

	/*type TestStruct struct {
		name string
		age  []int
		pref []struct {
			color string
			food  string
		}
		awesome bool
	}*/

	/*a := TestStruct{
		name: "chris",
		age:  []int{27, 3, 0},
		pref: []struct {
			color string
			food  string
		}{
			{
				color: "blue",
				food:  "pizza",
			},
			{
				color: "red",
				food:  "ice cream",
			},
		},
		awesome: true,
	}*/

	/*a := TestStruct{
		name: "chris",
		age:  []int{27, 3, 0},
		pref: []struct {
			color string
			food  string
		}{
			{
				color: "blue",
				food:  "pizza",
			},
			{
				color: "red",
				food:  "ice cream",
			}},
		awesome: true,
	}*/

	ret += ","

	return ret
}
