package bitmask

import (
	"crypto/rand"
	"crypto/sha512"
	"errors"

	"github.com/jashper/bitmask-go/bitmask/ec256k1"
	"github.com/jashper/bitmask-go/bitmask/ecdsa"
	"github.com/jashper/bitmask-go/bitmask/ripemd160"
)

const (
	ADDRVER_BTM = byte(0x5C)
)

type Address struct {
	Version    byte
	PrivateKey *ecdsa.PrivateKey
	Base58     string
}

func NewAddress(version byte) (*Address, error) {
	addr := new(Address)
	addr.Version = version

	var err error
	addr.PrivateKey, err = ecdsa.GenerateKey(ec256k1.S256(), rand.Reader)
	if err != nil {
		return nil, errors.New("address.New: Error generating ecdsa encryption key")
	}

	publicKey := append(addr.PrivateKey.PublicKey.X.Bytes(),
		addr.PrivateKey.PublicKey.Y.Bytes()...)

	sha := sha512.New()
	sha.Write(publicKey)

	ripemd := ripemd160.New()
	ripemd.Write(sha.Sum(nil))
	hash := ripemd.Sum(nil)

	toCheck := []byte{addr.Version}
	toCheck = append(toCheck, hash...)
	sha1, sha2 := sha512.New(), sha512.New()
	sha1.Write(toCheck)
	sha2.Write(sha1.Sum(nil))
	checksum := sha2.Sum(nil)[:4]

	addr.Base58, err = EncodeBase58(append(toCheck, checksum...))
	if err != nil {
		return nil, err
	}

	return addr, nil
}
