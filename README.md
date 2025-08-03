## `json2go - convert json to go code`

This module converts a json blob to go code. The generated code includes both a struct definition and an instance for the input json.

Reading from stdin is supported, as well as reading from input files, making it suitable for use with `go generate`.

### `quick start - generate code from stdin`

try the following example to generate code from stdin:

```bash
cat <<EOF | go run .
{
  "name": "bob",
  "age": 28
}
EOF
```
the generated go source will be written to the console, including both the struct definition and an instance.

### `how to use as a go generator`

Be sure to check out the [example app](example/)

to use this module with `go generate`, use a directive like the following in your go source:

```
//go:generate go run github.com/christopherriley/json2go -in config.json -out config.go -struct Config -var config
```

#### `step 1 - create json file`

create a simple json file to start with, for example

config.json
```json
{
  "name": "bob",
  "age": 28
}
```

#### `step 2 - create a go main`

create a go program that will consume the json file. include the `go generate` directive.

main.go
```go
package main

import "fmt"

//go:generate go run github.com/christopherriley/json2go -in config.json -out config.go -struct Config -var config

func main() {
	// note the instance variable name matches the -var flag to the generate directive, above
	fmt.Println("name: ", config.name)
	fmt.Println("age: ", config.age)
}
```

#### `step 3 - generate the code`

in order for the app to run, the go code will need to be generated from the json

```bash
go generate
```

you can examine the generated source

```bash
cat config.go
```

```go
// this file was generated from config.json
// do not modify

package main

type Config struct {
    age int
    name string
}

var config Config = Config{
    age: 28,
    name: "bob",
}
```

#### `step 4 - run the app`

```bash
go run main.go config.go
```

you should see the following output:

```bash
name: bob
age: 28
```

### `unsupported features`

json/javascript has support for [mixed arrays](https://www.geeksforgeeks.org/types-of-arrays-in-javascript/#mixed-arrays). go can achieve the same thing using `[]any`, but defeats the purpose of this project. therefore json arrays with mixed types are not supported. 
