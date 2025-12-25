package generate

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	t.Run("basic json", func(t *testing.T) {
		const basicJson = `
			{
				"name": "chris",
				"age": [27, 18, -3, 4.2],
				"height": 6,
				"pref": [{
					"color": "blue",
					"food": "pizza"
				},
				{
					"color": "red",
					"food": "ice cream"
				}],
				"awesome": true
			}
`

		const basicJsonMain = `
			package main

			import "fmt"

			func main() {
				fmt.Println(sample.Name)
				fmt.Println(sample.Height)
				fmt.Println("length of sample.age is", len(sample.Age))
				fmt.Println("length of sample.pref is", len(sample.Pref))
				fmt.Println(sample.Age[0])
				fmt.Println(sample.Age[1])
				fmt.Println(sample.Age[2])
				fmt.Println(sample.Age[3])
			}
`

		const basicJsonMainExpectedOutput = `
			chris
			6
			length of sample.age is 4
			length of sample.pref is 2
			27
			18
			-3
			4.2
		`

		testGenerate(t, basicJson, basicJsonMain, basicJsonMainExpectedOutput)
	})

	t.Run("json from example app", func(t *testing.T) {
		const exampleAppJson = `
			{
				"version": "0.56",
				"platform": {
					"windows": {
						"arch": ["amd64", "arm64"],
						"output": "example.exe"
					},
					"linux": {
						"arch": ["amd64"],
						"output": "example_linux"
					},
					"darwin": {
						"arch": ["amd64", "aarch64"],
						"output": "example_macos"
					}
				}
			}
`

		const exampleAppMain = `
			package main

			import "fmt"

			func main() {
				fmt.Println(sample.Version)
			}
`

		const exampleAppMainExpectedOutput = `
			0.56
`

		testGenerate(t, exampleAppJson, exampleAppMain, exampleAppMainExpectedOutput)
	})

	t.Run("empty array", func(t *testing.T) {
		const exampleAppJson = `
			{
				"empty": []
			}
`

		const exampleAppMain = `
			package main

			import "fmt"

			func main() {
				fmt.Println(sample.Empty)
				fmt.Println(len(sample.Empty))
			}
`

		const exampleAppMainExpectedOutput = `
			[]
			0
`

		testGenerate(t, exampleAppJson, exampleAppMain, exampleAppMainExpectedOutput)
	})

	t.Run("array with optional fields", func(t *testing.T) {
		const exampleAppJson = `
				{
					"some_array": [{
						"field_a": 1,
						"field_b": 2
					},
					{
						"field_b": 3,
						"field_c": 4
					}]
				}
	`

		const exampleAppMain = `
				package main

				import "fmt"

				func main() {
					fmt.Println(len(sample.Some_array))
				}
	`

		const exampleAppMainExpectedOutput = `
				2
	`

		testGenerate(t, exampleAppJson, exampleAppMain, exampleAppMainExpectedOutput)
	})

	t.Run("openapi generator response example", func(t *testing.T) {
		// https://github.com/OpenAPITools/openapi-generator-cli/blob/master/examples/v2.0/json/api-with-examples.json

		const exampleAppJson = `
			{
			"swagger": "2.0",
			"info": {
				"title": "Simple API overview",
				"version": "v2"
			},
			"paths": {
				"/": {
					"get": {
						"operationId": "listVersionsv2",
						"summary": "List API versions",
						"produces": [
						"application/json"
						],
						"responses": {
							"200": {
								"description": "200 300 response",
								"examples": {
								"application/json": "{\n    \"versions\": [\n        {\n            \"status\": \"CURRENT\",\n            \"updated\": \"2011-01-21T11:33:21Z\",\n            \"id\": \"v2.0\",\n            \"links\": [\n                {\n                    \"href\": \"http://127.0.0.1:8774/v2/\",\n                    \"rel\": \"self\"\n                }\n            ]\n        },\n        {\n            \"status\": \"EXPERIMENTAL\",\n            \"updated\": \"2013-07-23T11:33:21Z\",\n            \"id\": \"v3.0\",\n            \"links\": [\n                {\n                    \"href\": \"http://127.0.0.1:8774/v3/\",\n                    \"rel\": \"self\"\n                }\n            ]\n        }\n    ]\n}"
								}
							},
							"300": {
								"description": "200 300 response",
								"examples": {
								"application/json": "{\n    \"versions\": [\n        {\n            \"status\": \"CURRENT\",\n            \"updated\": \"2011-01-21T11:33:21Z\",\n            \"id\": \"v2.0\",\n            \"links\": [\n                {\n                    \"href\": \"http://127.0.0.1:8774/v2/\",\n                    \"rel\": \"self\"\n                }\n            ]\n        },\n        {\n            \"status\": \"EXPERIMENTAL\",\n            \"updated\": \"2013-07-23T11:33:21Z\",\n            \"id\": \"v3.0\",\n            \"links\": [\n                {\n                    \"href\": \"http://127.0.0.1:8774/v3/\",\n                    \"rel\": \"self\"\n                }\n            ]\n        }\n    ]\n}"
								}
							}
						}
					}
				},
					"/v2": {
						"get": {
							"operationId": "getVersionDetailsv2",
							"summary": "Show API version details",
							"produces": [
							"application/json"
							],
							"responses": {
							"200": {
								"description": "200 203 response",
								"examples": {
								"application/json": "{\n    \"version\": {\n        \"status\": \"CURRENT\",\n        \"updated\": \"2011-01-21T11:33:21Z\",\n        \"media-types\": [\n            {\n                \"base\": \"application/xml\",\n                \"type\": \"application/vnd.openstack.compute+xml;version=2\"\n            },\n            {\n                \"base\": \"application/json\",\n                \"type\": \"application/vnd.openstack.compute+json;version=2\"\n            }\n        ],\n        \"id\": \"v2.0\",\n        \"links\": [\n            {\n                \"href\": \"http://127.0.0.1:8774/v2/\",\n                \"rel\": \"self\"\n            },\n            {\n                \"href\": \"http://docs.openstack.org/api/openstack-compute/2/os-compute-devguide-2.pdf\",\n                \"type\": \"application/pdf\",\n                \"rel\": \"describedby\"\n            },\n            {\n                \"href\": \"http://docs.openstack.org/api/openstack-compute/2/wadl/os-compute-2.wadl\",\n                \"type\": \"application/vnd.sun.wadl+xml\",\n                \"rel\": \"describedby\"\n            },\n            {\n              \"href\": \"http://docs.openstack.org/api/openstack-compute/2/wadl/os-compute-2.wadl\",\n              \"type\": \"application/vnd.sun.wadl+xml\",\n              \"rel\": \"describedby\"\n            }\n        ]\n    }\n}"
								}
							},
							"203": {
								"description": "200 203 response",
								"examples": {
								"application/json": "{\n    \"version\": {\n        \"status\": \"CURRENT\",\n        \"updated\": \"2011-01-21T11:33:21Z\",\n        \"media-types\": [\n            {\n                \"base\": \"application/xml\",\n                \"type\": \"application/vnd.openstack.compute+xml;version=2\"\n            },\n            {\n                \"base\": \"application/json\",\n                \"type\": \"application/vnd.openstack.compute+json;version=2\"\n            }\n        ],\n        \"id\": \"v2.0\",\n        \"links\": [\n            {\n                \"href\": \"http://23.253.228.211:8774/v2/\",\n                \"rel\": \"self\"\n            },\n            {\n                \"href\": \"http://docs.openstack.org/api/openstack-compute/2/os-compute-devguide-2.pdf\",\n                \"type\": \"application/pdf\",\n                \"rel\": \"describedby\"\n            },\n            {\n                \"href\": \"http://docs.openstack.org/api/openstack-compute/2/wadl/os-compute-2.wadl\",\n                \"type\": \"application/vnd.sun.wadl+xml\",\n                \"rel\": \"describedby\"\n            }\n        ]\n    }\n}"
								}
						}
						}
					}
				}
			},
			"consumes": [
				"application/json"
			]
		}
`

		const exampleAppMain = `
				package main

				import "fmt"

				func main() {
					fmt.Println(sample.Paths._fwdslash_.Get.Summary)
				}
	`

		const exampleAppMainExpectedOutput = `
				List API versions
	`

		testGenerate(t, exampleAppJson, exampleAppMain, exampleAppMainExpectedOutput)
	})
}

