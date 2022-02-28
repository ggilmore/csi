package main

import "fmt"

type bitset struct {
	capacity int
	data     []uint8
}

func newBitset(numElements int) *bitset {
	length := numElements / 8
	if (numElements % 8) > 0 {
		length++
	}

	return &bitset{
		capacity: numElements,
		data:     make([]uint8, length),
	}
}

func (b *bitset) mark(i int) error {
	if i < 0 || i >= b.capacity {
		return fmt.Errorf("index %d out of bounds [0, %d)", i, b.capacity)
	}

	b.data[i/8] |= 1 << (i % 8)
	return nil
}

func (b *bitset) isMarked(i int) (bool, error) {
	if i < 0 || i >= b.capacity {
		return false, fmt.Errorf("index %d out of bounds [0, %d)", i, b.capacity)
	}

	return b.data[i/8]&(1<<(i%8)) != 0, nil
}

func (b *bitset) size() int {
	return b.capacity
}
