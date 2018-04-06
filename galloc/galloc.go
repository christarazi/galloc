package galloc

import (
	"fmt"
	"golang.org/x/sys/unix"
	"syscall"
	"unsafe"
)

type MemBlock struct {
	Freed bool
	Size  int
	Prev  *MemBlock
	Next  *MemBlock
	Data  unsafe.Pointer
}

var base *MemBlock

const (
	SYSBRK    uintptr = 12
	ALIGN     int     = 64
	BLOCKSIZE uintptr = unsafe.Sizeof(*base)
)

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
	//   fatal error: bulkBarrierPreWrite: unaligned arguments
	aligned := aligned_size(size + int(BLOCKSIZE))
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
	block.Prev = prev
	block.Next = nil

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

func Free(ptr unsafe.Pointer) {
	b := uintptr(ptr) - (unsafe.Offsetof(base.Data) -
		unsafe.Offsetof(base.Freed)) // Retrieve start of MemBlock.

	// Mark block as freed.
	block := (*MemBlock)(unsafe.Pointer(b))
	block.Freed = true

	// Check if it is the last (end) block (farthest to the right if you think
	// of a linked list). If it is, make sure to destroy the any links to it,
	// and call brk(2) to deallocate the memory back.
	if block.Next == nil {
		if block.Prev != nil {
			block.Prev.Next = nil
		}
		block.Prev = nil

		p, err := brk(b)
		if err != 0 {
			panic(err)
		}

		// Don't forget to set the base pointer to nil as well.
		if block == base {
			base = nil
		}
		block = nil

		fmt.Printf("deallocated : 0x%x\n", p)
	}
}
