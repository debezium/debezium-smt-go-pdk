package debezium

import (
	"strconv"
	"unsafe"

	"github.com/andreaTP/debezium-smt-go-pdk/internal"
)

// access a nested field in the record data structure provided by Debezium
func Get(proxyPtr uint32, fieldName string) uint32 {
	var fieldNameLen = len(fieldName) + 1
	var fieldNamePtr = internal.Malloc(uintptr(fieldNameLen))
	internal.WriteCString(uintptr(fieldNamePtr), fieldName)

	return envGet(proxyPtr, uint32(uintptr(fieldNamePtr)))
}

// materialize the String content referenced
func GetString(proxyPtr uint32) string {
	var resultPtr = envGetString(proxyPtr)
	var result = internal.ReadCString(resultPtr)
	internal.Free(unsafe.Pointer(uintptr(resultPtr)))
	return result
}

// materialize the Numeric content referenced
func GetInt(proxyPtr uint32) uint32 {
	return envGetInt(proxyPtr)
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

// set a Numeric content for the Debezium Host
func SetInt(value uint32) uint32 {
	bs := []byte(strconv.Itoa(int(value)))
	var valuePtr = internal.Malloc(uintptr(len(bs) + 1))
	internal.WriteCString(uintptr(valuePtr), string(bs))

	return envSetInt(uint32(uintptr(valuePtr)))
}

//go:wasm-module env
//export get_string
func envGetString(proxyPtr uint32) uint32

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
//export set_int
func envSetInt(valuePtr uint32) uint32

//go:wasm-module env
//export set_bool
func envSetBool(valuePtr uint32) uint32

//go:wasm-module env
//export get
func envGet(proxyPtr, fieldNamePtr uint32) uint32

//go:wasm-module env
//export get_int
func envGetInt(proxyPtr uint32) uint32
