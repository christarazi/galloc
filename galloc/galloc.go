package galloc

import (
	"golang.org/x/sys/unix"
	// "reflect"
	"fmt"
	"syscall"
	"unsafe"
)

type MemBlock struct {
	Freed bool
	Size  int
	Next  *MemBlock
	Prev  *MemBlock
	Data  unsafe.Pointer
}

const SYSBRK uintptr = 12
const ALIGN uintptr = 64

var base *MemBlock

func brk(size uintptr) (uintptr, syscall.Errno) {
	r, _, err := unix.Syscall(SYSBRK, size, 0, 0)
	return r, err
}

func aligned_size(size int) int {
	size = size - 1
	size |= size >> 1
	size |= size >> 2
	size |= size >> 4
	size |= size >> 8
	size |= size >> 16
	size = size + 1

	return size
}

func extend(size int, prev *MemBlock) *MemBlock {
	var block *MemBlock

	ptr, err := brk(0)
	if err != 0 {
		panic(err)
	}
	fmt.Printf("base: 0x%x\n", ptr)

	// Hack to get aligned memory blocks. Otherwise, Go will puke out something
	// like this:
	// fatal error: bulkBarrierPreWrite: unaligned arguments
	aligned := aligned_size(size)
	for aligned%ALIGN != 0 {
		aligned = aligned_size(aligned + 1)
	}
	println("allocating bytes:", aligned)

	end, err := brk(uintptr(aligned) + ptr)
	if err != 0 {
		panic(err)
	}
	fmt.Printf("new end: 0x%x\n", end)

	block = (*MemBlock)(unsafe.Pointer(ptr))
	block.Freed = false
	block.Size = size
	block.Next = nil
	block.Prev = prev

	// Do not forget to set the next block from the previous.
	if prev != nil {
		prev.Next = block
	}

	return block
}

func Debug() {
	b := base
	for b != nil {
		fmt.Println(*b, "->")
		b = b.Next
	}
}

func Malloc(size int) unsafe.Pointer {
	var block, prev *MemBlock

	if base == nil {
		block = extend(size, prev)
		fmt.Println("first call to brk:", block)
		base = block
	} else {
		prev = base
		cur := base
		for cur != nil && cur.Freed == false && cur.Size <= size {
			prev = cur
			cur = cur.Next
		}

		if cur == nil {
			block = extend(size, prev)
		} else {
			block = cur
		}
	}

	return unsafe.Pointer(&block.Data)
}

// TODO: CT
// func Free(ptr unsafe.Pointer) {
// 	b := (*MemBlock)(unsafe.Pointer(ptr))
// 	b.Freed = true
// }