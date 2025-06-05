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
	typeName  string
	array     bool
	t         fieldType
	subStruct *GoStruct
	val       any
	depth     int
	indent    string
}

func newArrayField(v []any, name string, depth int, indent string) GoField {
	goField := newScalarField(v, name, depth, indent)
	goField.array = true
	goField.depth = depth

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

func newScalarField(v []any, name string, depth int, indent string) GoField {
	f := GoField{
		name:   name,
		val:    v,
		indent: indent,
		depth:  depth,
	}

	switch v[0].(type) {
	case float64:
		f.t = fieldFloat
	case string:
		f.t = fieldString
	case int:
		f.t = fieldInt
	case bool:
		f.t = fieldBool
	default:
		sub, ok := v[0].(map[string]any)
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
		f = newScalarField([]any{v}, k, depth, indent)
	}

	switch f.t {
	case fieldString:
		f.typeName = "string"
	case fieldFloat:
		f.typeName = "float64"
	case fieldInt:
		f.typeName = "int"
	case fieldBool:
		f.typeName = "bool"
	case fieldStruct:
		f.typeName = f.subStruct.String()
	default:
		panic(fmt.Sprintf("field '%s' has unknown type %d", f.name, f.t))
	}

	return f
}

func (f GoField) String() string {
	s := f.name + " "

	if f.array {
		s += "[]"
	}

	s += f.typeName

	return s
}

func (f GoField) Value() GoValue {
	ret := GoValue{
		any:           f.val,
		fieldName:     f.name,
		fieldTypeName: f.typeName,
		t:             f.t,
		array:         f.array,
		indent:        f.indent,
		depth:         f.depth + 1,
	}

	return ret
}
