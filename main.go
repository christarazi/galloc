package main

import (
	"./galloc"
	"fmt"
	// "unsafe"
)

func main() {
	ret := (*int)(galloc.Malloc(1))
	ret2 := (*int)(galloc.Malloc(2))
	ret3 := (*int)(galloc.Malloc(3))
	ret4 := (*int)(galloc.Malloc(4))
	ret5 := (*int)(galloc.Malloc(5))

	*ret = 0x02
	*ret2 = 0x02
	*ret3 = 0x02
	*ret4 = 0x02
	*ret5 = 0x02
	fmt.Println(*ret)

	// TODO: CT
	// galloc.Free(unsafe.Pointer(ret))

	galloc.Debug()
}
