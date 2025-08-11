package main

import (
	"fmt"
	"sort"
	"strings"
)

type GoStruct struct {
	name   string
	depth  int
	field  []GoField
	indent string
}

func (s GoStruct) String() string {
	ret := "struct {\n"
	for _, f := range s.field {
		ret += strings.Repeat(s.indent, s.depth+1)
		ret += fmt.Sprintf("%s\n", f)
	}
	ret += strings.Repeat(s.indent, s.depth)
	ret += "}"

	return ret
}

func BuildGoStruct(m map[string]any, name string, depth int, indent string) GoStruct {
	s := GoStruct{
		depth:  depth,
		indent: indent,
		name:   name,
	}

	fields := make([]string, 0, len(m))
	for f := range m {
		fields = append(fields, f)
	}
	sort.Strings(fields)

	for _, f := range fields {
		f := NewField(f, m[f], depth, indent)
		s.field = append(s.field, f)
	}

	return s
}
