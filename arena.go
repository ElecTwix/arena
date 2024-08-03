package arena

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

type Arena struct {
	// The size of the chunk
	// The Current
	Current   *ArenaChunk
	ChunkSize uintptr
}

type ArenaChunk struct {
	Prev   *ArenaChunk
	Cursor unsafe.Pointer
	Used   uintptr
}

func ArenaAlloc(chunkSize uintptr) *Arena {
	arena := &Arena{
		ChunkSize: chunkSize,
		Current:   nil,
	}

	sizePtr := uintptr(chunkSize)

	arena.Current = arena.allocChunk(sizePtr, arena)

	return arena
}

func (a *Arena) allocChunk(size uintptr, arenaPtr *Arena) *ArenaChunk {
	ptr, err := unix.MmapPtr(0, 0, nil, size, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_ANONYMOUS|unix.MAP_PRIVATE)
	if err != nil {
		return nil
	}

	chunk := &ArenaChunk{
		Used:   0,
		Cursor: ptr,
		Prev:   nil,
	}

	arenaPtr.Current = chunk

	return chunk
}

func (a *Arena) Alloc(size uintptr) unsafe.Pointer {
	if a.ChunkSize < uintptr(a.Current.Used)+size {
		arena := a.allocChunk(size, a)
		a.Current.Cursor = arena.Cursor
	}

	a.Current.Used += size
	a.Current.Cursor = unsafe.Add(a.Current.Cursor, size)
	return a.Current.Cursor
}

func (a *Arena) Free() {
	arenaChunk := a.Current
	for arenaChunk != nil {
		unix.MunmapPtr(arenaChunk.Cursor, a.ChunkSize)
		arenaChunk = arenaChunk.Prev
	}
	a.Current = nil
}

func (a *Arena) Reset() {
	arenaChunk := a.Current
	for arenaChunk.Prev != nil {
		arenaChunk.Used = 0
		arenaChunk = arenaChunk.Prev
	}
}

func (a *Arena) ResetCurrent() {
	arenaChunk := a.Current
	unix.MunmapPtr(arenaChunk.Cursor, a.ChunkSize)
}

// Helper function to allocate memory
func SetMemory[T any](ptr unsafe.Pointer, value T) {
	*(*T)(ptr) = value
}

func GetMemory[T any](ptr unsafe.Pointer) T {
	return *(*T)(ptr)
}

func GetMemorySlice[T any](ptr unsafe.Pointer, size int) []T {
	return *(*[]T)(ptr)
}

func Next(ptr unsafe.Pointer, size int) unsafe.Pointer {
	return unsafe.Add(ptr, size)
}

func Prev(ptr unsafe.Pointer, size int) unsafe.Pointer {
	return unsafe.Add(ptr, size*-1)
}
