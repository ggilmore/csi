package memtable

import (
	"sort"
)

type HashIndex struct {
	hashTable map[string][]byte
}

func NewHashIndex() *HashIndex {
	return &HashIndex{
		hashTable: make(map[string][]byte),
	}
}

func (h HashIndex) Get(key []byte) (value []byte, err error) {
	k := string(key)
	value, ok := h.hashTable[k]

	if !ok {
		return nil, KeyNotFound
	}

	return value, nil
}

func (h HashIndex) Has(key []byte) (ret bool, err error) {
	_, ok := h.hashTable[string(key)]
	return ok, nil
}

func (h HashIndex) Put(key, value []byte) error {
	k := string(key)
	h.hashTable[k] = value

	return nil
}

func (h HashIndex) Delete(key []byte) error {
	delete(h.hashTable, string(key))
	return nil
}

func (h HashIndex) RangeScan(start, limit []byte) (Iterator, error) {
	var entries []pair

	for k, v := range h.hashTable {
		entries = append(entries, pair{[]byte(k), v})
	}
	sort.Slice(entries, func(i, j int) bool {
		return string(entries[i].key) < string(entries[j].key)
	})

	startIndex := 0
	for i := 0; i < len(entries); i++ {
		if string(entries[i].key) >= string(start) {
			break
		}
		startIndex++
	}

	endIndex := len(entries)
	for i := len(entries) - 1; i >= 0; i-- {
		if string(entries[i].key) < string(limit) {
			break
		}
		endIndex--
	}

	return NewSliceIterator(entries[startIndex:endIndex]), nil

}

type pair struct {
	key   []byte
	value []byte
}

type sliceIterator struct {
	storage []pair
	err     error
	index   int
}

func NewSliceIterator(storage []pair) *sliceIterator {
	return &sliceIterator{
		storage: storage,
		index:   -1,
	}
}

func (s *sliceIterator) Next() bool {
	if s.index < len(s.storage)-1 {
		s.index++
		return true
	}

	return false
}

func (s *sliceIterator) Error() error {
	return s.err
}

func (s *sliceIterator) Key() []byte {
	if s.index >= 0 && s.index < len(s.storage) {
		e := s.storage[s.index]
		return e.key
	}

	return nil
}

func (s *sliceIterator) Value() []byte {
	if s.index >= 0 && s.index < len(s.storage) {
		e := s.storage[s.index]
		return e.value
	}

	return nil
}
