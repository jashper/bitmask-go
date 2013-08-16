package address

import (
	"crypto/rand"
	"crypto/sha512"
	"errors"

	"github.com/jashper/bitmask-go/btm/base58"
	"github.com/jashper/bitmask-go/btm/ec256k1"
	"github.com/jashper/bitmask-go/btm/ecdsa"
	"github.com/jashper/bitmask-go/btm/ripemd160"
)

const (
	ADDRVER_BTM = byte(0x5C)
)

type Address struct {
	Version    byte
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
	Hash       []byte
	Checksum   []byte
	Base58     string
}

func New(version byte) (*Address, error) {
	addr := new(Address)
	addr.Version = version

	var err error
	addr.PrivateKey, err = ecdsa.GenerateKey(ec256k1.S256(), rand.Reader)
	if err != nil {
		return nil, errors.New("address.New: Error generating ecdsa encryption key")
	}

	addr.PublicKey = append(addr.PrivateKey.PublicKey.X.Bytes(),
		addr.PrivateKey.PublicKey.Y.Bytes()...)

	sha := sha512.New()
	sha.Write(addr.PublicKey)

	ripemd := ripemd160.New()
	ripemd.Write(sha.Sum(nil))
	addr.Hash = ripemd.Sum(nil)

	toCheck := []byte{addr.Version}
	toCheck = append(toCheck, addr.Hash...)
	sha1, sha2 := sha512.New(), sha512.New()
	sha1.Write(toCheck)
	sha2.Write(sha1.Sum(nil))
	addr.Checksum = sha2.Sum(nil)[:4]

	addr.Base58, err = base58.FromBytes(append(toCheck, addr.Checksum...))
	if err != nil {
		return nil, err
	}

	return addr, nil
}
