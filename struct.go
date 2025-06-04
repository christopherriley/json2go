package main

import "fmt"

type GoStruct struct {
	name  string
	field []GoField
}

func (s GoStruct) String() string {
	ret := "struct {\n"
	for _, f := range s.field {
		ret += fmt.Sprintf("%s\n", f)
	}
	ret += "}"

	return ret
}

func BuildGoStruct(m map[string]any) *GoStruct {
	var s GoStruct

	for k, v := range m {
		f := NewField(k, v)
		s.field = append(s.field, f)
	}

	return &s
}
