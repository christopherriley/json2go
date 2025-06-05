package main

type fieldType int

const (
	fieldString fieldType = iota
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
	var ft fieldType
	var ftn string

	ti := typeInfo{
		fieldName:     fieldName,
		fieldTypeName: ftn,
		t:             ft,
		indent:        indent,
		depth:         depth,
	}

	switch v.(type) {
	case float64:
		ti.setType(fieldFloat)
	case string:
		ti.setType(fieldString)
	case int:
		ti.setType(fieldInt)
	case bool:
		ti.setType(fieldBool)
	default:
		ti.setType(fieldStruct)
	}

	return ti
}

func (ti *typeInfo) setType(t fieldType) {
	ti.t = t

	switch t {
	case fieldBool:
		ti.fieldTypeName = "bool"
	case fieldString:
		ti.fieldTypeName = "string"
	case fieldInt:
		ti.fieldTypeName = "int"
	case fieldFloat:
		ti.fieldTypeName = "float64"
	case fieldStruct:
		ti.fieldTypeName = "struct"
	}
}
