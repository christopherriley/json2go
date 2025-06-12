package main

import (
	"fmt"
	"math"
	"reflect"
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
	fieldName     string
	fieldTypeName string
	array         bool
	t             fieldType
	depth         int
	indent        string
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
		ti.fieldTypeName = "nil"
		return
	}

	switch v := v.(type) {
	case float64:
		ti.t = fieldFloat
		ti.fieldTypeName = "float64"
		if v == math.Trunc(v) {
			ti.t = fieldInt
			ti.fieldTypeName = "int"
		}
	case string:
		ti.t = fieldString
		ti.fieldTypeName = "string"
	case bool:
		ti.t = fieldBool
		ti.fieldTypeName = "bool"
	default:
		ti.t = fieldStruct
		ti.fieldTypeName = fmt.Sprintf("struct - %s", reflect.TypeOf(v))
	}
}
