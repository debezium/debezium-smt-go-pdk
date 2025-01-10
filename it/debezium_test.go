package it

import (
	"testing"

	"context"
	_ "embed"
	"log"
	"strconv"

	"github.com/stretchr/testify/assert"

	"github.com/tetratelabs/wazero"
)

var get_string_ptr uint32
var set_string_ptr uint32
var set_int uint32

// wazero module builder
func wazeroStub(ctx context.Context) wazero.Runtime {
	var r = wazero.NewRuntime(ctx)

	var _, err = r.NewHostModuleBuilder("env").
		NewFunctionBuilder().
		WithFunc(func(v uint32) uint32 {
			get_string_ptr = v
			return get_string_ptr
		}).
		Export("get_string").
		NewFunctionBuilder().
		WithFunc(func(v uint32) uint32 {
			return 2
		}).
		Export("get_int").
		NewFunctionBuilder().
		WithFunc(func(v1, v2 uint32) uint32 {
			return 3
		}).
		Export("get").
		NewFunctionBuilder().
		WithFunc(func(v uint32) uint32 {
			return 4
		}).
		Export("set_bool").
		NewFunctionBuilder().
		WithFunc(func(v uint32) uint32 {
			set_int = v
			return 5
		}).
		Export("set_int").
		NewFunctionBuilder().
		WithFunc(func(v uint32) uint32 {
			set_string_ptr = v
			return 6
		}).
		Export("set_string").
		NewFunctionBuilder().
		WithFunc(func() uint32 {
			return 7
		}).
		Export("set_null").
		NewFunctionBuilder().
		WithFunc(func(v uint32) uint32 {
			return 8
		}).
		Export("is_null").
		Instantiate(ctx)

	if err != nil {
		log.Panicln(err)
	}

	return r
}

//go:embed testdata/test1.wasm
var test1Wasm []byte

func TestGuestNull(t *testing.T) {
	var ctx = context.Background()
	var r = wazeroStub(ctx)
	defer r.Close(ctx)

	mod, _ := r.Instantiate(ctx, test1Wasm)
	res, _ := mod.ExportedFunction("process").Call(ctx, 0)

	assert.Equal(t, res[0], uint64(7))
}

//go:embed testdata/test2.wasm
var test2Wasm []byte

func TestGuestString(t *testing.T) {
	var ctx = context.Background()
	var r = wazeroStub(ctx)
	defer r.Close(ctx)

	mod, _ := r.Instantiate(ctx, test2Wasm)
	res, _ := mod.ExportedFunction("process").Call(ctx, 0)

	assert.Equal(t, res[0], uint64(6))

	buf, _ := mod.Memory().Read(set_string_ptr, 3)
	assert.Equal(t, "foo", string(buf))
}

//go:embed testdata/test3.wasm
var test3Wasm []byte

func TestHostString(t *testing.T) {
	var ctx = context.Background()
	var r = wazeroStub(ctx)
	defer r.Close(ctx)

	mod, _ := r.Instantiate(ctx, test3Wasm)

	malloc := mod.ExportedFunction("malloc")
	results, _ := malloc.Call(ctx, 4)
	namePtr := results[0]
	mod.Memory().Write(uint32(namePtr), []byte("baz"))

	res, _ := mod.ExportedFunction("process").Call(ctx, namePtr)

	assert.Equal(t, res[0], uint64(6))

	buf, _ := mod.Memory().Read(set_string_ptr, 3)
	assert.Equal(t, "baz", string(buf))
}

//go:embed testdata/test4.wasm
var test4Wasm []byte

func TestGuestNumbers(t *testing.T) {
	var ctx = context.Background()
	var r = wazeroStub(ctx)
	defer r.Close(ctx)

	mod, _ := r.Instantiate(ctx, test4Wasm)

	res, _ := mod.ExportedFunction("process").Call(ctx, 123)

	assert.Equal(t, res[0], uint64(5))

	buf, _ := mod.Memory().Read(set_int, 3)
	myInt, _ := strconv.Atoi(string(buf))

	assert.Equal(t, myInt, 123)
}
