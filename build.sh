#!/bin/bash

export GOPATH=$HOME/go

go install github.com/jashper/bitmask-go/src/bitmask

go install github.com/jashper/bitmask-go/src/bitmask/address
go install github.com/jashper/bitmask-go/src/bitmask/base58
go install github.com/jashper/bitmask-go/src/bitmask/ec256k1
go install github.com/jashper/bitmask-go/src/bitmask/ecdsa
go install github.com/jashper/bitmask-go/src/bitmask/node
go install github.com/jashper/bitmask-go/src/bitmask/ripemd160