package bitmask

import (
	"sync/atomic"
)

type Buffer struct {
	values [][]byte
	maxIdx int

	head      int32
	tail      int
	peerTail  int
	peerCount int32

	putReserve int32
	cleanLock  int32
}

func NewBuffer(size int) (this *Buffer) {
	this = new(Buffer)
	this.values = make([][]byte, size)
	this.maxIdx = size - 1

	this.head = 0
	this.tail = 0
	this.peerTail = 0
	this.peerCount = 0

	this.putReserve = 0
	this.cleanLock = 0

	return this
}

func (b *Buffer) Clean() {
	if b.cleanLock == 1 || !atomic.CompareAndSwapInt32(&b.cleanLock, 0, 1) {
		return
	}

	for b.tail != b.peerTail {
		b.values[b.tail] = nil

		b.tail++
		if b.tail > b.maxIdx {
			b.tail = 0
		}
	}

	b.cleanLock = 0

}

func (b *Buffer) Put(value []byte) {

	tempSlice := make([]byte, len(value))
	copy(tempSlice, value)

	for {
		putIdx := b.putReserve

		newIdx := putIdx + 1
		if int(newIdx) > b.maxIdx {
			newIdx = 0
		}

		for int(newIdx) == b.tail {
			b.Clean()
		}

		if atomic.CompareAndSwapInt32(&b.putReserve, putIdx, newIdx) {
			b.values[putIdx] = tempSlice

			for {
				temp := b.head

				// TODO: Possible race condition with writing the wrong head

				if atomic.CompareAndSwapInt32(&b.head, temp, newIdx) {
					break
				}
			}

			break
		}
	}

	b.Clean()
}

func (b *Buffer) Get(peerStart int) (values [][]byte, newStart int) {
	if peerStart == -1 { // New peer starting to read/get
		for {
			temp := b.peerCount
			if temp != -1 && atomic.CompareAndSwapInt32(&b.peerCount, temp, -1) {
				peerStart = b.peerTail
				b.peerCount = temp + 1
				break
			}
		}
	}

	if peerStart == int(b.head) {
		return nil, peerStart
	}

	newStart = int(b.head)

	isPeerTail := false
	if peerStart == b.peerTail {
		isPeerTail = true
	}

	var count int
	if peerStart > newStart {
		count = b.maxIdx - peerStart + newStart + 1
	} else {
		count = newStart - peerStart
	}

	values = make([][]byte, count)
	idx := 0
	for peerStart != newStart {
		values[idx] = make([]byte, len(b.values[peerStart]))
		copy(values[idx], b.values[peerStart])

		peerStart++
		if peerStart > b.maxIdx {
			peerStart = 0
		}
		idx++
	}

	for isPeerTail == true {
		temp := b.peerCount
		if atomic.CompareAndSwapInt32(&b.peerCount, temp, temp-1) {
			break
		}
	}

	if atomic.CompareAndSwapInt32(&b.peerCount, 0, -1) {
		b.peerTail = newStart
		b.peerCount = 1
	}

	return values, newStart
}
