package main

import (
	"fmt"
	. "github.com/jashper/bitmask-go/bitmask"
	"runtime"
	"strconv"
	"sync"
)

func testA(b *Buffer) {
	lA := b.SpawnListener()
	lB := b.SpawnListener()

	for i := 0; i < 15; i++ {
		lB.Put([]byte(strconv.FormatInt(int64(i), 10)))

		if i == 8 {
			testB(lA)
			testB(lB)
		}
	}

	lC := b.SpawnListener()
	testB(lA)
	testB(lC)
	testB(lB)

}

func testB(b *BufferListener) {
	values := b.Get()
	fmt.Println(values)
}

func main() {
	runtime.GOMAXPROCS(2)

	buffer := NewBuffer(10)

	go testA(buffer)

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
