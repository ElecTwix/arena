package arena_test

import (
	"testing"

	"github.com/ElecTwix/arena"
)

func TestArenaAlloc(t *testing.T) {
	arena := arena.ArenaAlloc(4096)
	if arena == nil {
		t.Error("Failed to allocate arena")
	}

	if arena.ChunkSize != 4096 {
		t.Error("Chunk size is not 4096")
	}

	ptr := arena.Alloc(1024)
	if ptr == nil {
		t.Error("Failed to allocate memory")
	}

	// set some for testing
	for i := 0; i < 1024; i++ {
		*(*byte)(ptr) = 0x41
	}

	arena.Free()
}

func TestArenaAllocChunk(t *testing.T) {
	arena := arena.ArenaAlloc(4096)
	if arena == nil {
		t.Error("Failed to allocate arena")
	}

	for i := 0; i < 4; i++ {

		ptr := arena.Alloc(1024)
		if ptr == nil {
			t.Error("Failed to allocate memory")
		}

	}

	if arena.Current.Used != uintptr(4096) {
		t.Error("Offset is not 4096 is ", arena.Current.Used)
	}

	ptr := arena.Alloc(1024)
	if ptr == nil {
		t.Error("Failed to allocate memory")
	}

	if arena.Current.Used != 1024 {
		t.Error("Offset is not 1024")
	}

	arena.Free()
}

func TestArenaReset(t *testing.T) {
	arena := arena.ArenaAlloc(4096)
	if arena == nil {
		t.Error("Failed to allocate arena")
	}
	for i := 0; i < 4; i++ {
		ptr := arena.Alloc(1024)
		if ptr == nil {
			t.Error("Failed to allocate memory")
		}
	}
	if arena.Current.Used != uintptr(4096) {
		t.Error("Offset is not 4096 is ", arena.Current.Used)
	}
	arena.Reset()
	if arena.Current == nil {
		t.Error("Current is not nil")
	}
	arena.Free()
}

func TestArenaResetCurrent(t *testing.T) {
	arena := arena.ArenaAlloc(4096)
	if arena == nil {
		t.Error("Failed to allocate arena")
	}
	for i := 0; i < 4; i++ {
		ptr := arena.Alloc(1024)
		if ptr == nil {
			t.Error("Failed to allocate memory")
		}
	}
	if arena.Current.Used != uintptr(4096) {
		t.Error("Offset is not 4096 is ", arena.Current.Used)
	}
	arena.ResetCurrent()
	if arena.Current == nil {
		t.Error("Current is not nil")
	}
	arena.Free()
}

func TestArenaFree(t *testing.T) {
	arena := arena.ArenaAlloc(4096)
	if arena == nil {
		t.Error("Failed to allocate arena")
	}
	for i := 0; i < 4; i++ {
		ptr := arena.Alloc(1024)
		if ptr == nil {
			t.Error("Failed to allocate memory")
		}
	}
	if arena.Current.Used != uintptr(4096) {
		t.Error("Offset is not 4096 is ", arena.Current.Used)
	}
	arena.Free()
	if arena.Current != nil {
		t.Error("Current is not nil")
	}
}

func TestSetMemory(t *testing.T) {
	a := arena.ArenaAlloc(4096)
	if a == nil {
		t.Error("Failed to allocate arena")
	}
	ptr := a.Alloc(1024)
	if ptr == nil {
		t.Error("Failed to allocate memory")
	}
	arena.SetMemory(ptr, byte(0x41))
	if *(*byte)(ptr) != 0x41 {
		t.Error("Memory is not 0x41")
	}
	a.Free()
}

func TestGetMemory(t *testing.T) {
	a := arena.ArenaAlloc(4096)
	if a == nil {
		t.Error("Failed to allocate arena")
	}
	ptr := a.Alloc(1024)
	if ptr == nil {
		t.Error("Failed to allocate memory")
	}
	val := int32(111)

	arena.SetMemory(ptr, val)
	if *(*int32)(ptr) != 111 {
		t.Error("Memory is not 0x41", *(*int32)(ptr))
	}

	someVal := arena.GetMemory[int32](ptr)
	if someVal != val {
		t.Error("Memory is not 0x41", someVal)
	}

	a.Free()
}

func TestGetMemorySlice(t *testing.T) {
	a := arena.ArenaAlloc(4096)
	if a == nil {
		t.Error("Failed to allocate arena")
	}
	ptr := a.Alloc(1024)
	if ptr == nil {
		t.Error("Failed to allocate memory")
	}
	val := []int32{1, 2, 3, 4, 5}
	arena.SetMemory(ptr, val)
	if arena.GetMemorySlice[int32](ptr, 5)[0] != 1 {
		t.Error("Memory is not 1")
	}
	if arena.GetMemorySlice[int32](ptr, 5)[1] != 2 {
		t.Error("Memory is not 2")
	}
	if arena.GetMemorySlice[int32](ptr, 5)[2] != 3 {
		t.Error("Memory is not 3")
	}
	if arena.GetMemorySlice[int32](ptr, 5)[3] != 4 {
		t.Error("Memory is not 4")
	}
	if arena.GetMemorySlice[int32](ptr, 5)[4] != 5 {
		t.Error("Memory is not 5")
	}
	a.Free()
}

func TestNext(t *testing.T) {
	a := arena.ArenaAlloc(4096)
	if a == nil {
		t.Error("Failed to allocate arena")
	}
	ptr := a.Alloc(1024)
	if ptr == nil {
		t.Error("Failed to allocate memory")
	}
	arena.SetMemory(ptr, byte(0x41))
	if *(*byte)(ptr) != 0x41 {
		t.Error("Memory is not 0x41")
	}
	next := arena.Next(ptr, 1)
	if *(*byte)(next) != 0 {
		t.Error("Memory is not empty")
	}
	arena.SetMemory(next, byte(0x41))
	if *(*byte)(next) != 0x41 {
		t.Error("Memory is not 0x41")
	}
	a.Free()
}

func TestPrev(t *testing.T) {
	a := arena.ArenaAlloc(4096)
	if a == nil {
		t.Error("Failed to allocate arena")
	}
	ptr := a.Alloc(1024)
	if ptr == nil {
		t.Error("Failed to allocate memory")
	}
	arena.SetMemory(ptr, byte(0x41))
	if *(*byte)(ptr) != 0x41 {
		t.Error("Memory is not 0x41")
	}
	next := arena.Next(ptr, 1)
	if *(*byte)(next) != 0 {
		t.Error("Memory is not empty")
	}
	arena.SetMemory(next, byte(0x41))
	if *(*byte)(next) != 0x41 {
		t.Error("Memory is not 0x41")
	}
	prev := arena.Prev(next, 1)
	if *(*byte)(prev) != 0x41 {
		t.Error("Memory is not 0x41")
	}
	a.Free()
}
