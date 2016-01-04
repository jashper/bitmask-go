package bitmask

import (
	"errors"
	"math/big"
)

const (
	BASE58_ALPHABET = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

var (
	BASE58_REV_ALPHABET = map[string]int64{
		"1": 0, "2": 1, "3": 2, "4": 3, "5": 4, "6": 5, "7": 6, "8": 7, "9": 8, "A": 9,
		"B": 10, "C": 11, "D": 12, "E": 13, "F": 14, "G": 15, "H": 16, "J": 17, "K": 18,
		"L": 19, "M": 20, "N": 21, "P": 22, "Q": 23, "R": 24, "S": 25, "T": 26, "U": 27,
		"V": 28, "W": 29, "X": 30, "Y": 31, "Z": 32, "a": 33, "b": 34, "c": 35, "d": 36,
		"e": 37, "f": 38, "g": 39, "h": 40, "i": 41, "j": 42, "k": 43, "m": 44, "n": 45,
		"o": 46, "p": 47, "q": 48, "r": 49, "s": 50, "t": 51, "u": 52, "v": 53, "w": 54,
		"x": 55, "y": 56, "z": 57,
	}
)

func EncodeBase58(input []byte) (string, error) {
	if len(input) < 1 {
		return "", errors.New("base58.FromBytes: Byte slice is too short")
	}
	output := ""

	n := big.NewInt(0).SetBytes(input)
	r := big.NewInt(0)
	base := big.NewInt(58)

	zero := big.NewInt(0)
	for n.Cmp(zero) == 1 {
		r.Mod(n, base)
		n.Div(n, base)
		idx := r.Int64()
		output = BASE58_ALPHABET[idx:idx+1] + output
	}

	for i := 0; i < len(input); i++ {
		if input[i] > 0 {
			break
		}
		output = "1" + output
	}

	return output, nil
}

func DecodeBase58(input string) ([]byte, error) {
	output := big.NewInt(0)
	tmp := big.NewInt(0)

	base := big.NewInt(58)
	exp := big.NewInt(0)
	val := big.NewInt(0)
	count := len(input)
	for i := 0; i < count; i++ {
		v, ok := BASE58_REV_ALPHABET[input[i:i+1]]
		if !ok {
			return nil, errors.New(
				"base58.ToBytes: Character not present in base58")
		}

		val.SetInt64(v)
		exp.SetInt64(int64(count - (i + 1)))

		tmp.Exp(base, exp, nil)
		tmp.Mul(tmp, val)
		output.Add(output, tmp)
	}

	return output.Bytes(), nil
}