func annotatedSource(src string) string {
	ret := ""
	scanner := bufio.NewScanner(strings.NewReader(src))
	line := 1
	for scanner.Scan() {
		ret += fmt.Sprintf("%3d:    %s", line, scanner.Text()) + "\n"
		line++
	}

	return ret
}

func testGenerate(t *testing.T, rawJson, mainFunc, expected string) {
	testDir, err := os.MkdirTemp("", "TestGenerate-")
	require.NoError(t, err)
	defer os.RemoveAll(testDir)

	generated, err := Generate(fmt.Sprintf("generated for unittest - %s", t.Name()), rawJson, "main", "Sample", "sample")
	require.NoError(t, err)

	t.Logf("generated source:\n%s", generated)

	generatedSampleFile := filepath.Join(testDir, "sample.go")
	err = os.WriteFile(generatedSampleFile, []byte(generated), 0644)
	require.NoError(t, err)

	mainFile := filepath.Join(testDir, "main.go")
	err = os.WriteFile(mainFile, []byte(mainFunc), 0644)
	require.NoError(t, err)

	buf := bytes.Buffer{}
	writer := io.Writer(&buf)
	cmd := exec.Command("go", "mod", "init", "unittest")
	cmd.Dir = testDir
	cmd.Stdout = writer
	cmd.Stderr = writer
	err = cmd.Run()
	require.NoError(t, err, buf.String())

	buf = bytes.Buffer{}
	writer = io.Writer(&buf)

	cmd = exec.Command("go", "run", ".")
	cmd.Dir = testDir
	cmd.Stdout = writer
	cmd.Stderr = writer
	err = cmd.Run()
	require.NoError(t, err, buf.String()+"\n\nGENERATED SOURCE FOLLOWS\n\n"+annotatedSource(generated))

	t.Logf("output:\n%s", buf.String())

	// strip out tabs and spaces from the constant string
	expected = strings.TrimSpace(strings.ReplaceAll(expected, "\t", ""))

	// strip spaces from the output
	actual := strings.TrimSpace(buf.String())

	assert.Equal(t, expected, actual)
}
