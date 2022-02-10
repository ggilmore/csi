package memtable

import (
	"errors"
	"sort"
)

type DB interface {
	// Get gets the value for the given key. It returns a KeyNotFound error if the
	// DB does not contain the key.
	Get(key []byte) (value []byte, err error)

	// Has returns true if the DB contains the given key.
	Has(key []byte) (ret bool, err error)

	// Put sets the value for the given key. It overwrites any previous value
	// for that key; a DB is not a multi-map.
	Put(key, value []byte) error

	// Delete deletes the value for the given key.
	Delete(key []byte) error

	// RangeScan returns an Iterator (see below) for scanning through all
	// key-value pairs in the given range, ordered by key ascending.
	RangeScan(start, limit []byte) (Iterator, error)
}

type InMemory struct {
	hashTable map[string][]byte
}

func NewInMemory() *InMemory {
	return &InMemory{
		hashTable: make(map[string][]byte),
	}
}

var KeyNotFound = errors.New("key not found")

func (m InMemory) Get(key []byte) (value []byte, err error) {
	k := string(key)
	value, ok := m.hashTable[k]

	if !ok {
		return nil, KeyNotFound
	}

	return value, nil
}

func (m InMemory) Has(key []byte) (ret bool, err error) {
	_, ok := m.hashTable[string(key)]
	return ok, nil
}

func (m InMemory) Put(key, value []byte) error {
	k := string(key)
	m.hashTable[k] = value

	return nil
}

func (m InMemory) Delete(key []byte) error {
	delete(m.hashTable, string(key))
	return nil
}

func (m InMemory) RangeScan(start, limit []byte) (Iterator, error) {
	var entries []pair

	for k, v := range m.hashTable {
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

type Iterator interface {
	// Next moves the iterator to the next key/value pair.
	// It returns false if the iterator is exhausted.
	Next() bool

	// Error returns any accumulated error. Exhausting all the key/value pairs
	// is not considered to be an error.
	Error() error

	// Key returns the key of the current key/value pair, or nil if done.
	Key() []byte

	// Value returns the value of the current key/value pair, or nil if done.
	Value() []byte
}
