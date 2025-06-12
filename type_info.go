package main

import (
	"fmt"
	"math"
)

type fieldType int

const (
	fieldNil fieldType = iota
	fieldString
	fieldInt
	fieldFloat
	fieldBool
	fieldStruct
)

type typeInfo struct {
	fieldName string
	array     bool
	t         fieldType
	depth     int
	indent    string
}

func NewTypeInfo(v any, fieldName, indent string, depth int) typeInfo {
	ti := typeInfo{
		fieldName: fieldName,
		indent:    indent,
		depth:     depth,
	}

	if v == nil {
		ti.setType(nil)
		return ti
	}

	switch v := v.(type) {
	case []any:
		ti.array = true
		ti.setType(v[0])

		if ti.t == fieldInt {
			for _, elem := range v {
				if elem.(float64) != math.Trunc(elem.(float64)) {
					ti.setType(elem)
				}
			}
		}
	default:
		ti.setType(v)
	}

	return ti
}

func (ti *typeInfo) setType(v any) {
	if v == nil {
		ti.t = fieldNil
		return
	}

	switch v := v.(type) {
	case float64:
		ti.t = fieldFloat
		if v == math.Trunc(v) {
			ti.t = fieldInt
		}
	case string:
		ti.t = fieldString
	case bool:
		ti.t = fieldBool
	default:
		ti.t = fieldStruct
	}
}

func (ti typeInfo) String() string {
	switch ti.t {
	case fieldNil:
		return "nil"
	case fieldBool:
		return "bool"
	case fieldInt:
		return "int"
	case fieldFloat:
		return "float64"
	case fieldString:
		return "string"
	case fieldStruct:
		return "struct unknown"
	default:
		panic(fmt.Sprintf("unknown fieldtype %d", ti.t))
	}
}
