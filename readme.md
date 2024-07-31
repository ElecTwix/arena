# Go Memory Arena

## Introduction

this is project meant to how to do memory arena in go.

## How to use

```go
package main


import (
        "fmt
        "github.com/ElecTwix/arena"

)


func main() {
        a := arena.NewArena(1024) // 1kb arena
        p := a.Alloc(10) // 10 byte

        arena.SetMemory(p, []byte("hello world"))

        v = arena.GetMemory(p, 10) // v is []byte("hello world")

        a.Free(p) // free memory
}
```
