package generate

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

func newField(v any, name string, depth int, indent string) GoField {
	ti := NewTypeInfo(v, name, indent, depth)

	var subStruct GoStruct

	if ti.t == fieldStruct {
		switch v := v.(type) {
		case []any:
			combinedFields := make(map[string]any)
			for _, elem := range v {
				m, ok := elem.(map[string]any)
				if !ok {
					panic(fmt.Sprintf("error parsing json field '%s' - type is '%T'", name, v))
				}

				maps.Copy(combinedFields, m)
			}

			subStruct = BuildGoStruct(combinedFields, name, depth+1, indent)

		default:
			sub, ok := v.(map[string]any)
			if !ok {
				panic(fmt.Sprintf("error parsing json field '%s' - type is '%+v'", name, v))
			}

			subStruct = BuildGoStruct(sub, name, depth+1, indent)
		}
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

	for i := range 10 {
		if strings.HasPrefix(name, fmt.Sprintf("%d", i)) {
			name = "_" + name
			break
		}
	}

	name = strings.ToUpper(name[:1]) + name[1:]

	return name
}

func NewField(k string, v any, depth int, indent string) GoField {
	return newField(v, sanitizeFieldName(k), depth, indent)
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
