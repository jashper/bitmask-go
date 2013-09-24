package main

import (
	"fmt"
	. "github.com/jashper/bitmask-go/bitmask"
)

func main() {
	_, err := NewAddress(ADDRVER_BTM)
	fmt.Println(err)
}
