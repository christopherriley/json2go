package main

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
				fmt.Println(sample.name)
				fmt.Println(sample.height)
				fmt.Println("length of sample.age is", len(sample.age))
				fmt.Println("length of sample.pref is", len(sample.pref))
				fmt.Println(sample.age[0])
				fmt.Println(sample.age[1])
				fmt.Println(sample.age[2])
				fmt.Println(sample.age[3])
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
				fmt.Println(sample.version)
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
				fmt.Println(sample.empty)
				fmt.Println(len(sample.empty))
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
						"field_a": "a1",
						"field_b": "b1"
					},
					{
						"field_b": "b2",
						"field_c": "c2"
					}]
				}
	`

		const exampleAppMain = `
				package main

				import "fmt"

				func main() {
					fmt.Println(len(sample.some_array))
				}
	`

		const exampleAppMainExpectedOutput = `
				2
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

	generated := Generate(fmt.Sprintf("generated for unittest - %s", t.Name()), rawJson, "main", "Sample", "sample")

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
