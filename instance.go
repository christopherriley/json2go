package main

import (
	"fmt"
	"strings"
)

type GoInstance struct {
	GoStruct
}

func (i GoInstance) String() string {
	ret := fmt.Sprintf("%s{\n", i.name)
	for _, f := range i.field {
		if f.Value() != nil {
			ret += strings.Repeat(i.indent, i.depth+1)
			ret += fmt.Sprintf("%s: %s\n", f.fieldName, f.Value())
		}
	}

	ret += strings.Repeat(i.indent, i.depth) + "}"
	return ret
}
