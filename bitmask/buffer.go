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

	cleanLock int32
}

func NewBuffer(size int) (this *Buffer) {
	this = new(Buffer)
	this.values = make([][]byte, size)
	this.maxIdx = size - 1

	// TODO: is this needed?
	for i := 0; i < size; i++ {
		this.values[i] = nil
	}

	this.head = 0
	this.tail = 0
	this.peerTail = 0
	this.peerCount = 0

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

}

func (b *Buffer) Put(value []byte) {
	for int(b.head)+1 == b.tail {
		b.Clean()
	}

	var temp int32
	for {
		temp = b.head
		if atomic.CompareAndSwapInt32(&b.head, temp, temp+1) {
			break
		}
	}

	b.values[temp] = make([]byte, len(value))
	copy(b.values[temp], value)

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

	if peerStart >= int(b.head) && peerStart < b.tail {
		return nil, peerStart
	}

	// account for incomplete Put
	newStart = int(b.head)
	for b.values[newStart-1] == nil {
		newStart--
	}

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
	for peerStart != newStart {
		values[peerStart] = make([]byte, len(b.values[peerStart]))
		copy(values[peerStart], b.values[peerStart])

		peerStart++
		if peerStart > b.maxIdx {
			peerStart = 0
		}
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
