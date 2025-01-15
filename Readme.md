# Debezium SMT Go PDK

This library can be used to write
[Debezium SMT](https://debezium.io/documentation/reference/stable/transformations/index.html) in Go.

## Install

Include the library with Go get:

```bash
go get github.com/debezium/debezium-smt-go-pdk
```

## Reference Documentation

You can find the reference documentation for this library on
[pkg.go.dev](https://pkg.go.dev/github.com/debezium/debezium-smt-go-pdk).

## Getting Started

A simple Debezium SMT written in Go should include a `process` function exported like:

```go
package main

import (
	"github.com/debezium/debezium-smt-go-pdk"
)

//export process
func process(proxyPtr uint32) uint32 {
	return debezium.SetNull()
}

func main() {}
```

You can compile the program using [TinyGo](https://tinygo.org/) (version > `0.34.0`):

```bash
tinygo build --no-debug -target=wasm-unknown -o smt.wasm main.go
```

### Data In/Out

For efficiency reasons the full content of the record is not transferred to the Go function, but it can be lazyily accessed using PDK functionalities:

```go
debezium.GetString(debezium.Get(proxyPtr, "value.op"))
```

where `debezium.Get` is used to access the required field with a familiar dot(`.`) syntax, and `debezium.GetString` (or `debezium.IsNull`, `debezium.GetInt32`, etc.) materialize the value.

Similarly, returning a value to Debezium is performed using the PDK functionalities:

```go
return debezium.SetString("foobar")
```

the value returned by the `Set` function (or `SetNull`, `SetBool`, `SetString` ...) should be returned as the result of the `process` function.
