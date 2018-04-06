package main

import (
	"./galloc"
	"fmt"
	"unsafe"
)

func main() {
	ret := (*int)(galloc.Malloc(1))
	ret2 := (*int)(galloc.Malloc(2))
	ret3 := (*int)(galloc.Malloc(3))
	ret4 := (*int)(galloc.Malloc(4))
	ret5 := (*int)(galloc.Malloc(512))

	*ret = 0x01
	*ret2 = 0x02
	*ret3 = 0x03
	*ret4 = 0x04
	*ret5 = 0x05
	fmt.Println(*ret)

	galloc.Debug()

	galloc.Free(unsafe.Pointer(ret5))
	galloc.Free(unsafe.Pointer(ret4))
	galloc.Free(unsafe.Pointer(ret3))
	galloc.Free(unsafe.Pointer(ret2))
	galloc.Free(unsafe.Pointer(ret))

	ret = (*int)(galloc.Malloc(1))
	galloc.Debug()
	galloc.Free(unsafe.Pointer(ret))
	galloc.Debug()
}
