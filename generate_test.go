package main

import (
	"bytes"
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
				"age": [27, 18, -3],
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
				fmt.Println("length of sample.age is", len(sample.age))
				fmt.Println("length of sample.pref is", len(sample.pref))
				fmt.Println(sample.age[0])
				fmt.Println(sample.age[1])
				fmt.Println(sample.age[2])
			}
`

		const basicJsonMainExpectedOutput = `
			chris
			length of sample.age is 3
			length of sample.pref is 2
			27
			18
			-3
		`

		testGenerate(t, basicJson, basicJsonMain, basicJsonMainExpectedOutput)
	})
}

func testGenerate(t *testing.T, rawJson, mainFunc, expected string) {
	testDir, err := os.MkdirTemp("", "TestGenerate-")
	require.NoError(t, err)
	defer os.RemoveAll(testDir)

	generated := Generate(rawJson, "main", "Sample", "sample")

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
	require.NoError(t, err, buf.String())

	// strip out tabs and spaces from the constant string
	expected = strings.TrimSpace(strings.ReplaceAll(expected, "\t", ""))

	// strip spaces from the output
	actual := strings.TrimSpace(buf.String())

	assert.Equal(t, expected, actual)
}
