package main

import (
	"encoding/json"
	"fmt"
)

const sampleJson = `
{
    "name": "chris",
    "age": 27,
	"pref": {
		"color": "blue",
		"food": "pizza"
	}
}
`

type temp struct {
	a string
	b float64
	c int
	d *temp
}

type aTest struct {
	a string
	b struct {
		c string
		d int
	}
}

type fieldType int

const (
	fieldString fieldType = iota
	fieldInt
	fieldFloat
	fieldStruct
)

type goField struct {
	name      string
	array     bool
	t         fieldType
	subStruct *goStruct
}

type goStruct struct {
	name  string
	field []goField
}

func (s goStruct) String() string {
	ret := "struct {\n"
	for _, f := range s.field {
		ret += fmt.Sprintf("%s\n", f)
	}
	ret += "}"

	return ret
}

func (f goField) String() string {
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
	case fieldStruct:
		s += f.subStruct.String()
	default:
		panic(fmt.Sprintf("field '%s' has unknown type %d", f.name, f.t))
	}

	return s
}

func buildStruct(m map[string]any) *goStruct {
	var s goStruct

	for k, v := range m {
		var f goField

		switch v.(type) {
		case float64:
			f = goField{
				name:  k,
				array: false,
				t:     fieldFloat,
			}
		case string:
			f = goField{
				name:  k,
				array: false,
				t:     fieldString,
			}
		case int:
			f = goField{
				name:  k,
				array: false,
				t:     fieldString,
			}
		default:
			sub, ok := v.(map[string]any)
			if !ok {
				panic("error parsing json")
			}

			f = goField{
				name:      k,
				array:     false,
				t:         fieldStruct,
				subStruct: buildStruct(sub),
			}
		}

		s.field = append(s.field, f)
	}

	return &s
}

func main() {
	var rawJson map[string]any

	if err := json.Unmarshal([]byte(sampleJson), &rawJson); err != nil {
		panic(err)
	}

	goRep := buildStruct(rawJson)

	fmt.Println(goRep)
}
