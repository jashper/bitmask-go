package base58

import (
	"testing"
)

func TestConversion(t *testing.T) {

	bytes := []byte("This is a string")

	encoded, err := FromBytes(bytes)
	if err != nil {
		t.Error(err.Error())
	}

	decoded, err := ToBytes(encoded)
	if err != nil {
		t.Error(err.Error())
	}

	if string(decoded) != string(bytes) {
		t.Error("Decoded base58 does not match. Expected %s, got %s", string(bytes), string(decoded))
	}

}
