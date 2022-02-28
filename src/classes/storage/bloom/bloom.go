package main

import (
	"hash/fnv"
	"math"

	"github.com/spaolacci/murmur3"
)

type bloomFilter interface {
	add(item string)

	// `false` means the item is definitely not in the set
	// `true` means the item might be in the set
	maybeContains(item string) bool

	// Number of bytes used in any underlying storage
	memoryUsage() int
}

type trivialBloomFilter struct {
	capacity int
	bitset   *bitset
}

func newTrivialBloomFilter(capacity int, falsePositiveRate float64) *trivialBloomFilter {
	storageSize := int(math.Ceil((float64(capacity) / math.Ln2) * math.Log2(1.0/falsePositiveRate)))
	return &trivialBloomFilter{
		bitset:   newBitset(storageSize),
		capacity: capacity,
	}
}

func (b *trivialBloomFilter) add(item string) {
	i1, i2 := b.getIndexes(item)
	err := b.bitset.mark(i1)
	if err != nil {
		panic(err)
	}

	err = b.bitset.mark(i2)
	if err != nil {
		panic(err)
	}
}

func (b *trivialBloomFilter) getIndexes(item string) (int, int) {
	hash := fnv.New64()
	_, err := hash.Write([]byte(item))
	if err != nil {
		panic(err)
	}

	index := hash.Sum64() % uint64(b.bitset.capacity)

	hash2 := murmur3.New64()
	_, err = hash2.Write([]byte(item))
	if err != nil {
		panic(err)
	}

	index2 := hash2.Sum64() % uint64(b.bitset.capacity)

	return int(index), int(index2)
}

func (b *trivialBloomFilter) maybeContains(item string) bool {
	index, index2 := b.getIndexes(item)

	contains, err := b.bitset.isMarked(index)
	if err != nil {
		panic(err)
	}

	contains2, err := b.bitset.isMarked(index2)
	if err != nil {
		panic(err)
	}

	return contains && contains2
}

func (b *trivialBloomFilter) memoryUsage() int {
	return b.bitset.size()
}
