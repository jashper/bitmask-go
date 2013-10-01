package bitmask

import (
	"sync/atomic"
)

type Buffer struct {
	values [][]byte
	maxIdx int

	head     int32
	peerTail int
	tail     int

	peerCount    int
	voteCount    int
	voteRound    int
	nextPeerTail int

	putReserve int32
	cleanLock  int32
	voteLock   int32
}

type BufferListener struct {
	idx    int
	buffer *Buffer

	round    int
	hasVoted bool
}

func NewBuffer(size int) (this *Buffer) {
	this = new(Buffer)
	this.values = make([][]byte, size)
	this.maxIdx = size - 1

	this.head = 0
	this.peerTail = 0
	this.tail = 0

	this.peerCount = 0
	this.voteCount = 0
	this.voteRound = 0
	this.nextPeerTail = 0

	this.putReserve = 0
	this.cleanLock = 0
	this.voteLock = 0

	return this
}

func (b *Buffer) SpawnListener() (bl *BufferListener) {
	bl = new(BufferListener)
	bl.buffer = b
	bl.hasVoted = false

	for {
		if atomic.CompareAndSwapInt32(&b.voteLock, 0, 1) {
			bl.round = b.voteRound
			bl.idx = b.peerTail
			b.peerCount++

			b.voteLock = 0
			break
		}
	}

	return bl
}

func (bl *BufferListener) Stop() {
	b := bl.buffer
	for {
		if atomic.CompareAndSwapInt32(&b.voteLock, 0, 1) {
			b.peerCount--
			b.voteLock = 0
			break
		}
	}
}

func (bl *BufferListener) Get() (values [][]byte) {
	b := bl.buffer

	values, bl.idx = b.get(bl.idx)
	bl.vote()

	return values
}

func (bl *BufferListener) vote() {
	b := bl.buffer

	if bl.round != b.voteRound {
		bl.round = b.voteRound
		bl.hasVoted = false
	}

	if !bl.hasVoted {
		for {
			if atomic.CompareAndSwapInt32(&b.voteLock, 0, 1) {
				if b.voteCount == 0 {
					b.nextPeerTail = bl.idx
					b.voteCount++
				} else {
					if bl.idx < b.nextPeerTail {
						b.nextPeerTail = bl.idx
						bl.round++
						b.voteCount = 1
					} else {
						b.voteCount++
					}
				}

				bl.hasVoted = true

				b.voteLock = 0
				break
			}
		}
	}
}

func (b *Buffer) clean() {
	if b.cleanLock == 1 || !atomic.CompareAndSwapInt32(&b.cleanLock, 0, 1) {
		return
	}

	if b.voteCount >= b.peerCount {
		for {
			if atomic.CompareAndSwapInt32(&b.voteLock, 0, 1) {
				b.peerTail = b.nextPeerTail
				b.voteCount = 0
				b.voteRound++

				b.voteLock = 0
				break
			}
		}
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

func (bl *BufferListener) Put(value []byte) {
	b := bl.buffer

	tempSlice := make([]byte, len(value))
	copy(tempSlice, value)

	for {
		putIdx := b.putReserve

		newIdx := putIdx + 1
		if int(newIdx) > b.maxIdx {
			newIdx = 0
		}

		for int(newIdx) == b.tail {
			b.clean()
		}

		if atomic.CompareAndSwapInt32(&b.putReserve, putIdx, newIdx) {
			b.values[putIdx] = tempSlice

			for {
				temp := b.head

				if atomic.CompareAndSwapInt32(&b.head, temp, newIdx) {
					break
				}
			}

			break
		}
	}

	b.clean()
}

func (b *Buffer) get(start int) (values [][]byte, end int) {
	if start == int(b.head) {
		return nil, start
	}

	end = int(b.head)

	var count int
	if start > end {
		count = b.maxIdx - start + end + 1
	} else {
		count = end - start
	}

	values = make([][]byte, count)
	idx := 0
	for start != end {
		values[idx] = make([]byte, len(b.values[start]))
		copy(values[idx], b.values[start])

		start++
		if start > b.maxIdx {
			start = 0
		}
		idx++
	}

	return values, end
}
