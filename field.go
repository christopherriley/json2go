package main

import (
	"math"
)

type GoField struct {
	typeInfo
	subStruct *GoStruct
	val       any
}

func newArrayField(v []any, name string, depth int, indent string) GoField {
	goField := newScalarField(v, name, depth, indent)
	goField.array = true
	goField.depth = depth

	if goField.t == fieldFloat {
		goField.setType(fieldInt)
	loop:
		for _, elem := range v {
			switch v := elem.(type) {
			case float64:
				if v != math.Trunc(v) {
					goField.setType(fieldFloat)
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
	ti := NewTypeInfo(v[0], name, indent, depth)

	var subStruct *GoStruct
	if ti.t == fieldStruct {
		sub, ok := v[0].(map[string]any)
		if !ok {
			panic("error parsing json")
		}

		subStruct = BuildGoStruct(sub, name, depth+1, indent)
		ti.fieldTypeName = subStruct.String()
	}

	return GoField{
		val:       v,
		typeInfo:  ti,
		subStruct: subStruct,
	}
}

func NewField(k string, v any, depth int, indent string) GoField {
	var f GoField

	switch v := v.(type) {
	case []any:
		f = newArrayField(v, k, depth, indent)
	default:
		f = newScalarField([]any{v}, k, depth, indent)
	}

	return f
}

func (f GoField) String() string {
	s := f.fieldName + " "

	if f.array {
		s += "[]"
	}

	s += f.fieldTypeName

	return s
}

func (f GoField) Value() GoValue {
	ret := GoValue{
		val:      f.val,
		typeInfo: f.typeInfo,
	}

	ret.depth = ret.depth + 1

	return ret
}
