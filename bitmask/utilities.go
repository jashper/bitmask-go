package bitmask

import (
	"reflect"
	"unsafe"
)

func parseVarInt(bytes []byte, start int) (i uint64, end int) {
	max := len(bytes)

	if start >= max {
		return 0, -1
	}

	end = start + 1

	var varInt interface{}
	varInt = *(*uint8)(unsafe.Pointer(&bytes[start]))

	start++

	if varInt == uint8(0xfd) {
		end += 2

		if end > max {
			return 0, -1
		}

		varInt = *(*uint16)(unsafe.Pointer(&bytes[start]))
	} else if varInt == uint8(0xfe) {
		end += 4

		if end > max {
			return 0, -1
		}

		varInt = *(*uint32)(unsafe.Pointer(&bytes[start]))
	} else if varInt == uint8(0xff) {
		end += 8

		if end > max {
			return 0, -1
		}

		varInt = *(*uint64)(unsafe.Pointer(&bytes[start]))
	}

	return reflect.ValueOf(varInt).Uint(), end
}

func parseVarStr(bytes []byte, start int) (str *string, end int) {
	max := len(bytes)

	length, start := parseVarInt(bytes, start)
	if start == -1 {
		return nil, -1
	}

	// Assumption: length is never > 2^31 - 1 (ie: 32-bit systems won't overflow to negatives)

	end = start + int(length)
	if end > max {
		return nil, -1
	}
	str = (*string)(unsafe.Pointer(&bytes[start]))

	return str, end
}
