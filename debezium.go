package debezium

import (
	"unsafe"

	"github.com/debezium/debezium-smt-go-pdk/internal"
)

// access a nested field in the record data structure provided by Debezium
func Get(proxyPtr uint32, fieldName string) uint32 {
	var fieldNameLen = len(fieldName) + 1
	var fieldNamePtr = internal.Malloc(uintptr(fieldNameLen))
	internal.WriteCString(uintptr(fieldNamePtr), fieldName)

	return envGet(proxyPtr, uint32(uintptr(fieldNamePtr)))
}

// access the elemnt at index of the array
func GetArrayElem(proxyPtr uint32, elementIdx uint32) uint32 {
	return envGetArrayElem(proxyPtr, uint32(elementIdx))
}

// returns the dimension of the array
func GetArraySize(proxyPtr uint32) uint32 {
	return envGetArraySize(proxyPtr)
}

// materialize the String content referenced
func GetString(proxyPtr uint32) string {
	var resultPtr = envGetString(proxyPtr)
	var result = internal.ReadCString(resultPtr)
	internal.Free(unsafe.Pointer(uintptr(resultPtr)))
	return result
}

// materialize the Schema Name
func GetSchemaName(proxyPtr uint32) string {
	var resultPtr = envGetSchemaName(proxyPtr)
	var result = internal.ReadCString(resultPtr)
	internal.Free(unsafe.Pointer(uintptr(resultPtr)))
	return result
}

// materialize the Schema Type
func GetSchemaType(proxyPtr uint32) string {
	var resultPtr = envGetSchemaType(proxyPtr)
	var result = internal.ReadCString(resultPtr)
	internal.Free(unsafe.Pointer(uintptr(resultPtr)))
	return result
}

// materialize the Boolean content referenced
func GetBool(proxyPtr uint32) bool {
	return envGetBool(proxyPtr) > 0
}

// materialize the Bytes content referenced
func GetBytes(proxyPtr uint32) []byte {
	var resultPtr = envGetBytes(proxyPtr)
	var result = internal.ReadCString(resultPtr)
	internal.Free(unsafe.Pointer(uintptr(resultPtr)))
	return []byte(result)
}

// materialize the Float32 content referenced
func GetFloat32(proxyPtr uint32) float32 {
	return envGetFloat32(proxyPtr)
}

// materialize the Float64 content referenced
func GetFloat64(proxyPtr uint32) float64 {
	return envGetFloat64(proxyPtr)
}

// materialize the Int16 content referenced
func GetInt16(proxyPtr uint32) int16 {
	return int16(envGetInt16(proxyPtr))
}

// materialize the Int32 content referenced
func GetInt32(proxyPtr uint32) int32 {
	return envGetInt32(proxyPtr)
}

// materialize the Int64 content referenced
func GetInt64(proxyPtr uint32) int64 {
	return envGetInt64(proxyPtr)
}

// materialize the Int8 content referenced
func GetInt8(proxyPtr uint32) int8 {
	return int8(envGetInt8(proxyPtr))
}

// check whenever the referenced content is Null
func IsNull(valuePtr uint32) bool {
	return (envIsNull(valuePtr) > 0)
}

// set a Boolean content for the Debezium Host
func SetBool(value bool) uint32 {
	var valuePtr = internal.Malloc(1)

	if value {
		*(*byte)(unsafe.Pointer(uintptr(valuePtr))) = 1
	} else {
		*(*byte)(unsafe.Pointer(uintptr(valuePtr))) = 0
	}

	return envSetBool(uint32(uintptr(valuePtr)))
}

// set a Null content for the Debezium Host
func SetNull() uint32 {
	return envSetNull()
}

// set a String content for the Debezium Host
func SetString(value string) uint32 {
	var valueLen = len(value) + 1
	var valuePtr = internal.Malloc(uintptr(valueLen))
	internal.WriteCString(uintptr(valuePtr), value)
	return envSetString(uint32(uintptr(valuePtr)))
}

//go:wasm-module env
//export get_string
func envGetString(proxyPtr uint32) uint32

//go:wasm-module env
//export get_schema_name
func envGetSchemaName(proxyPtr uint32) uint32

//go:wasm-module env
//export get_schema_type
func envGetSchemaType(proxyPtr uint32) uint32

//go:wasm-module env
//export set_null
func envSetNull() uint32

//go:wasm-module env
//export is_null
func envIsNull(valuePtr uint32) uint32

//go:wasm-module env
//export set_string
func envSetString(valuePtr uint32) uint32

//go:wasm-module env
//export set_bool
func envSetBool(valuePtr uint32) uint32

//go:wasm-module env
//export get
func envGet(proxyPtr, fieldNamePtr uint32) uint32

//go:wasm-module env
//export get_array_elem
func envGetArrayElem(proxyPtr, elemIdx uint32) uint32

//go:wasm-module env
//export get_array_size
func envGetArraySize(proxyPtr uint32) uint32

//go:wasm-module env
//export get_bool
func envGetBool(proxyPtr uint32) uint32

//go:wasm-module env
//export get_bytes
func envGetBytes(proxyPtr uint32) uint32

//go:wasm-module env
//export get_float32
func envGetFloat32(proxyPtr uint32) float32

//go:wasm-module env
//export get_float64
func envGetFloat64(proxyPtr uint32) float64

//go:wasm-module env
//export get_int16
func envGetInt16(proxyPtr uint32) uint32

//go:wasm-module env
//export get_int32
func envGetInt32(proxyPtr uint32) int32

//go:wasm-module env
//export get_int64
func envGetInt64(proxyPtr uint32) int64

//go:wasm-module env
//export get_int8
func envGetInt8(proxyPtr uint32) uint32
