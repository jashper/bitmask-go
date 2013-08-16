package address

import (
	"fmt"
	"testing"
)

func TestAddress(t *testing.T) {
	addr, err := New(ADDRVER_BTM)
	if err != nil {
		t.Error(err.Error())
	}

	fmt.Println(addr.Base58)
}
