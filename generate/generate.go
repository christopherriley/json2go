package generate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"text/template"
)

const generatedFileTemplate string = `
// {{.Comment}}
// do not modify

package {{.PackageName}}

type {{.StructName}} {{.StructDef}}

var {{.VarName}} {{.StructName}} = {{.Instance}}
`

type TemplateParams struct {
	Comment     string
	PackageName string
	StructName  string
	StructDef   string
	VarName     string
	Instance    string
}

func Generate(comment, rawJson, pkgName, structName, varName string) (string, error) {
	var rawJsonMap map[string]any

	if err := json.Unmarshal([]byte(rawJson), &rawJsonMap); err != nil {
		return "", fmt.Errorf("failed to unmarshal input json: %s", err)
	}

	goStruct := BuildGoStruct(rawJsonMap, structName, 0, "    ")
	instance := GoInstance{goStruct}

	params := TemplateParams{
		Comment:     comment,
		PackageName: pkgName,
		StructName:  structName,
		StructDef:   goStruct.String(),
		VarName:     varName,
		Instance:    instance.String(),
	}

	buf := bytes.Buffer{}
	writer := io.Writer(&buf)
	err := template.Must(template.New("gen").Parse(generatedFileTemplate[1:])).Execute(writer, params)
	if err != nil {
		panic(fmt.Sprintf("template generation failed: %s", err))
	}

	return buf.String(), nil
}
