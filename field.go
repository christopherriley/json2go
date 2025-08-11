package main

import (
	"fmt"
	"maps"
	"strings"
)

type GoField struct {
	typeInfo
	subStruct *GoStruct
	val       any
}

var fieldTokenReplace map[string]string

func init() {
	fieldTokenReplace = map[string]string{
		"*":  "_star_",
		"/":  "_fwdslash_",
		"\\": "_backslash_",
		"\"": "_dquote_",
		"'":  "_squote_",
		":":  "_colon_",
		"$":  "_dollar_",
		"(":  "_open_",
		")":  "_close_",
		"{":  "_open_",
		"}":  "_close_",
	}
}

func newEmptyArrayField(v any, name string, depth int, indent string) GoField {
	goField := newScalarField(v, name, depth, indent)
	goField.array = true
	goField.depth = depth
	goField.val = []any{}

	return goField
}

func newArrayField(v []any, name string, depth int, indent string) GoField {
	ti := NewTypeInfo(v, name, indent, depth)

	var subStruct GoStruct

	if ti.t == fieldStruct {
		combinedFields := make(map[string]any)
		for _, elem := range v {
			m, ok := elem.(map[string]any)
			if !ok {
				panic(fmt.Sprintf("error parsing json field '%s' - type is '%T'", name, v))
			}

			maps.Copy(combinedFields, m)
		}

		subStruct = BuildGoStruct(combinedFields, name, depth+1, indent)
	}

	ti.array = true

	goField := GoField{
		typeInfo:  ti,
		subStruct: &subStruct,
	}

	goField.array = true
	goField.depth = depth
	goField.val = v

	return goField
}

func newScalarField(v any, name string, depth int, indent string) GoField {
	ti := NewTypeInfo(v, name, indent, depth)

	var subStruct GoStruct
	if ti.t == fieldStruct {
		sub, ok := v.(map[string]any)
		if !ok {
			panic(fmt.Sprintf("error parsing json field '%s' - type is '%+v'", name, v))
		}

		subStruct = BuildGoStruct(sub, name, depth+1, indent)
	}

	return GoField{
		val:       v,
		typeInfo:  ti,
		subStruct: &subStruct,
	}
}

func sanitizeFieldName(name string) string {
	if name = strings.TrimSpace(name); len(name) == 0 {
		panic("empty field name not allowed")
	}

	for token, replace := range fieldTokenReplace {
		name = strings.ReplaceAll(name, token, replace)
	}

	for i := 0; i < 10; i++ {
		if strings.HasPrefix(name, fmt.Sprintf("%d", i)) {
			name = "_" + name
			break
		}
	}

	name = strings.ToUpper(name[:1]) + name[1:]

	return name
}

func NewField(k string, v any, depth int, indent string) GoField {
	var f GoField

	k = sanitizeFieldName(k)

	switch v := v.(type) {
	case []any:
		if len(v) == 0 {
			f = newEmptyArrayField("", k, depth, indent)
		} else {
			f = newArrayField(v, k, depth, indent)
		}
	default:
		f = newScalarField(v, k, depth, indent)
	}

	return f
}

func (f GoField) String() string {
	s := f.fieldName + " "

	if f.array {
		s += "[]"
	}

	if f.t == fieldStruct {
		s += f.subStruct.String()
	} else {
		s += f.typeInfo.String()
	}

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
