package main

import (
	"fmt"
	"math"
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
}

func newArrayField(v []any, name string) GoField {
	goField := newScalarField(v[0], name)
	goField.array = true

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

func newScalarField(v any, name string) GoField {
	f := GoField{
		name: name,
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
		f.subStruct = BuildGoStruct(sub)
	}

	return f
}

func NewField(k string, v any) GoField {
	var f GoField

	switch v := v.(type) {
	case []any:
		f = newArrayField(v, k)
	default:
		f = newScalarField(v, k)
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
