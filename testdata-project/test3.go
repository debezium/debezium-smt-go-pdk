package main

import (
	"github.com/debezium/debezium-smt-go-pdk"
)

//export process
func process(proxyPtr uint32) uint32 {
	str := debezium.GetString(proxyPtr)
	return debezium.SetString(str)
}

func main() {}
