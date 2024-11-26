package debezium
// Debezium PDK POC
// TODO: test me properly

// low level primitives to be used:

import (
	"unsafe"
	"strconv"
)

// alloc/free implementation from:
// https://github.com/tinygo-org/tinygo/blob/2a76ceb7dd5ea5a834ec470b724882564d9681b3/src/runtime/arch_tinygowasm_malloc.go#L7
var allocs = make(map[uintptr][]byte)

//export malloc
func libc_malloc(size uintptr) unsafe.Pointer {
	if size == 0 {
		return nil
	}
	buf := make([]byte, size)
	ptr := unsafe.Pointer(&buf[0])
	allocs[uintptr(ptr)] = buf
	return ptr
}

//export free
func libc_free(ptr unsafe.Pointer) {
	if ptr == nil {
		return
	}
	if _, ok := allocs[uintptr(ptr)]; ok {
		delete(allocs, uintptr(ptr))
	} else {
		panic("free: invalid pointer")
	}
}

// we decide to use C string format as it requires a single int pointer

func readCString(offset uint32) string {
	length := 0
	for {
		s := *(*int32)(unsafe.Pointer(uintptr(offset) + uintptr(length)))
		if byte(s) == 0 {
			break
		}
		length++
	}

	buffer := make([]byte, length)
	for i := 0; i < int(length); i++ {
		s := *(*int32)(unsafe.Pointer(uintptr(offset) + uintptr(i)))
		buffer[i] = byte(s)
	}
	return string(buffer)
}

// inspired by:
// https://github.com/tinygo-org/tinygo/blob/2a76ceb7dd5ea5a834ec470b724882564d9681b3/src/runtime/string.go#L278
func writeCString(offset uintptr, str string) {
	stringData := []byte(str)
	for i := 0; i < len(stringData); i++ {
		*(*byte)(unsafe.Pointer(uintptr(offset) + uintptr(i))) = stringData[i]
	}
	*(*byte)(unsafe.Pointer(uintptr(offset) + uintptr(len(stringData)))) = 0 // trailing 0 byte
}

// this is the "Debezium Guest SDK" implementation
// wrapping together the low level primitives and using the
// "Debezium Host SDK" functionality

//go:wasm-module env
//export get_string
func envGetString(proxyPtr uint32) uint32

func GetString(proxyPtr uint32) string {
	var resultPtr = envGetString(proxyPtr)
	var result = readCString(resultPtr)
	libc_free(unsafe.Pointer(uintptr(resultPtr)))
	return result
}

//go:wasm-module env
//export get_uint32
func envGetUInt32(proxyPtr uint32) uint32

func GetUInt32(proxyPtr uint32) uint32 {
	return envGetUInt32(proxyPtr)
}

//go:wasm-module env
//export get
func envGet(proxyPtr, fieldNamePtr uint32) uint32

func Get(proxyPtr uint32, fieldName string) uint32 {
	var fieldNameLen = len(fieldName) + 1
	var fieldNamePtr = libc_malloc(uintptr(fieldNameLen))
	writeCString(uintptr(fieldNamePtr), fieldName)

	return envGet(proxyPtr, uint32(uintptr(fieldNamePtr)))
}

//go:wasm-module env
//export set_bool
func envSetBool(valuePtr uint32) uint32

func SetBool(value bool) uint32 {
	var valuePtr = libc_malloc(1)

	if (value) {
		*(*byte)(unsafe.Pointer(uintptr(valuePtr))) = 1
	} else {
		*(*byte)(unsafe.Pointer(uintptr(valuePtr))) = 0
	}

	return envSetBool(uint32(uintptr(valuePtr)))
}

//go:wasm-module env
//export set_int
func envSetInt(valuePtr uint32) uint32

func SetInt(value uint32) uint32 {
	var valuePtr = libc_malloc(4)
	bs := []byte(strconv.Itoa(int(value)))

	for i := 0; i < len(bs); i++ {
		*(*byte)(unsafe.Pointer(uintptr(valuePtr) + uintptr(i))) = bs[i]
	}

	return envSetInt(uint32(uintptr(valuePtr)))
}

//go:wasm-module env
//export set_string
func envSetString(valuePtr uint32) uint32

func SetString(value string) uint32 {
	var valueLen = len(value) + 1
	var valuePtr = libc_malloc(uintptr(valueLen))
	writeCString(uintptr(valuePtr), value)
	return envSetString(uint32(uintptr(valuePtr)))
}

//go:wasm-module env
//export is_null
func envIsNull(valuePtr uint32) uint32

func IsNull(valuePtr uint32) bool {
	return (envIsNull(valuePtr) > 0)
}

//go:wasm-module env
//export set_null
func envSetNull() uint32

func SetNull() uint32 {
	return envSetNull()
}
