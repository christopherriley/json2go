package main

import (
	"fmt"
	"strings"
)

type GoInstance struct {
	GoStruct
}

func (i GoInstance) String() string {
	/*a := GoField{
		name: "abc",
		subStruct: &GoStruct{
			name:  "xyz",
			depth: 1,
		},
	}*/

	/*b := GoStruct{
		field: []GoField{
			{
				name: "chris",
			},
			{
				name: "bob",
			},
		},
	}*/

	ret := fmt.Sprintf("%s{\n", i.name)
	for _, f := range i.field {
		ret += strings.Repeat(i.indent, i.depth+1)
		ret += fmt.Sprintf("%s: %s\n", f.name, f.Value())
	}

	ret += strings.Repeat(i.indent, i.depth) + "}"
	return ret
}
