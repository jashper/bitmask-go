package main

import (
	"fmt"
	. "github.com/jashper/bitmask-go/bitmask"
	"runtime"
	"strconv"
	"sync"
)

func testA(b *Buffer) {
	for i := 0; i < 11; i++ {
		b.Put([]byte(strconv.FormatInt(int64(i), 10)))

		if i == 8 {
			go testB(b)
		}
	}

	b.Get(-1)
	//fmt.Println(len(values))

}

func testB(b *Buffer) {
	values, _ := b.Get(-1)

	fmt.Println(len(values))
}

func main() {
	runtime.GOMAXPROCS(2)

	buffer := NewBuffer(10)

	go testA(buffer)

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
