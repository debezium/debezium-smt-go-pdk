package it

import (
	"testing"

	"context"
	_ "embed"
	"log"

	"github.com/stretchr/testify/assert"

	"github.com/tetratelabs/wazero"
)

var get_string_ptr uint32
var set_string_ptr uint32

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
		Export("get_bool").
		NewFunctionBuilder().
		WithFunc(func(v uint32) uint32 {
			return 3
		}).
		Export("get_bytes").
		NewFunctionBuilder().
		WithFunc(func(v uint32) uint32 {
			return 4
		}).
		Export("get_float32").
		NewFunctionBuilder().
		WithFunc(func(v uint32) uint64 {
			return 5
		}).
		Export("get_float64").
		NewFunctionBuilder().
		WithFunc(func(v uint32) uint32 {
			return 6
		}).
		Export("get_int16").
		NewFunctionBuilder().
		WithFunc(func(v uint32) uint32 {
			return 7
		}).
		Export("get_int32").
		NewFunctionBuilder().
		WithFunc(func(v1 uint32) uint32 {
			return 8
		}).
		Export("get_int64").
		NewFunctionBuilder().
		WithFunc(func(v uint32) int64 {
			return 9
		}).
		Export("get_int8").
		NewFunctionBuilder().
		WithFunc(func(v1, v2 uint32) uint32 {
			return 10
		}).
		Export("get").
		NewFunctionBuilder().
		WithFunc(func(v uint32) uint32 {
			return 11
		}).
		Export("set_bool").
		NewFunctionBuilder().
		WithFunc(func(v uint32) uint32 {
			set_string_ptr = v
			return 12
		}).
		Export("set_string").
		NewFunctionBuilder().
		WithFunc(func() uint32 {
			return 13
		}).
		Export("set_null").
		NewFunctionBuilder().
		WithFunc(func(v uint32) uint32 {
			return 14
		}).
		Export("is_null").
		NewFunctionBuilder().
		WithFunc(func(v1, v2 uint32) uint32 {
			return 15
		}).
		Export("get_array_elem").
		NewFunctionBuilder().
		WithFunc(func(v uint32) uint32 {
			return 16
		}).
		Export("get_array_size").
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

	mod, err := r.Instantiate(ctx, test1Wasm)
	if err != nil {
		log.Panicln(err)
	}
	res, err := mod.ExportedFunction("process").Call(ctx, 0)
	if err != nil {
		log.Panicln(err)
	}

	assert.Equal(t, uint64(13), res[0])
}

//go:embed testdata/test2.wasm
var test2Wasm []byte

func TestGuestString(t *testing.T) {
	var ctx = context.Background()
	var r = wazeroStub(ctx)
	defer r.Close(ctx)

	mod, _ := r.Instantiate(ctx, test2Wasm)
	res, _ := mod.ExportedFunction("process").Call(ctx, 0)

	assert.Equal(t, uint64(12), res[0])

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

	assert.Equal(t, uint64(12), res[0])

	buf, _ := mod.Memory().Read(set_string_ptr, 3)
	assert.Equal(t, "baz", string(buf))
}
