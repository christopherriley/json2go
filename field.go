package main

import (
	"fmt"
	"math"
)

type GoField struct {
	typeInfo
	subStruct *GoStruct
	val       any
}

func newEmptyArrayField(v any, name string, depth int, indent string) GoField {
	goField := newScalarField(v, name, depth, indent)
	goField.array = true
	goField.depth = depth
	goField.val = []any{}

	return goField
}

func newArrayField(v []any, name string, depth int, indent string) GoField {
	ti := NewTypeInfo(v[0], name, indent, depth)

	if ti.t == fieldStruct {
		var subs []map[string]any
		for _, elem := range v {
			sub, ok := elem.(map[string]any)
			if !ok {
				panic(fmt.Sprintf("error parsing json field '%s' - type is '%+v'", name, v))
			}

			subs = append(subs, sub)
		}

		subStructs := BuildGoStructArray(subs, name, depth+1, indent)

		ti.fieldTypeName = subStructs[0].String()
	}

	ti.array = true

	goField := GoField{
		typeInfo: ti,
	}

	goField.array = true
	goField.depth = depth
	goField.val = v

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

func newScalarField(v any, name string, depth int, indent string) GoField {
	ti := NewTypeInfo(v, name, indent, depth)

	var subStruct *GoStruct
	if ti.t == fieldStruct {
		sub, ok := v.(map[string]any)
		if !ok {
			panic(fmt.Sprintf("error parsing json field '%s' - type is '%+v'", name, v))
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
		if len(v) == 0 {
			f = newEmptyArrayField("", k, depth, indent)
		} else {
			f = newArrayField(v, k, depth, indent)
		}
	default:
		f = newScalarField(v, k, depth, indent)
		if f.t == fieldFloat {
			if f.val.(float64) == math.Trunc(f.val.(float64)) {
				f.setType(fieldInt)
			}
		}
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

func (f GoField) Value() *GoValue {
	if f.val == nil {
		return nil
	}

	ret := GoValue{
		val:      f.val,
		typeInfo: f.typeInfo,
	}

	ret.depth = ret.depth + 1

	return &ret
}
