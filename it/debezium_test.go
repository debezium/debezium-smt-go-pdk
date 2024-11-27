package it

import (
	"testing"

	"context"
	_ "embed"
	"log"

	"github.com/orsinium-labs/tinytest/is"

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
			return 1
		}).
		Export("get_string").
		NewFunctionBuilder().
		WithFunc(func(v uint32) uint32 {
			return 2
		}).
		Export("get_uint32").
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

func TestNull(t *testing.T) {
	c := is.NewRelaxed(t)

	var ctx = context.Background()
	var r = wazeroStub(ctx)
	defer r.Close(ctx)

	mod, _ := r.Instantiate(ctx, test1Wasm)
	res, _ := mod.ExportedFunction("process").Call(ctx, 0)

	is.Equal(c, res[0], 7)
}

//go:embed testdata/test2.wasm
var test2Wasm []byte

func TestString(t *testing.T) {
	c := is.NewRelaxed(t)

	var ctx = context.Background()
	var r = wazeroStub(ctx)
	defer r.Close(ctx)

	mod, _ := r.Instantiate(ctx, test2Wasm)
	res, _ := mod.ExportedFunction("process").Call(ctx, 0)

	is.Equal(c, res[0], 6)

	buf, _ := mod.Memory().Read(set_string_ptr, 3)
	is.Equal(c, "foo", string(buf))
}
